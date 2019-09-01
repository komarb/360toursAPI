package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type Link struct {
	ID string `json:"id"`
	F  int    `json:"f"`
	O  int    `json:"o"`
}
type Picture struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Links []Link `json:"links"`
}
type Tour struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MainID   string             `json:"mainId" bson:"mainId"`
	Pictures []Picture          `json:"pictures"`
}

func createTour(w http.ResponseWriter, r *http.Request) {
	var tour Tour
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
	var tour Tour
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
	if !stringInSlice(tour.ID.Hex(), loadedToursID) {
		loadedToursID = append(loadedToursID, tour.ID.Hex())
		setLoadedToursCookie(w, loadedToursID)
	}

	json.NewEncoder(w).Encode(tour)
}

func getAllTours(w http.ResponseWriter, r *http.Request) {
	IDs := make([]string, 0)
	tours := make([]Tour, 0)

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
	tours := make([]Tour, 0)
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
		var tour Tour
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
	var updateTour Tour
	var updatedTour Tour
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

	loadedToursID = removeString(loadedToursID, i)
	setLoadedToursCookie(w, loadedToursID)
	w.WriteHeader(200)
}

func setLoadedToursCookie(w http.ResponseWriter, value []string) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Panic(err)
	}

	c := http.Cookie{}
	c.Name = "loadedTours"
	c.Value	= base64.StdEncoding.EncodeToString(jsonValue)
	c.Path = "/"

	http.SetCookie(w, &c)
	//w.Header().Set("Set-Cookie", "loadedTours="+string(jsonValue))
}
func getLoadedToursCookie(r *http.Request) ([]string, error) {
	var loadedToursID []string
	c, err := r.Cookie("loadedTours")
	if err != nil {
		return nil, err
	}
	jsonValue, err :=  base64.StdEncoding.DecodeString(c.Value)
	err = json.Unmarshal(jsonValue, &loadedToursID)
	return loadedToursID, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeString(a []string, i int) []string {
	emptySlice := make([]string, 0)
	if len(a) == 1 {
		return emptySlice
	}
	a[i] = a[len(a)-1]
	a[len(a)-1] = ""
	return a[:len(a)-1]
}

var collection *mongo.Collection

func main() {
	// Create db client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
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
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	collection = client.Database("36t").Collection("tours")

	r := mux.NewRouter()
	r.HandleFunc("/tours", getAllTours).Methods("GET")
	r.HandleFunc("/tours/unloaded", getUnloadedTours).Methods("GET")
	r.HandleFunc("/tours/{id}", getTour).Methods("GET")
	r.HandleFunc("/tours/{id}", deleteTour).Methods("DELETE")
	r.HandleFunc("/tours", createTour).Methods("POST")
	r.HandleFunc("/tours/{id}", updateTour).Methods("PUT")

	http.ListenAndServe(":8080", r)
}