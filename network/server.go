package network

import (
	"encoding/json"
	"fmt"
	"my-go-project/filesystem"
	"net/http"

	"github.com/gorilla/mux"
)

func Server(drive filesystem.Drive) {
	GetDrive := func(w http.ResponseWriter, r *http.Request) {
		dirEntries, _ := drive.List("/")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dirEntries)
	}
	WriteDrive := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		if drive.Exists(path) {
			w.Header().Set("Status", "404")
			w.Write([]byte("Path: " + path + ", does not exist"))
			return
		}
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error retrieving file:", err)
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}

	}

	r := mux.NewRouter()
	port := ":8080"

	r.HandleFunc("/", GetDrive).Methods("GET")
	r.HandleFunc("/drive/{path:.*}", WriteDrive).Methods("POST")

	http.ListenAndServe(port, r)

}
