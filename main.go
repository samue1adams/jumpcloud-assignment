package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Password struct {
	Password string
	id       int
}

type Stats struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

var inc = 0

var durations []int64
var passwordsMap = make(map[int]string)
var wg sync.WaitGroup

func acceptPassword(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.ParseForm()
	req := r.FormValue("password")

	password := Password{
		Password: req,
		id:       inc + 1}
	response := json.NewEncoder(w).Encode(password.id)
	if response != nil {
		return
	}

	inc += 1

	wg.Add(1)
	go hashPassword(password, &wg, passwordsMap)

	end := time.Now()
	duration := end.Sub(start)
	durations = append(durations, duration.Microseconds())
}

func hashPassword(p Password, wg *sync.WaitGroup, m map[int]string) {
	hash := hmac.New(sha512.New, []byte(p.Password))
	encodedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	m[p.id] = encodedPassword
	time.Sleep(5 * time.Second)
	wg.Done()
}

func getHashedPasswordById(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wg.Wait()

	vars := mux.Vars(r)
	key := vars["id"]
	keyInt, _ := strconv.ParseInt(key, 0, 64)
	fmt.Fprintf(w, passwordsMap[int(keyInt)])

	end := time.Now()
	duration := end.Sub(start)
	durations = append(durations, duration.Microseconds())
}

func stats(w http.ResponseWriter, r *http.Request) {
	stats := Stats{
		Total:   len(passwordsMap),
		Average: getAverage(durations)}
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		return
	}
}

func getAverage(intArray []int64) float64 {
	var sum int64 = 0
	size := len(intArray)
	for i := 0; i < size; i++ {
		sum += intArray[i]
	}
	return float64(sum) / (float64(size))
}

func shutDown(w http.ResponseWriter, r *http.Request) {
	//wg.Wait()

	fmt.Fprintf(w, "couldn't figure this one out")
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hash", acceptPassword).Methods("POST")
	router.HandleFunc("/hash/{id}", getHashedPasswordById).Methods("GET")
	router.HandleFunc("/stats", stats).Methods("GET")
	router.HandleFunc("/shutdown", shutDown).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequest()
}
