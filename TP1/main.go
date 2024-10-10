package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

type Etudiant struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string
}

type ListeEtudiant struct {
	Nom_de_Classe string
	Filiere       string
	Niveau        string
	Nbr_Etudiant  int
	Etudiant      []Etudiant
}

var (
	// Compteur de vues
	viewCount int
	mutex     sync.Mutex // Pour gérer les accès concurrents au compteur
)

type PageData struct {
	Count  int
	IsEven bool
}

func main() {
	temp, err := template.ParseGlob("./template/*.html")
	if err != nil {
		fmt.Println(fmt.Sprintf("ERREUR => %s", err.Error()))
		os.Exit(2)
	}
	// Chalenge 1
	// Gérer les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("asset"))))

	http.HandleFunc("/cours", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bonjour bienvenue sur la route /cours"))
	})

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		classe := ListeEtudiant{
			Nom_de_Classe: "B1 Informatique",
			Filiere:       "Informatique",
			Niveau:        "Bachelor 1",
			Nbr_Etudiant:  7,
			Etudiant: []Etudiant{
				{"纪", "建锋", 18, "M"},
				{"Skibidi", "Dimitri", 17, "F"},
				{"Amir", "Eddy", 15, "M"},
				{"M", "Matheo", 9, "F"},
				{"Thibaut", "Eddy", 20, "M"},
				{"Daniel", "R", 2, "F"},
				{"Panda", "Etienne", 20, "M"},
			},
		}
		err := temp.ExecuteTemplate(w, "AffichageDonnes", classe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// Chalenge 2

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock() // Verrouille l'accès au compteur
		viewCount++
		isEven := viewCount%2 == 0
		mutex.Unlock() // Déverrouille l'accès au compteur

		data := PageData{
			Count:  viewCount,
			IsEven: isEven,
		}
		err := temp.ExecuteTemplate(w, "Condition", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	// Ajout de la gestion d'erreur pour ListenAndServe
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println(fmt.Sprintf("ERREUR => %s", err.Error()))
		os.Exit(1)
	}
}
