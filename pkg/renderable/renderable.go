package renderable

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

type Renderable interface {
	HTML() template.HTML
	Verify() error
}

func html[T interface{ getType() string }](m map[string]string, r T) template.HTML {
	filepath, inMap := m[r.getType()]
	if !inMap {
		panic(fmt.Sprintf("unknown background type: %v", r.getType()))
	}

	t, err := template.ParseFiles(filepath)
	if err != nil {
		panic(fmt.Sprintf("cannot parse template file: %v", err))
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, r)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return template.HTML(buf.String())
}

func verify[T interface{ getType() string }](m map[string]string, r T) error {
	filepath, inMap := m[r.getType()]
	if !inMap {
		return fmt.Errorf("unknown background type: %v", r.getType())
	}

	if _, err := os.Stat(filepath); err != nil {
		return err
	}

	t, err := template.ParseFiles(filepath)
	if err != nil {
		return fmt.Errorf("cannot parse template file: %v", err)
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, r)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	return nil
}

type Background struct {
	Type     string       `json:"type"`
	Content  string       `json:"content"`
	Style    template.CSS `json:"style"`
	Autoplay bool         `json:"autoplay"`
	Loop     bool         `json:"loop"`
	Muted    bool         `json:"muted"`
}

func (b Background) getType() string {
	return b.Type
}

func (b Background) HTML() template.HTML {
	return html(
		map[string]string{
			"image": "./assets/templates/renderable/background_image.html",
			"video": "./assets/templates/renderable/background_video.html",
		},
		b,
	)
}

func (b Background) Verify() error {
	return verify(
		map[string]string{
			"image": "./assets/templates/renderable/background_image.html",
			"video": "./assets/templates/renderable/background_video.html",
		},
		b,
	)
}

type Element struct {
	Type    string       `json:"type"`
	Content string       `json:"content"`
	Layout  template.CSS `json:"layout"`
	Style   template.CSS `json:"style"`
	Link    string       `json:"link"`
}

func (e Element) getType() string {
	return e.Type
}

func (e Element) HTML() template.HTML {
	return html(
		map[string]string{
			"image": "./assets/templates/renderable/element_image.html",
		},
		e,
	)
}

func (e Element) Verify() error {
	return verify(
		map[string]string{
			"image": "./assets/templates/renderable/element_image.html",
		},
		e,
	)
}
