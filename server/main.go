package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func logrequest(r *http.Request) {
	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
}

func getroot(w http.ResponseWriter, r *http.Request) {
	logrequest(r)
	http.FileServer(http.Dir("./")).ServeHTTP(w, r)
}

type Response struct {
    Down string
	Up string
	Msg string
}

func handleservices(w http.ResponseWriter, r *http.Request) {
    // Set the response header to application/json
    w.Header().Set("Content-Type", "application/json")

    // Create a response object
	params := r.URL.Query()
    up := params.Get("up") // Get param1 value
    down := params.Get("down") // Get param2 value
	downafter, upafter, msg := boardlogic(down, up);
    response := Response {
		Down: downafter,
        Up: upafter,
		Msg: msg,
    }

    // Encode the response as JSON and write it to the response writer
    json.NewEncoder(w).Encode(response)
}


func getboard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(board)
}

func main() {
	http.HandleFunc("/", getroot)
	http.HandleFunc("/services/", handleservices)
	http.HandleFunc("/getboard/", getboard)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

