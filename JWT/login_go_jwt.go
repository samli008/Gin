package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// define the URL and method
	url := "http://localhost:8080/login"
	method := "POST"

	payload := strings.NewReader(`{
    "username": "liyang",
    "password": "123456"
	}`)

	client := &http.Client{}
	req1, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req1.Header.Add("Content-Type", "application/json")

	res1, err := client.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res1.Body.Close()

	body, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("token is %s\n", string(body))

	// get the token
	token := string(body)
	//time.Sleep(time.Minute * 2)

	url2 := "http://localhost:8080/getAllBooks"
	method2 := "GET"

	req2, err := http.NewRequest(method2, url2, nil)

	if err != nil {
		fmt.Println(err)
		return

	}
	req2.Header.Add("Token", token)

	res2, err := client.Do(req2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res2.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body2))
}
