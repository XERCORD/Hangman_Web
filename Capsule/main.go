package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type PageConditionSimple struct {
	Check      bool
	CheckOwner bool
}

func main() {
	temp, err := template.ParseGlob("./template/*.html")
	if err != nil {
		fmt.Println("ERREUR")
		os.Exit(02)
		return
	}
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello !")
	})

	http.HandleFunc("/condition", func(w http.ResponseWriter, r *http.Request) {
		dataPage := PageConditionSimple{true, false}
		temp.ExecuteTemplate(w, "exempleConditionSimple", dataPage)
	})

	http.ListenAndServe("localhost:8000", nil)

}
