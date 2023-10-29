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

func addUrl(longUrl, shortUrl string) bool {
	if len(longUrl) > 0 && len(shortUrl) > 0 {
		urls[shortUrl] = longUrl
		return true
	}

	return false
}

func findLongUrl(shortUrl string) (string, bool) {

	if len(shortUrl) > 0 {
		if val, ok := urls[shortUrl]; ok {
			return val, true
		}
	}

	return "", false
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		longUrl := r.FormValue("longurl")
		shortUrl := r.FormValue("shorturl")
		if addUrl(longUrl, shortUrl) {
			path := "http://" + r.Host + "/" + shortUrl
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
		if longUrl, ok := findLongUrl(r.URL.Path[1:]); ok {
			w.Header().Set("Location", longUrl)
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
