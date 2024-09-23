package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
			fmt.Printf("error: %g\n", err)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			fmt.Printf("error: %g\n", err)
		}
	}
}

type Person struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

func (p Person) Key() int {
	return p.Number
}

func genDrawingHandleFunc(cfg config.Config, drawingData *drawing.Data[Person, int], list *List) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		number, err := strconv.Atoi(r.PathValue("number"))
		if err != nil {
			return
		}
		result, err := drawingData.Draw(Person{Number: number})
		if err != nil {
			genHandleFunc("./assets/templates/error.html", err.Error())(w, r)
			return
		}
		fmt.Printf("number: %v, result: %+v, waiting: %v\n", number, result, drawingData.WaitingFreshmenCount())

		SaveResults(drawingData.Results(), list)

		names := make([]string, len(result))
		for i, v := range result {
			names[i] = v.Name
		}
		// names := slices.Collect(slices.Values(result).Map(func(p Person) string { return p.Name }))

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
	Freshmen []Person `json:"freshmen"`
	Seniors  []Person `json:"seniors"`
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
		list.Freshmen,
		list.Seniors,
	)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/static"))))
	http.HandleFunc("/{$}", genHandleFunc("./assets/templates/page.html", cfg.HomepageConfig))
	http.HandleFunc("/drawing", genHandleFunc("./assets/templates/page.html", cfg.DrawingConfig))
	http.HandleFunc("/settings", genHandleFunc("./assets/templates/settings.html", struct {
		DrawingData *drawing.Data[Person, int]
		List        List
	}{
		&drawingData,
		list,
	}))
	http.HandleFunc("/result/{number}", genDrawingHandleFunc(cfg.ResultConfig, &drawingData, &list))
	http.HandleFunc("/results.json", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./results.json") })
	http.Handle("/", http.NotFoundHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func LoadList(filepath string) (List, error) {
	listFile, err := os.ReadFile(filepath)
	if err != nil {
		return List{}, err
	}

	var list List
	err = json.Unmarshal(listFile, &list)

	return list, err
}

type result struct {
	Number int      `json:"number"`
	Name   string   `json:"name"`
	Paired []Person `json:"paired"`
}

func SaveResults(results drawing.Results[Person, int], list *List) error {
	rs := make([]result, 0, results.Len())
	for _, v := range list.Freshmen {
		if !results.Contains(v) {
			continue
		}
		rs = append(rs, result{Number: v.Number, Name: v.Name, Paired: results.Index(v)})
	}
	slices.SortFunc(rs, func(a result, b result) int {
		return a.Number - b.Number
	})
	b, err := json.MarshalIndent(rs, "", "	")
	if err != nil {
		return err
	}

	err = os.WriteFile("./results.json", b, 0644)
	if err != nil {
		return err
	}

	return nil
}
