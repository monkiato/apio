package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	mk_os "monkiato/apio/internal/os"
	"monkiato/apio/pkg/server"
	"net/http"
	"os"
	"time"
)

const defaultManifestPath = "/app/manifest.json"

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)

	debugMode := mk_os.GetIntEnv("DEBUG_MODE", 0)
	if debugMode == 1 {
		log.SetLevel(log.DebugLevel)
	}

	port := mk_os.GetEnv("SERVER_PORT", "80")
	storageType := mk_os.GetEnv("STORAGE_TYPE", server.StorageTypeMongoDB)

	server.InitStorage(readManifest(), storageType)

	mainRoute := mux.NewRouter().PathPrefix("/api/").Subrouter()
	addListRoutesEndpoint(mainRoute)
	addAPIRoutes(mainRoute)

	srv := &http.Server{
		Handler: mainRoute,
		Addr:    fmt.Sprintf(":%s", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Debugf("server ready. Running at %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func readManifest() string {
	manifestPath, found := os.LookupEnv("MANIFEST_PATH")
	if !found {
		manifestPath = defaultManifestPath
	}
	file, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("can't readmin manifest file. err: %s", err.Error())
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	return string(data)
}

func addListRoutesEndpoint(route *mux.Router) {
	log.Debug("adding all routes list...")
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
	log.Debug("add API routes...")
	for _, collection := range server.Storage.GetCollectionDefinitions() {
		log.Debugf("adding routes for collection '%s'", collection.Name)
		apiRoute := router.PathPrefix(fmt.Sprintf("/%s/", collection.Name)).Subrouter()
		apiRoute.Use(server.ValidateID(collection))
		apiRoute.HandleFunc("/{id}", server.GetHandler(collection)).Methods(http.MethodGet)
		apiRoute.HandleFunc("/", server.ParseBody(server.PutHandler(collection))).Methods(http.MethodPut)
		apiRoute.HandleFunc("/{id}", server.ParseBody(server.PostHandler(collection))).Methods(http.MethodPost)
		apiRoute.HandleFunc("/{id}", server.DeleteHandler(collection)).Methods(http.MethodDelete)
		apiRoute.HandleFunc("/", server.ListCollectionHandler(collection)).Methods(http.MethodGet)
	}
}
