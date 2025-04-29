package main

import (
	"fmt";
	"log";
	"strconv";
	"net/http";
	"time";
	"golang.org/x/crypto/bcrypt";
)

func hash(str string) string {
//	start := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(str), 13);
//	end := time.Now()
// fmt.Println(end - start);
	if err != nil {
		log.Fatal(err);
	}
	return string(hash);
}

func register(w http.ResponseWriter, r *http.Request,username string, password string) {
	query := "select username from users where username = ?";
	var result string;
	db.QueryRow(query, username).Scan(&result);
	if result != "" {
		fmt.Println("That username is taken!");
		return;
	}
	hash := hash(password);
	query = "insert into users (username, password) values (?, ?)";
	db.QueryRow(query, username, hash);
	fmt.Println("User: " + username + " successfully added to database!");
	http.Redirect(w, r, "/game/", http.StatusSeeOther)

}

func login(w http.ResponseWriter, r *http.Request, username string, password string) {
	// if the user has a cookie and it is valid, simply log them in. Otherwise
	// generate a new cookie if they login correctly
	if (!validatecookie(w, r)) {
		query := "select userid, password from users where username = ?";
		var hash string;
		var userid int;
		err := db.QueryRow(query, username).Scan(&userid, &hash);
		if err != nil {
			fmt.Println("Attempted login for user " + username + " but that user does not exist in db");
			return;
		}
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password));
		if err != nil {
			fmt.Println("Unsuccessful login for user " + username);
			return;
		}
		// successful login
		cookie := generatecookie(w, r, username);

		http.SetCookie(w, cookie);
		fmt.Println("login: setting new cookie for user ======>" + cookie.Value);

		currenttime := time.Now().UTC();

		// delete any old cookies before adding a new one

		var id int;
		query = "select id from sessions where userid = ?";
		db.QueryRow(query, userid).Scan(&id);
		// if a session already exists we update it, otherwise we create a new
		// one
		if (id > 0) { 
			query = "update sessions set sessionid = ?, path = ?, expires = ?, createdat = ?, updatedat = ? where id = ?";
			db.QueryRow(query, cookie.Value, cookie.Path, cookie.Expires, currenttime, currenttime, id);
		} else {
			query = "insert into sessions (userid,sessionid,path,expires,createdat,updatedat) values (?,?,?,?,?,?)";
			db.QueryRow(query, userid, cookie.Value, cookie.Path, cookie.Expires, currenttime, currenttime);
		}
	} else {
		fmt.Println("login: valid cookie found!");
	}
	fmt.Println("login: Successful login for user " + username);
	http.Redirect(w, r, "/home/", http.StatusSeeOther);
}

func logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id");
	if err == http.ErrNoCookie {
		fmt.Println("logout: No cookie found");
		return 
	}
	fmt.Println("LOGOUT: SUCCESSFULLY LOGGED OUT USER ====> " + cookie.Value);
	query := "delete from sessions where sessionid = ?";
	db.QueryRow(query, cookie.Value);
	http.Redirect(w, r, "/login/", http.StatusSeeOther);
}

func validatecookie(w http.ResponseWriter, r *http.Request) bool {
    cookie, err := r.Cookie("session_id");
    if err != nil {
		// check client side and see if they have a session cookie
        if err == http.ErrNoCookie {
			fmt.Println("validatecookie: No cookie found");
            // http.Error(w, "No cookie found", http.StatusUnauthorized)
            return false
        }
        // Handle other errors
		fmt.Println("validatecookie: Error retrieving cookie");
        // http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
        return false
    }

	var userid int;
	var expires string;
	layout := "2006-01-02 15:04:05";

	fmt.Println("validating the cookie ====> " + cookie.Value);

	query := "select userid, expires from sessions where sessionid = ?";
	db.QueryRow(query, cookie.Value).Scan(&userid, &expires);
	datetimeexpires, err := time.Parse(layout, expires);
	datetimeexpires = datetimeexpires.UTC();

	// check to ensure cookie is in the db
	if (userid == 0) { 
		// http.Error(w, "Cookie invalid", http.StatusUnauthorized);
		fmt.Println("validatecookie: Cookie is not in the database.");
		return false;
	}

	currenttime := time.Now().UTC();
	//check if the cookie is not expired
	if (currenttime.After(datetimeexpires)) {
		fmt.Println("validatecookie: COOKIE EXPIRED");
		return false;
	}
	// cookie is valid! add to cookie timer, and allow whatever action was
	// asked
	newexpirationtime := time.Now().UTC().Add(1 * time.Hour);
	query = "update sessions set expires = ?, updatedat = ? where sessionid = ?";
	db.QueryRow(query, newexpirationtime, currenttime, cookie.Value);
	fmt.Println("validatecookie: Valid cookie! Updated the expiration!");
	cookie.Expires = newexpirationtime;
	http.SetCookie(w, cookie);
	return true;
}

func generatecookie(w http.ResponseWriter, r *http.Request, username string) *http.Cookie {
	expirationtime := time.Now().UTC().Add(1 * time.Hour) // Set expiration to 1 hour from now
	seconds := expirationtime.Unix();
	sessionid := hash(username + "-"+ salt + "-" + strconv.FormatInt(seconds, 10));
	cookie := &http.Cookie{
        Name: "session_id",
        Value: sessionid,
		Path: "/",
        Expires: expirationtime, // Set expiration to 1 hour
        HttpOnly: true, // Prevent JavaScript access
        Secure:   true,  // Ensure it's sent only over HTTPS
    }
	return cookie
}
