package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	// yamlPtr := flag.String("yaml", "paths.yaml", "Path to yaml file with path mapping")

	// flag.Parse()

	// yamlData, err := ioutil.ReadFile(*yamlPtr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// yamlHandler, err := urlshort.YAMLHandler(yamlData, mapHandler)
	// if err != nil {
	// 	panic(err)
	// }

	filePtr := flag.String("yaml", "paths.json", "Path to yaml file with path mapping")

	flag.Parse()

	data, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		log.Fatal(err)
	}

	pathHandler, err := urlshort.JSONHandler(data, mapHandler)
	if err != nil {
		panic(err)
	}
	runServer(pathHandler)

}

func runServer(h http.HandlerFunc) {
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", h)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
