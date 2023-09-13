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
	MakeFolder := func(w http.ResponseWriter, r *http.Request) {
		type RequestBody struct {
			ParentDirectory string `json:"parentDirectory"`
		}
		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		// drive.

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
		defer file.Close()

		name := r.FormValue("filename")
		targetBranch := r.FormValue("targetBranch")

		err = drive.UploadFile(name, targetBranch, file)
		if err != nil {
			fmt.Println("Error creating file:", err)
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		branch, err := drive.Tree.FindBranchDescending(targetBranch)
		drive.List(path)

	}

	r := mux.NewRouter()
	port := ":8080"

	r.HandleFunc("/", GetDrive).Methods("GET")
	r.HandleFunc("/drive/{path:.*}", WriteDrive).Methods("POST")

	http.ListenAndServe(port, r)

}
