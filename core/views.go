package core

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"go-abbreviation/templates"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func MakeHttp(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("HTTP handler error", "error", err, "path", r.URL.Path)
		}
	}
}

func TemplRender(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}

// No longer maintained, going to remove this
func ShowListJson(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// keyword := r.Form.Get("keyword")
	// fmt.Println("THE FORM PARAM IS: ", keyword)
	// jsonFile, err := os.Open("static/data.json")
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// fmt.Println("Successfully opened json file")
	// defer jsonFile.Close()

	// // Decode json file
	// var abv models.List
	// json.NewDecoder(jsonFile).Decode(&abv)

	// // Render Templ template
	// TemplRender(w, r, templates.List("Abbreviations", abv, keyword))
}

func ShowListDb(w http.ResponseWriter, r *http.Request) {
	// keyword := chi.URLParam(r, "keyword") //Only for file routing params
	keyword := r.URL.Query().Get("q")
	// fmt.Println(r.Header.Get("Hx-Request"))
	results := DbRetrieveKeywordJaroWinker(keyword)
	if r.Header.Get("Hx-Request") == "true" {
		fmt.Println("Hx request received, using filtered results")
		TemplRender(w, r, templates.ListFilteredResult(results, strconv.Itoa(len(results))))
	} else {
		fmt.Println("Doing normal list")
		TemplRender(w, r, templates.ListAll("Abbreviations", results, strconv.Itoa(len(results)), keyword))
	}
}

//
// Note: No need for this because I changed to GET
//
// func ShowListDbFilter(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	keyword := r.Form.Get("keyword")
// 	results := DbRetrieveKeywordJaroWinker(keyword)
// 	TemplRender(w, r, templates.ListFilteredResult(results, strconv.Itoa(len(results))))
// }

func ShowListDbAlphabets(w http.ResponseWriter, r *http.Request) {
	results := DbRetrieveAlphabet("")
	if r.Header.Get("Hx-Request") == "true" {
		TemplRender(w, r, templates.ListFilteredResult(results, strconv.Itoa(len(results))))
	} else {
		TemplRender(w, r, templates.ListAll("Abbreviations", results, strconv.Itoa(len(results)), ""))
	}
}

func ShowListDbFilterAlphabets(w http.ResponseWriter, r *http.Request) {
	// err := r.ParseForm()
	// if err != nil {
	// 	fmt.Println("Error parsing form: ", err)
	// }
	// alphabet := r.PostFormValue("alphabet")
	alphabet := chi.URLParam(r, "alphabet")

	results := DbRetrieveAlphabet(alphabet)
	if r.Header.Get("Hx-Request") == "true" {
		TemplRender(w, r, templates.ListFilteredResult(results, strconv.Itoa(len(results))))
	} else {
		fmt.Println("Doing normal list")
		TemplRender(w, r, templates.ListAll("Abbreviations", results, strconv.Itoa(len(results)), ""))
	}
}

func SyncJsonToDb(w http.ResponseWriter, r *http.Request) {
	// Function to populate using json, I find it easier to work with
	// DbPopulateFromJson()
}

func Test(w http.ResponseWriter, r *http.Request) {
	DbBase()

	// e, _ := beidermorse.NewEncoder()
	// primary := e.Encode("centre")
	// dbvalue := e.Encode("centre")
	// dbvalue2 := e.Encode("Economist Intelligence Unit- research organization")
	// dbvalue3 := e.Encode("senter")
	// fmt.Println(dbvalue2)
	// // longArr := []string{"centre", "cen", "centerrr"}
	// // var results bool
	// // var results2 bool

	// for _, v := range dbvalue {
	// 	result := slices.Contains(primary, v)
	// 	fmt.Println("dbvalue: ", result)
	// 	if result {
	// 		break
	// 	}
	// }
	// for _, v := range dbvalue2 {
	// 	result2 := slices.Contains(primary, v)
	// 	// fmt.Println("dbvalue2: ", result2)
	// 	if result2 {
	// 		break
	// 	}
	// }
	// for _, v := range dbvalue3 {
	// 	result3 := slices.Contains(primary, v)
	// 	fmt.Println("dbvalue3: ", result3)
	// 	if result3 {
	// 		break
	// 	}
	// }
	// for _, v := range dbvalue2 {
	// 	fmt.Println("dbvalue2:", slices.Contains(primary, v))
	// }
	// for _, v := range dbvalue3 {
	// 	fmt.Println("dbvalue3:", slices.Contains(primary, v))
	// }

	// fmt.Println("primary: ", primary)
	// fmt.Println("dbvalue: ", dbvalue)
	// fmt.Println("dbvalue2: ", dbvalue2)
	// fmt.Println("dbvalue3: ", dbvalue3)

	// fmt.Println(emptySlice)
	w.Write([]byte("Testing"))
}
