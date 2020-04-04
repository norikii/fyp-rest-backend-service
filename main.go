package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemName string `json:"item_name,omitempty" bson:"item_name,omitempty"`
	ItemDescription string `json:"item_description,omitempty" bson:"item_description,omitempty"`
	ItemType string `json:"item_type,omitempty" bson:"item_type,omitempty"`
	ItemImg string `json:"item_img,omitempty" bson:"item_img,omitempty"`
	ItemPrice float32 `json:"item_price,omitempty" bson:"item_price,omitempty"`
	EstimatePrepareTime  int64 `json:"estimate_prepare_time,omitempty" bson:"estimate_prepare_time,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt int64 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

var client *mongo.Client
var clientErr error

func CreateItem(w http.ResponseWriter, request *http.Request) {
	var item Item
	item.CreatedAt = time.Now().Unix()
	item.UpdatedAt =  time.Now().Unix()

	json.NewDecoder(request.Body).Decode(&item)
	collection := client.Database("mydb").Collection("items")
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, item)

	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var item Item
	collection := client.Database("mydb").Collection("items")
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Item{ID: id,}).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(item)
}

func GetItems(response http.ResponseWriter, r *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("Access-Control-Allow-Origin", "*")
	response.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	response.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Add("Access-Control-Max-Age", "86400")

	var items []Item
	collection := client.Database("mydb").Collection("items")
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var item Item
		cursor.Decode(&item)
		items = append(items, item)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(items)
}

func DeleteItem(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Add("Access-Control-Max-Age", "86400")

	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("mydb").Collection("items")
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	result, _ := collection.DeleteOne(ctx, filter)

	json.NewEncoder(writer).Encode(result)
}

func UpdateItem(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Add("Access-Control-Max-Age", "86400")

	params := mux.Vars(request)
	var item Item

	// get id
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("mydb").Collection("items")
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	// read update model
	json.NewDecoder(request.Body).Decode(&item)

	// prepare model update
	update := bson.D{
		{
			"$set", bson.D{
				{"item_name", item.ItemName},
				{"item_description", item.ItemDescription},
				{"item_type", item.ItemType},
				{"item_img", item.ItemImg},
				{"item_price", item.ItemPrice},
				{"estimate_prepare_time", item.EstimatePrepareTime},
				{"updated_at", time.Now().Unix()},

			},
		},
	}

	collection.FindOneAndUpdate(ctx, filter, update)
	item.ID = id

	json.NewEncoder(writer).Encode(item)
}

func main() {
	fmt.Println("Starting the application...")

	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	client, clientErr = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	connErr := client.Connect(ctx)
	if connErr != nil {
		log.Fatal(connErr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/item", CreateItem).Methods("POST")
	router.HandleFunc("/items", GetItems).Methods("GET")
	router.HandleFunc("/item/{id}", GetItem).Methods("GET")
	router.HandleFunc("/items/{id}", UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", DeleteItem).Methods("DELETE")
	err := http.ListenAndServe(":12345", router)
	if err != nil {
		log.Fatal(err)
	}
}

