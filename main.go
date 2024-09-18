package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

func genHandleFuncEx[T any](filepath string, data T, fn func(T, *http.Request) (any, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath)
		if err != nil {
			fmt.Printf("eror: %g\n", err)
			return
		}

		d, err := fn(data, r)
		if err != nil {
			fmt.Printf("eror: %g\n", err)
			return
		}

		err = t.Execute(w, d)
		if err != nil {
			fmt.Printf("eror: %g\n", err)
		}
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/{$}", genHandleFunc("./assets/templates/homepage.html", cfg.HomepageConfig))
	http.HandleFunc("/drawing", genHandleFunc("./assets/templates/drawing.html", cfg.DrawingConfig))
	http.HandleFunc("/result/{number}", genHandleFuncEx("./assets/templates/result.html", cfg.ResultConfig, DrawingResult))
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
