package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

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
	Category        string
	Assumption      string
}
type Situation struct {
	Name        string
	Loses       string
	Probability string
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

func MainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/main_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Index1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
func Index2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func MakeSlice(n int) []struct{} {
	return make([]struct{}, n)
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	numAssets := r.FormValue("numAssets")
	numAssetsInt := 0
	if numAssets != "" {
		fmt.Sscanf(numAssets, "%d", &numAssetsInt)
	}

	tmpl, err := template.New("asset_entry_form.html").Funcs(template.FuncMap{"makeSlice": MakeSlice}).ParseFiles("templates/asset_entry_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, numAssetsInt)
}

func SituationsHandler(w http.ResponseWriter, r *http.Request) {
	numSituations := r.FormValue("numSituation")
	numSituationsInt := 0
	if numSituations != "" {
		fmt.Sscanf(numSituations, "%d", &numSituationsInt)
	}

	tmpl, err := template.New("situation_entry_form.html").Funcs(template.FuncMap{"makeSlice": MakeSlice}).ParseFiles("templates/situation_entry_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, numSituationsInt)
}

func SaveAssetsHandler(w http.ResponseWriter, r *http.Request) {
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

		confInt, err := strconv.Atoi(confidentiality)
		if err != nil {
			fmt.Println("can't convert confidentiality")
			log.Fatal(err)
		}
		integrInt, err := strconv.Atoi(integrity)
		if err != nil {
			fmt.Println("can't convert integrity")
			log.Fatal(err)
		}
		availInt, err := strconv.Atoi(availability)
		if err != nil {
			fmt.Println("can't convert availability")
			log.Fatal(err)
		}
		weightInt, err := strconv.Atoi(weight)
		if err != nil {
			fmt.Println("can't convert weight")
			log.Fatal(err)
		}

		totalAsset := confInt + integrInt + availInt
		tav := totalAsset * weightInt
		category := "0"
		assumption := "no assumption"
		if tav >= 20 && tav <= 27 {
			category = "1"
			assumption = "require very serious and more attention"
		} else if tav >= 12 && tav <= 18 {
			category = "2"
			assumption = "require serious attention"
		} else if tav <= 10 {
			category = "3"
			assumption = "require less attention"
		}

		asset := Asset{
			Name:            name,
			Confidentiality: confidentiality,
			Integrity:       integrity,
			Availability:    availability,
			TotalAsset:      strconv.Itoa(totalAsset),
			WeightOfAsset:   weight,
			TAV:             strconv.Itoa(tav),
			Category:        category,
			Assumption:      assumption,
		}

		assets = append(assets, asset)
	}

	// Log assets before saving
	log.Println("Assets before saving:", assets)

	// Save assets to the database
	SaveAssetsToDatabase(assets)

	// Log assets after saving
	log.Println("Assets after saving:", assets)

	// Display the saved assets
	DisplayAssets(w, nil, assets)
}

func SaveSituationsHandler(w http.ResponseWriter, r *http.Request) {
	numSituations := r.FormValue("numSituation")

	numSituationsInt := 0

	fmt.Sscanf(numSituations, "%d", &numSituationsInt)

	var situations []Situation
	var totalSituationValue int
	for i := 0; i < numSituationsInt; i++ {
		name := r.FormValue(fmt.Sprintf("name_%d", i))
		loses := r.FormValue(fmt.Sprintf("loses_%d", i))
		probability := r.FormValue(fmt.Sprintf("probability_%d", i))

		losesInt, err := strconv.Atoi(loses)
		if err != nil {
			fmt.Println("can't convert loses to integer")
			log.Fatal(err)
		}

		probabilityInt, err := strconv.Atoi(probability)
		if err != nil {
			fmt.Println("can't convert probability to integer")
			log.Fatal(err)
		}

		product := losesInt * probabilityInt
		totalSituationValue += product
		situation := Situation{
			Name:        name,
			Loses:       loses,
			Probability: probability,
		}


		situations = append(situations, situation)
	}
	tmpl, err := template.ParseFiles("templates/all_situations.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	data := struct {
        Situations            []Situation
        TotalSituationValue   int
    }{
        Situations:            situations,
        TotalSituationValue:   totalSituationValue,
    }

    tmpl.Execute(w, data)
	

	
	//SaveAssetsToDatabase(situations)



	
//	DisplayAssets(w, nil, situations)
}

func SaveAssetsToDatabase(assets []Asset) {
	for _, asset := range assets {
		_, err := db.Exec(
			"INSERT INTO assets_risk (name, confidentiality, integrity, availability, total_asset, weight_of_asset, tav) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			asset.Name, asset.Confidentiality, asset.Integrity, asset.Availability, asset.TotalAsset, asset.WeightOfAsset, asset.TAV,
		)
		if err != nil {
			log.Println("Error inserting asset:", err)
		} else {
			log.Println("Asset inserted successfully:", asset)
		}
	}
}

func DisplayAssets(w http.ResponseWriter, r *http.Request, assets []Asset) {
	log.Println("Received assets:", assets)

	tmpl, err := template.ParseFiles("templates/all_assets.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, assets)
}
