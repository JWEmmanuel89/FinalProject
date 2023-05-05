//Filename: main

package main

import (
	"errors"
	"log"
	"net/http"
)

// Function to initialize new cookie
func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "myCookie",
		Value:    "Final Project on Cookies!",
		Path:     "/",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	//Send cookie to client and adds header to response
	http.SetCookie(w, &cookie)
	w.Write([]byte("Cookie has been set!"))
}

// Function to retrieve cookie from request using cookie name
func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("myCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "no cookie found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	//Cookie value in respnse body
	w.Write([]byte(cookie.Value))
}

func main() {
	//Create web server
	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)

	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
