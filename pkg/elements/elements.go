package elements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
)

type Renderable interface {
	HTML() template.HTML
	Verify() error
}

func html[T interface{ getType() string }](m map[string]string, r T, index int) template.HTML {
	filepath, inMap := m[r.getType()]
	if !inMap {
		panic(fmt.Sprintf("unknown background type: %v", r.getType()))
	}

	t, err := template.ParseFiles(filepath)
	if err != nil {
		panic(fmt.Sprintf("cannot parse template file: %v", err))
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, struct {
		Value T
		Index int
	}{r, index})
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return template.HTML(buf.String())
}

func verify[T interface{ getType() string }](m map[string]string, r T) error {
	ty := r.getType()
	if ty == "variable" {
		ty = "text"
	}
	filepath, inMap := m[ty]
	if !inMap {
		return fmt.Errorf("unknown background type: %v", r.getType())
	}

	if _, err := os.Stat(filepath); err != nil {
		return err
	}

	t, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, struct {
		Value T
		Index int
	}{r, 0})
	if err != nil {
		return err
	}

	return nil
}

type Element struct {
	Type    string         `json:"type"`
	Content string         `json:"content"`
	Style   template.CSS   `json:"style"`
	Link    template.URL   `json:"link"`
	Appear  int            `json:"appear"`
	Hide    int            `json:"hide"`
	Other   map[string]any `json:"-"`
}

func (e *Element) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &e.Other)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(e).Elem()
	t := reflect.TypeOf(e).Elem()
	for i := range t.NumField() {
		if t.Field(i).Name == "Other" {
			continue
		}

		key := t.Field(i).Tag.Get("json")
		if reflect.TypeOf(e.Other[key]) != reflect.TypeOf(nil) {
			v.Field(i).Set(reflect.ValueOf(e.Other[key]).Convert(t.Field(i).Type))
		}
		delete(e.Other, key)
	}
	return nil
}

func (e Element) getType() string {
	return e.Type
}

func (e Element) HTML(index int) template.HTML {
	return html(
		map[string]string{
			"image": "./assets/templates/elements/image.html",
			"video": "./assets/templates/elements/video.html",
			"text":  "./assets/templates/elements/text.html",
			"input": "./assets/templates/elements/input.html",
			"jump":  "./assets/templates/elements/jump.html",
		},
		e,
		index,
	)
}

func (e Element) Verify() error {
	return verify(
		map[string]string{
			"image": "./assets/templates/elements/image.html",
			"video": "./assets/templates/elements/video.html",
			"text":  "./assets/templates/elements/text.html",
			"input": "./assets/templates/elements/input.html",
			"jump":  "./assets/templates/elements/jump.html",
		},
		e,
	)
}
