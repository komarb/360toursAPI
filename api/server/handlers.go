package server

import (
	"../models"
	"../strutils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

func createTour(w http.ResponseWriter, r *http.Request) {
	var tour models.Tour
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&tour)


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, tour)
	if err != nil {
		log.Panic(err)
	}
	id := result.InsertedID
	tour.ID, err = primitive.ObjectIDFromHex(id.(primitive.ObjectID).Hex())

	json.NewEncoder(w).Encode(tour)
}

func getTour(w http.ResponseWriter, r *http.Request) {
	emptySlice := make([]string, 0)
	var tour models.Tour
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&tour)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	loadedToursID, err := getLoadedToursCookie(r)
	if err != nil {
		setLoadedToursCookie(w, emptySlice)
	}
	if !strutils.StringInSlice(tour.ID.Hex(), loadedToursID) {
		loadedToursID = append(loadedToursID, tour.ID.Hex())
		setLoadedToursCookie(w, loadedToursID)
	}

	json.NewEncoder(w).Encode(tour)
}

func getAllTours(w http.ResponseWriter, r *http.Request) {
	IDs := make([]string, 0)
	tours := make([]models.Tour, 0)

	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &tours)
	if err != nil {
		log.Fatal(err)
	}
	for _, tour := range tours {
		IDs = append(IDs, tour.ID.Hex())
	}

	setLoadedToursCookie(w, IDs)
	json.NewEncoder(w).Encode(tours)
}

// Gets "loadedTours" cookie and makes slice with tours which IDs are not in "loadedTours"
func getUnloadedTours(w http.ResponseWriter, r *http.Request) {
	tours := make([]models.Tour, 0)
	IDs := make([]string, 0)
	found := false
	w.Header().Set("Content-Type", "application/json")
	loadedToursID, err := getLoadedToursCookie(r)
	if err != nil {
		getAllTours(w, r)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	for cur.Next(ctx) {
		var tour models.Tour
		err := cur.Decode(&tour)
		if err != nil {
			log.Panic(err)
		}
		IDs = append(IDs, tour.ID.Hex())
		for _, id := range loadedToursID {
			if tour.ID.Hex() == id {
				found = true
			}
		}
		if !found {
			tours = append(tours, tour)
		}
		found = false
	}
	setLoadedToursCookie(w, IDs)
	json.NewEncoder(w).Encode(tours)
}
func updateTour(w http.ResponseWriter, r *http.Request) {
	var updateTour models.Tour
	var updatedTour models.Tour
	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&updateTour)
	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	updateResult, err := collection.ReplaceOne(ctx, filter, updateTour)
	if err != nil || updateResult.MatchedCount == 0{
		http.NotFound(w, r)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&updatedTour)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedTour)
}

func deleteTour(w http.ResponseWriter, r *http.Request) {
	var i = 0
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}
	loadedToursID, err := getLoadedToursCookie(r)
	if err != nil {
		log.Panic(err)
	}
	for _, id := range loadedToursID {
		if id == data["id"] {
			break
		}
		i++
	}

	loadedToursID = strutils.RemoveString(loadedToursID, i)
	setLoadedToursCookie(w, loadedToursID)
	w.WriteHeader(200)
}
