package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rodrigocollavo/apio/pkg/server"
	"time"
)

func main() {
	server.InitStorage(readManifest())

	mainRoute := mux.NewRouter().PathPrefix("/api/").Subrouter()
	addListRoutesEndoint(mainRoute)
	addAPIRoutes(mainRoute)

	srv := &http.Server{
		Handler:      mainRoute,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func readManifest() string {
	file, err := os.Open("manifest.sample.json")
	if err != nil {
		log.Fatalf("can't readmin manifest file. err: %s", err.Error())
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	return string(data)
}

func addListRoutesEndoint(route *mux.Router) {
	route.HandleFunc("/routes", func(writer http.ResponseWriter, request *http.Request) {
		route.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tpl, _ := route.GetPathTemplate()
			methods, _ := route.GetMethods()
			for _, method := range methods {
				writer.Write([]byte(fmt.Sprintf("%s	%s\n", method, tpl)))
				fmt.Println(fmt.Sprintf("%s	%s", method, tpl))
			}
			return nil
		})
	})
}

func addAPIRoutes(router *mux.Router) {
	for _, collection := range server.Storage.GetCollectionDefinitions() {
		apiRoute := router.PathPrefix(fmt.Sprintf("/%s/", collection.Name)).Subrouter()
		apiRoute.Use(server.ValidateId(collection))
		apiRoute.HandleFunc("/{id}", server.GetHandler(collection)).Methods(http.MethodGet)
		apiRoute.HandleFunc("/", server.ParseBody(server.PutHandler(collection))).Methods(http.MethodPut)
		apiRoute.HandleFunc("/{id}", server.ParseBody(server.PostHandler(collection))).Methods(http.MethodPost)
		apiRoute.HandleFunc("/{id}", server.DeleteHandler(collection)).Methods(http.MethodDelete)
	}
}
