package main

import (
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

const form = `<html>
	<head>
		<title></title>
	</head>
	<body>
		<form action="/" method="post">
			<label>long url <input type="text" name ="longurl"></label>
			<input type="submit" value ="submit">
		</form>
	</body>
</html>
`
const length = 8

var urls map[string]string = make(map[string]string)

func createShortURL(letters []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func addURL(LongURL string) (string, bool) {
	if len(LongURL) > 0 {
		letters := []rune(LongURL)
		shortURL := createShortURL(letters, length)
		if len(shortURL) > 0 {
			urls[shortURL] = LongURL
			return shortURL, true
		}
	}
	return "", false
}

func findLongURL(shortURL string) (string, bool) {
	if len(shortURL) > 0 {
		if val, ok := urls[shortURL]; ok {
			return val, true
		}
	}
	return "", false
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		LongURL := r.FormValue("longurl")
		if shortURL, ok := addURL(LongURL); ok {
			path := "http://" + r.Host + "/" + shortURL
			w.Write([]byte(path))
		}
		return
	}
	w.Write([]byte(form))
}

func SubPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if LongURL, ok := findLongURL(r.URL.Path[1:]); ok {
			w.Header().Set("Location", LongURL)
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundPage)
	router.HandleFunc("/", MainPage)
	router.HandleFunc("/{id}", SubPage)
	http.Handle("/", router)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
}
