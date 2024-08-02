package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// store all user
var userMap = map[string]string{
	"liyang": "123456",
	"samli":  "123456",
}

// USer type
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/getAllBooks", getAllBookHandler)
	fmt.Println("Listen on 8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		var user User

		// decode the request body into the struct
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			fmt.Fprintf(writer, "invalid body")
			return
		}

		if userMap[user.Username] == "" || userMap[user.Username] != user.Password {
			fmt.Fprintf(writer, "can not authenticate this user")
			return
		}

		token, err := generateJWT(user.Username)
		if err != nil {
			fmt.Fprintf(writer, "error in generating token")
		}

		fmt.Fprintf(writer, token)

	case "GET":
		fmt.Fprintf(writer, "only POST methods is allowed.")
		return
	}
}

var sampleSecretKey = []byte("samli008")

func generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func validateToken(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Header["Token"] == nil {
		fmt.Fprintf(w, "can not find token in header")
		return errors.New("Token error")
	}

	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return sampleSecretKey, nil
	})

	if token == nil {
		fmt.Fprintf(w, "invalid token")
		return errors.New("Token error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error")
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		fmt.Fprintf(w, "token expired")
		return errors.New("Token error")
	}

	return nil
}

func getAllBookHandler(w http.ResponseWriter, r *http.Request) {
	err := validateToken(w, r)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		books := getAllBook()
		json.NewEncoder(w).Encode(books)
	}
}

func getAllBook() []Book {
	return []Book{
		Book{
			Name:   "Book1",
			Author: "Author1",
		},
		Book{
			Name:   "Book2",
			Author: "Author2",
		},
		Book{
			Name:   "Book3",
			Author: "Author3",
		},
	}
}

type Book struct {
	Name   string
	Author string
}
