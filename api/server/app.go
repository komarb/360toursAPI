package server

import (
	"../config"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB *mongo.Client
}

var collection *mongo.Collection

func (a *App) Init(config *config.Config) {
	DBUri := "mongodb://" + config.DB.Host + ":" + config.DB.Port
	client, err := mongo.NewClient(options.Client().ApplyURI(DBUri))
	if err != nil {
		log.Panic(err)
	}

	// Create db connect
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection = client.Database(config.DB.DBName).Collection(config.DB.CollectionName)

	a.Router = mux.NewRouter()
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
func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Panic(err)
	}
}