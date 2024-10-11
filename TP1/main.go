package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
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

type User struct {
	Nom           string
	Prénom        string
	DateNaissance string
	Sexe          string
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

var (
	nomRegex    = regexp.MustCompile(`^[A-Za-zÀ-ÿ]{1,32}$`)
	prenomRegex = regexp.MustCompile(`^[A-Za-zÀ-ÿ]{1,32}$`)
	sexeOptions = []string{"masculin", "féminin", "autre"}
)

func main() {
	// Charger tous les templates HTML dans le dossier "template"
	temp, err := template.ParseGlob("./template/*.html")
	if err != nil {
		fmt.Println(fmt.Sprintf("ERREUR => %s", err.Error()))
		os.Exit(2)
	}

	// Gérer les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("asset"))))

	// Route /cours
	http.HandleFunc("/cours", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bonjour, bienvenue sur la route /cours"))
	})

	// Route /promo
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

	// Route /change avec un compteur de vues
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

	// Route /user/form pour le formulaire utilisateur
	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nom := r.FormValue("nom")
			prenom := r.FormValue("prenom")
			dateNaissance := r.FormValue("dateNaissance")
			sexe := r.FormValue("sexe")

			fmt.Println("Date de naissance reçue:", dateNaissance)

			if !nomRegex.MatchString(nom) || !prenomRegex.MatchString(prenom) || !isValidSexe(sexe) {
				http.Error(w, "Informations invalides", http.StatusBadRequest)
				return
			}

			// Traitez les données ici (enregistrement, affichage, etc.)
			http.Redirect(w, r, "/success", http.StatusSeeOther)
			return
		}

		err := temp.ExecuteTemplate(w, "form", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route pour afficher la page de succès
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		err := temp.ExecuteTemplate(w, "Succes", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Démarrage du serveur
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		fmt.Println(fmt.Sprintf("ERREUR => %s", err.Error()))
		os.Exit(1)
	}
}

func isValidSexe(sexe string) bool {
	for _, option := range sexeOptions {
		if sexe == option {
			return true
		}
	}
	return false
}
