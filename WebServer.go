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
	r.HandleFunc("/getsubset/{seed}/{startyear}/{endyear}", func(w http.ResponseWriter, r *http.Request) {
		// Function Here
		vars := mux.Vars(r)
		seed := vars["seed"]
		startyear := vars["startyear"]
		endyear := vars["endyear"]
		seed = string(seed)
		n, err := strconv.ParseInt(seed, 10, 64)
		s, err := strconv.ParseInt(startyear, 10, 64)
		e, err := strconv.ParseInt(endyear, 10, 64)
		fmt.Println("Seed", n)
		if err == nil {
			data := Gen.Gen(n, s, e)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		}
		
	}).Host("localhost:42069")

	r.HandleFunc("/getnewsubset/{startyear}/{endyear}", func(w http.ResponseWriter, r *http.Request) {
		// Function Here
		fmt.Println("No Seed")
		vars := mux.Vars(r)
		startyear := vars["startyear"]
		endyear := vars["endyear"]
		s, err := strconv.ParseInt(startyear, 10, 64)
		e, err := strconv.ParseInt(endyear, 10, 64)
		if err == nil {
			data := Gen.Gen(int64(time.Now().UnixNano()), s, e)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		}
	}).Host("localhost:42069")

	http.ListenAndServe(":42069", r)
	
}