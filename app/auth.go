package main

import (
	"fmt";
	"log";
	"golang.org/x/crypto/bcrypt";
)

func hashpassword(password string) string {
//	start := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13);
//	end := time.Now()
// fmt.Println(end - start);
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func register(username string, password string) {
	query := "select username from users where username = ?";
	var result string;
	db.QueryRow(query, username).Scan(&result)
	if result != "" {
		fmt.Println("That username is taken!");
		return;
	}
	hash := hashpassword(password);
	query = "insert into users (username, password) values (?, ?)";
	db.QueryRow(query, username, hash)
	fmt.Println("User: " + username + " successfully added to database!");
}

func login(username string, password string) {
	query := "select password from users where username = ?";
	var hash string;
	err := db.QueryRow(query, username).Scan(&hash)
	if err != nil {
		fmt.Println("Attempted login for user " + username + " but that user does not exist in db");
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Unsuccessful login for user " + username);
		return;
	}
	fmt.Println("Successful login for user " + username);
}
