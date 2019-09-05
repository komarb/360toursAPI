package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (a *App) Init() {
	a.Router := mux.NewRouter()
	a.setRouters()
}
func (a *App) setRouters() {
	a.Router.HandleFunc("/tours", getAllTours).Methods("GET")
	a.Router.HandleFunc("/tours/unloaded", getUnloadedTours).Methods("GET")
	a.Router.HandleFunc("/tours/{id}", getTour).Methods("GET")
	a.Router.HandleFunc("/tours/{id}", deleteTour).Methods("DELETE")
	a.Router.HandleFunc("/tours", createTour).Methods("POST")
	a.Router.HandleFunc("/tours/{id}", updateTour).Methods("PUT")
}
func (a *App) Run() {
	http.ListenAndServe(":8080", a.Router)
}