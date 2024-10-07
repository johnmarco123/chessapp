package main

import (
	"encoding/json";
	"errors";
	"fmt";
	"net/http";
	"os";
)

func logrequest(r *http.Request) {
	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path);
}

func getroot(w http.ResponseWriter, r *http.Request) {
	logrequest(r);
	http.FileServer(http.Dir("./")).ServeHTTP(w, r);
}

type Response struct {
	Board [8][8]piece `json:"board"`;
	Msg string        `json:"msg"`;
	Err string        `json:"err"`;
}

func handleservices(w http.ResponseWriter, r *http.Request) {
    // Set the response header to application/json
    w.Header().Set("Content-Type", "application/json");

    // Create a response object
	params := r.URL.Query();
	down := params.Get("down"); // Get param2 value
    up := params.Get("up"); // Get param1 value
	board, msg, err := boardlogic(down, up);
    response := Response {
		Board: board,
		Msg: msg,
		Err: err,
    }

    // Encode the response as JSON and write it to the response writer
    json.NewEncoder(w).Encode(response);
}


func getboard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json");
    response := Response {
		Board: board,
		Msg: "",
		Err: "",
    }
    json.NewEncoder(w).Encode(response);
}

func main() {
	http.HandleFunc("/", getroot);
	http.HandleFunc("/services/", handleservices);
	http.HandleFunc("/getboard/", getboard);

	err := http.ListenAndServe(":3333", nil);
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n");
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err);
		os.Exit(1);
	}
}

