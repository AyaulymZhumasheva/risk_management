package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type DataType struct {
	Name  string
	Value int
}

func assetValuationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("templates/index.html")
		if err != nil {
			return
		}
		temp.Execute(w, nil)
	} else if r.Method == "POST" {

		r.ParseForm()
		dataType := r.Form.Get("datatype")
		conf, _ := strconv.Atoi(r.Form.Get("conf"))
		integrity,_ := strconv.Atoi(r.Form.Get("integrity"))
		availability,_ := strconv.Atoi(r.Form.Get("availability"))
		weightAsset,_ := strconv.Atoi(r.Form.Get("weightAsset"))

		fmt.Fprintf(w, "Значение вашего актива: %s", dataType)
		fmt.Fprintf(w, "Значение вашего актива: %d", conf)
		fmt.Fprintf(w, "Значение вашего актива: %d", integrity)
		fmt.Fprintf(w, "Значение вашего актива: %d", availability)
		fmt.Fprintf(w, "Значение вашего актива: %d", weightAsset)

		TA := conf + availability + integrity
		TAV := TA * weightAsset
		fmt.Fprintf(w, "Значение вашего актива: %d", TAV)
	}
}
func main() {
	http.HandleFunc("/", assetValuationHandler)
	fmt.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", nil)
}
