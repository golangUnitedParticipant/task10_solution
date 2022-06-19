package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func handleName(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["name"]
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello, %s!", p)))
}

func handleBad(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("This call will always fail"))
}

func handleDataPost(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad body"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("I got message:\n%s", string(b))))
}

func handleDataGet(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	if len(b) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("I don't need body for GET request"))
		return
	}
	w.Write([]byte("Good! no body sent"))
}

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.Atoi(r.Header.Get("a"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	b, err := strconv.Atoi(r.Header.Get("b"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Add("a+b", strconv.FormatInt(int64(a+b), 10))
	w.Write([]byte("Good! no body sent"))
}

func Start(host string, port int) {
	router := mux.NewRouter()
	router.HandleFunc("/name/{name}", handleName).Methods(http.MethodGet)
	router.HandleFunc("/bad", handleBad).Methods(http.MethodGet)
	router.HandleFunc("/data", handleDataPost).Methods(http.MethodPost)
	router.HandleFunc("/headers", handleHeaders).Methods(http.MethodPost)
	router.HandleFunc("/data", handleDataGet).Methods(http.MethodGet)
	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}
}

func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
