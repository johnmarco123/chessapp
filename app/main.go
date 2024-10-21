package main

import (
	"log";
	"database/sql";
	"github.com/go-sql-driver/mysql";
	"encoding/json";
	"errors";
	"fmt";
	"net/http";
	"os";
)

var db *sql.DB;


func logrequest(r *http.Request) {
	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path);
}

func getroot(w http.ResponseWriter, r *http.Request) {
	logrequest(r);
}

type Response struct {
	Board [8][8]piece `json:"board"`;
	Msg string        `json:"msg"`;
	Err string        `json:"err"`;
}

// Define a struct for the request body
type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func registerhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req AuthRequest;
		err := json.NewDecoder(r.Body).Decode(&req);
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest);
		}
		register(req.Username, req.Password);
		return;
	}
	logrequest(r);
	http.ServeFile(w, r, "./static/register.html");
}
func loginhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req AuthRequest;
		err := json.NewDecoder(r.Body).Decode(&req);
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest);
		}
		login(req.Username, req.Password);
		return;
	}
	logrequest(r);
	http.ServeFile(w, r, "./static/login.html");
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

    // Capture connection properties.
    cfg := mysql.Config{
        User:   "root",
		// set your db pass in ur environment with this: 
		// export DB_PASS="your_password"
		Passwd: os.Getenv("DB_PASS"), 
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "chessapp",
    }
    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("DB Connected!")




	fs := http.FileServer(http.Dir("./static"));
	http.Handle("/", fs);

	http.HandleFunc("/login/", loginhandler);
	http.HandleFunc("/register/", registerhandler);
	http.HandleFunc("/services/", handleservices);
	http.HandleFunc("/getboard/", getboard);


	fmt.Println("Server started!");
	err = http.ListenAndServe(":3333", nil);



	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n");
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err);
		os.Exit(1);
	}
}

