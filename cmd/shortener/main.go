package main

import (
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
			<label>short url <input type="text" name ="shorturl"></label>
			<input type="submit" value ="submit">
		</form>
	</body>
</html>
`

var urls map[string]string = make(map[string]string)

func addURL(LongURL, shortURL string) bool {
	if len(LongURL) > 0 && len(shortURL) > 0 {
		urls[shortURL] = LongURL
		return true
	}

	return false
}

func findLongURL(shortURL string) (string, bool) {

	if len(shortURL) > 0 {
		if val, ok := urls[shortURL]; ok {
			return val, true
		}
	}

	return "", false
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		LongURL := r.FormValue("longurl")
		shortURL := r.FormValue("shorturl")
		if addURL(LongURL, shortURL) {
			path := "http://" + r.Host + "/" + shortURL
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(path))
		} else {
			http.Error(w, "Неверные данные ввода", http.StatusNotFound)
		}
		return
	}

	w.Write([]byte(form))
}

func subPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if LongURL, ok := findLongURL(r.URL.Path[1:]); ok {
			w.Header().Set("Location", LongURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", mainPage)
	router.HandleFunc("/{id}", subPage)
	http.Handle("/", router)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}

}
