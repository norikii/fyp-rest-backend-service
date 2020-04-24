package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tatrasoft/fyp-rest-backend-service/db"
	"github.com/tatrasoft/fyp-rest-backend-service/model"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

const (
	ordersCollection = "orders"
)

func CreateOrder(response http.ResponseWriter, request *http.Request) {
	order := &model.Order{}

	err := json.NewDecoder(request.Body).Decode(&order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}

	collection, err := db.GetDBCollection(dbname, ordersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	var totalPrice float32

	for _, orderItem := range order.Items {
		totalPrice += orderItem.ItemPrice
	}

	// updates the timestamps
	order.CreatedAt = time.Now().Unix()
	order.UpdatedAt =  time.Now().Unix()
	order.TotalPrice = totalPrice

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to insert the record: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(result)
}

func GetOrder(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to retrieve parameters: %v", err),
		})
		return
	}

	var order model.Order

	collection, err := db.GetDBCollection(dbname, ordersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	err = collection.FindOne(ctx, model.Order{ID: id,}).Decode(&order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to find the record: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(order)
}

func GetOrders(response http.ResponseWriter, request *http.Request) {
	var orders []model.Order
	collection, err := db.GetDBCollection(dbname, ordersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to retrieve the cursor: %v", err),
		})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		order :=  model.Order{}
		err := cursor.Decode(&order)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(response_models.ErrorResponse{
				ErrorCode: http.StatusInternalServerError,
				ErrorMessage:  fmt.Sprintf("unable to decode the order: %v", err),
			})
			return
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		//
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(orders)
}

func DeleteOrder(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to retrieve parameters: %v", err),
		})
		return
	}

	collection, err := db.GetDBCollection(dbname, ordersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to delete the record: %v", err),
		})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func UpdateOrder(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var order model.Order

	// get id
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to retrieve parameters: %v", err),
		})
		return
	}

	collection, err := db.GetDBCollection(dbname, ordersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": id}

	// read update model
	err = json.NewDecoder(request.Body).Decode(&order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}

	// prepare model update
	update := bson.D{
		{
			"$set", bson.D{
			{"table_id", order.TableID},
			{"staff_user_id", order.StaffUserID},
			{"guest_user_id", order.GuestUserID},
			{"items", order.Items},
			{"total_price", order.TotalPrice},
			{"created_at", order.CreatedAt},
			{"ready_at", order.ReadyAt},
			{"delivered_at", order.DeletedAt},
			{"payed_at", order.PayedAt},
			{"updated_at", time.Now().Unix()},
			{"deleted_at", order.DeletedAt},

		},
		},
	}

	sr := collection.FindOneAndUpdate(ctx, filter, update)
	if sr.Err() != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to update the record: %v", err),
		})
		return
	}
	order.ID = id

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(order)
}



