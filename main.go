package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Asset struct {
	Name            string
	Confidentiality string
	Integrity       string
	Availability    string
	TotalAsset      string
	WeightOfAsset   string
	TAV             string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/risk?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/assets", assetsHandler).Methods("POST")
	r.HandleFunc("/save_assets", saveAssetsHandler).Methods("POST")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func makeSlice(n int) []struct{} {
	return make([]struct{}, n)
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	numAssets := r.FormValue("numAssets")
	numAssetsInt := 0
	if numAssets != "" {
		fmt.Sscanf(numAssets, "%d", &numAssetsInt)
	}

	tmpl, err := template.New("asset_entry_form.html").Funcs(template.FuncMap{"makeSlice": makeSlice}).ParseFiles("templates/asset_entry_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, numAssetsInt)
}
func saveAssetsHandler(w http.ResponseWriter, r *http.Request) {
	numAssets := r.FormValue("numAssets")
	numAssetsInt := 0
	fmt.Sscanf(numAssets, "%d", &numAssetsInt)

	var assets []Asset

	for i := 0; i < numAssetsInt; i++ {
		name := r.FormValue(fmt.Sprintf("name_%d", i))
		confidentiality := r.FormValue(fmt.Sprintf("confidentiality_%d", i))
		integrity := r.FormValue(fmt.Sprintf("integrity_%d", i))
		availability := r.FormValue(fmt.Sprintf("availability_%d", i))
		weight := r.FormValue(fmt.Sprintf("weight_%d", i))
		
		totalAsset := confidentiality + integrity + availability
		totalAssetInt, _ := strconv.Atoi(totalAsset)
		weightInt, _ := strconv.Atoi(weight)
		tav := strconv.Itoa(totalAssetInt * weightInt)
		asset := Asset{
			Name:            name,
			Confidentiality: confidentiality,
			Integrity:       integrity,
			Availability:    availability,
			TotalAsset:      totalAsset,
			WeightOfAsset:   weight,
			TAV:             tav,
		}

		assets = append(assets, asset)
	}

	// Log assets before saving
	log.Println("Assets before saving:", assets)

	// Save assets to the database
	saveAssetsToDatabase(assets)

	// Log assets after saving
	log.Println("Assets after saving:", assets)

	// Display the saved assets
	displayAssets(w, assets)
}

func saveAssetsToDatabase(assets []Asset) {
	for _, asset := range assets {
		_, err := db.Exec(
			"INSERT INTO assets_risk (data_type, confidentiality, integrity, availability, total_asset, weight_of_asset, tav) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			asset.Name, asset.Confidentiality, asset.Integrity, asset.Availability, asset.TotalAsset, asset.WeightOfAsset, asset.TAV,
		)
		if err != nil {
			log.Println("Error inserting asset:", err)
		} else {
			log.Println("Asset inserted successfully:", asset)
		}
	}
}

func displayAssets(w http.ResponseWriter, assets []Asset) {
	log.Println("Received assets:", assets)

	tmpl, err := template.ParseFiles("templates/all_assets.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, assets)
}
