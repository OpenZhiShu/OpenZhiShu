package main

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "OpenZhiShu/pkg/drawing"
)

type Element interface {
	HTML() template.HTML
}

type HomepageConfig struct {
	Background Element
	Elements   []Element
}

type AllConfig struct {
	homepageConfig HomepageConfig
}

func handleHomepage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/templates/homepage.html")
	if err != nil {
		fmt.Printf("eror: %g\n", err)
		return
	}

	err = t.Execute(w, struct{}{})
	if err != nil {
		fmt.Printf("eror: %g\n", err)
	}
}

func handleDrawing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "the drawing page")
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/drawing", handleDrawing)
	http.HandleFunc("/{$}", handleHomepage)
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
