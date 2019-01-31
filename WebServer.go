package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	Gen "Ling499/packageGen"
)

func main() {
	fmt.Println("WebServer")
	r := mux.NewRouter()
	r.HandleFunc("/getsubset/{seed}", func(w http.ResponseWriter, r *http.Request) {
		// Function Here
		vars := mux.Vars(r)
		seed := vars["seed"]
		seed = string(seed)
		n, err := strconv.ParseInt(seed, 10, 64)
		fmt.Println("Seed", n)
		if err == nil {
			data := Gen.Gen(n)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		}
		
	}).Host("localhost:42069")

	r.HandleFunc("/getsubset", func(w http.ResponseWriter, r *http.Request) {
		// Function Here
		fmt.Println("No Seed")
		data := Gen.Gen(int64(time.Now().UnixNano()))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}).Host("localhost:42069")

	http.ListenAndServe(":42069", r)
	
}