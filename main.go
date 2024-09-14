package main

import (
	"fmt"
	// "html/template"
	"log"
	"net/http"
	"os"
	// "OpenZhiShu/pkg/drawing"
)

func handleHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "the home page")
}

func handleDrawing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "the drawing page")
}

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Printf("need a argument for address\n")
		return
	case 2:
		fmt.Printf("http://localhost:%v", os.Args[1])
	default:
		fmt.Printf("too many arguments: %v\n", len(os.Args))
		return
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/drawing", handleDrawing)
	http.HandleFunc("/{$}", handleHomepage)
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Args[1]), nil))
}
