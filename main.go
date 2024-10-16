package main

import (
	"encoding/csv"
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
		result, err := drawingData.Draw(number)
		if err != nil {
			genHandleFunc("./assets/templates/error.html", err.Error())(w, r)
			return
		}
		fmt.Printf("number: %v, result: %+v, waiting: %v\n", number, result, drawingData.WaitingFreshmenCount())

		SaveResults(drawingData.Results(), list)

		names := slices.Collect(func(yield func(string) bool) {
			for _, v := range result {
				if !yield(v.Name) {
					return
				}
			}
		})
		// names := slices.Collect(slices.Values(result).Map(func(p Person) string { return p.Name }))

		variables := map[string]string{
			"result":        strings.Join(names, " & "),
			"result_number": strconv.Itoa(result[0].Number),
		}

		elems := slices.Clone(cfg.Elements)
		for i := range elems {
			if elems[i].Type != "variable" {
				continue
			}
			elems[i].Type = elems[i].Other["to_type"].(string)

			value, inMap := variables[elems[i].Content]
			if !inMap {
				elems[i].Type = "text"
				elems[i].Content = fmt.Sprintf("no variable `%v`", elems[i].Content)
				continue
			}
			prefix, _ := elems[i].Other["prefix"].(string)
			suffix, _ := elems[i].Other["suffix"].(string)
			elems[i].Content = prefix + value + suffix
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

	cfg, err := config.LoadConfig[Config]("./configs/config.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	list, err := LoadList("./configs/list.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	drawingData := drawing.MakeData(
		list.Freshmen,
		list.Seniors,
	)

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./assets/style"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./configs/static"))))
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
	http.HandleFunc("/results/results.json", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./results.json") })
	http.HandleFunc("/results/results.csv", func(w http.ResponseWriter, r *http.Request) {
		results := drawingData.Results()
		rs := make([][]string, 0, len(results))
		for _, v := range list.Freshmen {
			paired, inMap := results[v.Key()]
			if !inMap {
				continue
			}
			rs = append(rs, append(
				[]string{strconv.Itoa(v.Number), v.Name},
				slices.Collect(func(yield func(string) bool) {
					for _, p := range paired {
						if !yield(strconv.Itoa(p.Number)) {
							return
						}
						if !yield(p.Name) {
							return
						}
					}
				})...,
			))
		}
		slices.SortFunc(rs, func(a []string, b []string) int {
			aAsInt, _ := strconv.Atoi(a[0])
			bAsInt, _ := strconv.Atoi(b[0])
			return aAsInt - bAsInt
		})
		csv.NewWriter(w).WriteAll(rs)
	})
	http.HandleFunc("PUT /drawing/all", func(w http.ResponseWriter, r *http.Request) {
		drawingData.DrawAll()
		SaveResults(drawingData.Results(), &list)
		w.Header().Add("HX-Refresh", "true")
	})
	http.HandleFunc("DELETE /results/delete", func(w http.ResponseWriter, r *http.Request) {
		drawingData.Reset()
		w.Header().Add("HX-Refresh", "true")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		genHandleFunc("./assets/templates/error.html", fmt.Sprintf("Sorry, the page `%v` could not be found.", r.URL.Path))(w, r)
	})

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
	Person
	Paired []Person `json:"paired"`
}

func SaveResults(results drawing.Results[Person, int], list *List) error {
	rs := make([]result, 0, len(results))
	for _, v := range list.Freshmen {
		paired, inMap := results[v.Key()]
		if !inMap {
			continue
		}
		rs = append(rs, result{Person: Person{Number: v.Number, Name: v.Name}, Paired: paired})
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
