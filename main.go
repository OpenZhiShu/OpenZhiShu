package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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

func genDrawingHandleFunc(config config.DynamicConfig, drawingData *drawing.Data[int]) func(http.ResponseWriter, *http.Request) {
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

		variables := map[string]string{}
		variables["result"] = fmt.Sprintf("%v", result)

		for i := range config.Elements {
			if config.Elements[i].Type != "variable" {
				continue
			}
			config.Elements[i].Type = "text"
			value, inMap := variables[config.Elements[i].Content]
			if !inMap {
				config.Elements[i].Content = fmt.Sprintf("no variable `%v`", config.Elements[i].Content)
				continue
			}
			config.Elements[i].Content = value
		}

		genHandleFunc("./assets/templates/result.html", config)(w, r)
	}
}

func DrawingResult(cfg config.DynamicConfig, r *http.Request) (any, error) {
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		return nil, err
	}
	d := drawing.MakeData([]int{}, []int{})
	result, err := d.Draw(number)
	if err != nil {
		return nil, err
	}
	_ = result
	return cfg, nil
}

type List struct {
	Freshmen []int `json:"freshmen"`
	Seniors  []int `json:"seniors"`
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

	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	listFile, err := os.ReadFile("./list.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	var list List
	err = json.Unmarshal(listFile, &list)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	drawingData := drawing.MakeData(list.Freshmen, list.Seniors)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/{$}", genHandleFunc("./assets/templates/homepage.html", cfg.HomepageConfig))
	http.HandleFunc("/drawing", genHandleFunc("./assets/templates/drawing.html", cfg.DrawingConfig))
	http.HandleFunc("/result/{number}", genDrawingHandleFunc(cfg.ResultConfig, &drawingData))
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
