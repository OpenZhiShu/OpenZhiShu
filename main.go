package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"OpenZhiShu/pkg/config"
	"OpenZhiShu/pkg/drawing"
)

func genHandleFunc(filepath string, data any) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath)
		if err != nil {
			fmt.Printf("eror: %g\n", err)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			fmt.Printf("eror: %g\n", err)
		}
	}
}

func genDrawingHandleFunc(cfg config.Config, drawingData *drawing.Data[int], list *List) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		number, err := strconv.Atoi(r.PathValue("number"))
		if err != nil {
			return
		}
		result, err := drawingData.Draw(number)
		if err != nil {
			genHandleFunc("./assets/templates/error.html", err.Error())(w, r)
			return
		}
		fmt.Printf("number: %v, result: %v, len: %v\n", number, result, len(drawingData.Results()))

		names := make([]string, len(result))
		for i, v := range result {
			names[i] = list.Seniors[v]
		}

		variables := map[string]string{}
		variables["result"] = strings.Join(names, " & ")

		elems := slices.Clone(cfg.Elements)
		for i := range elems {
			if elems[i].Type != "variable" {
				continue
			}
			elems[i].Type = "text"
			value, inMap := variables[elems[i].Content]
			if !inMap {
				elems[i].Content = fmt.Sprintf("no variable `%v`", elems[i].Content)
				continue
			}
			elems[i].Content = value
		}

		newCfg := config.Config{BodyColor: cfg.BodyColor, Ratio: cfg.Ratio, Elements: elems}

		genHandleFunc("./assets/templates/page.html", newCfg)(w, r)
	}
}

type Config struct {
	HomepageConfig config.Config `json:"homepage"`
	DrawingConfig  config.Config `json:"drawing"`
	ResultConfig   config.Config `json:"result"`
}

func (c Config) Verify() error {
	if err := c.HomepageConfig.Verify(); err != nil {
		return err
	}
	if err := c.DrawingConfig.Verify(); err != nil {
		return err
	}
	if err := c.ResultConfig.Verify(); err != nil {
		return err
	}
	return nil
}

type List struct {
	Freshmen map[int]string
	Seniors  map[int]string
}

func main() {
	fmt.Print("choose a port to listen: ")
	var port int
	_, err := fmt.Scanln(&port)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		fmt.Println("you can use something like: `choose a port to listen: 8080`")
		return
	}
	fmt.Printf("http://localhost:%v\n", port)

	cfg, err := config.LoadConfig[Config]("./config.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	list, err := LoadList("./list.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	drawingData := drawing.MakeData(
		slices.Collect(maps.Keys(list.Freshmen)),
		slices.Collect(maps.Keys(list.Seniors)),
	)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/{$}", genHandleFunc("./assets/templates/page.html", cfg.HomepageConfig))
	http.HandleFunc("/drawing", genHandleFunc("./assets/templates/page.html", cfg.DrawingConfig))
	http.HandleFunc("/result/{number}", genDrawingHandleFunc(cfg.ResultConfig, &drawingData, &list))
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func LoadList(filepath string) (List, error) {
	listFile, err := os.ReadFile("./list.json")
	if err != nil {
		return List{}, err
	}

	var tmpList struct {
		Freshmen map[string]string `json:"freshmen"`
		Seniors  map[string]string `json:"seniors"`
	}
	err = json.Unmarshal(listFile, &tmpList)
	if err != nil {
		return List{}, err
	}

	list := List{
		Freshmen: make(map[int]string, len(tmpList.Freshmen)),
		Seniors:  make(map[int]string, len(tmpList.Freshmen)),
	}
	for key := range tmpList.Freshmen {
		i, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		list.Freshmen[i] = tmpList.Freshmen[key]
	}
	for key := range tmpList.Seniors {
		i, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		list.Seniors[i] = tmpList.Seniors[key]
	}

	return list, nil
}
