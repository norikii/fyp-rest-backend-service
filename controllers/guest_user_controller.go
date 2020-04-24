package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"

	"github.com/tatrasoft/fyp-rest-backend-service/db"
	"github.com/tatrasoft/fyp-rest-backend-service/model"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/auth"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
)

const (
	guestUsersCollection = "guest_users"
)

// CreateGuestUser creates or registers a staff user to the system
func CreateGuestUser(response http.ResponseWriter, request *http.Request) {
	guestUser := &model.GuestUser{}

	// decoding the request body
	err := json.NewDecoder(request.Body).Decode(&guestUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}

	// getting database collection
	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}

	// check if user email already in database
	findRes := collection.FindOne(context.Background(), model.GuestUser{Email: guestUser.Email})
	if findRes.Err() == nil {
		response.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusNotAcceptable,
			ErrorMessage:  "email already used",
		})
		return
	}

	// setting up the creation timestamp
	guestUser.CreatedAt = time.Now().Unix()
	guestUser.UpdatedAt = time.Now().Unix()
	// hashing password
	pwd := guestUser.Password
	pwdHash, err := auth.HashAndSaltPwd(pwd)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to create hash: %v", err),
		})
		return
	}
	guestUser.Password = pwdHash

	// inserting the data into the database
	result, err := collection.InsertOne(context.Background(), guestUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to insert entry: %v", err),
		})
		return
	}

	// preparing the response
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(result)
}

// LoginGuestUser logs in already created user to the system
func LoginGuestUser(response http.ResponseWriter, request *http.Request) {
	guestUser := &model.GuestUser{}
	// decoding the login details from the request to the staff user struct
	err := json.NewDecoder(request.Body).Decode(guestUser)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("invalid user detials: %v", err),
		})
		return
	}

	// validating the staff user's credentials
	successResponse, err := checkGuestUser(guestUser.Email, guestUser.Password)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("invalid user detials: %v", err),
		})
		return
	}

	// preparing response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(&successResponse)
}

// checks if passed staff user's email and password are correct
func checkGuestUser(email string, password string) (*response_models.SuccessLoginResponse, error) {
	guestUser := &model.GuestUser{}

	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	err = collection.FindOne(context.TODO(), model.GuestUser{Email: email}).Decode(guestUser)
	if err != nil {
		return nil, fmt.Errorf("staff user cannot be found: %v", err)
	}

	isValid, err := auth.IsValidPassword(guestUser.Password, password)
	if err != nil && !isValid {
		return nil, fmt.Errorf("invalid password: %v", err)
	}

	token, err := auth.CreateJWTToken(guestUser.ID, guestUser.Email, guestUser.Email, false)
	if err != nil {
		return nil, fmt.Errorf("unable to create token: %v", err)
	}

	// setting up the filter
	filter := bson.M{"email": email}
	// prepare update object
	update := bson.D{
		{
			"$set", bson.D{
			{"logged_at", time.Now().Unix()},
		},
		},
	}
	// update is the logged_at timestamp
	res := collection.FindOneAndUpdate(context.Background(), filter, update)
	if res.Err() != nil {
		return nil, fmt.Errorf("unable to update logged_at field: %v", err)
	}

	// not send password in the response
	guestUser.Password = ""

	return &response_models.SuccessLoginResponse{
		StatusCode: http.StatusOK,
		Message:    "user is logged in",
		Token: 		token,
		User:       guestUser,
	}, nil
}

// GetGuestUsers retrieves all staff user objects from the database
func GetGuestUsers(response http.ResponseWriter, request *http.Request) {
	var guestUsers []model.GuestUser
	ctx := context.Background()

	// getting the database collection
	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("cannot connect to db: %v", err),
		})
		return
	}

	// returns a cursor with staff user records from the database
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("curosor err: %v", err),
		})
		return
	}
	defer cursor.Close(ctx)
	// iterate over the entries from the database and decoding every object into the staff user struct
	// and adding them into the collection of staff user objects
	for cursor.Next(ctx) {
		var guestUser model.GuestUser
		cursor.Decode(&guestUser)
		guestUser.Password = ""
		guestUsers = append(guestUsers, guestUser)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}

	// preparing the response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(guestUsers)
}

// FindUser controller retrieves the user from the
// database or returns error response if entry not present
func FindGuestUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// getting the staff user id parameter
	id, _ := primitive.ObjectIDFromHex(params["id"])

	guestUser := &model.GuestUser{}
	// creating the database collection
	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to get db collection: %v", err),
		})
		return
	}
	ctx, _ := context.WithCancel(context.Background())

	// checking if requested staff user id is in the database
	err = collection.FindOne(ctx, model.GuestUser{ID: id}).Decode(guestUser)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			ErrorMessage:  fmt.Sprintf("user not found: %v", err),
		})
		return
	}
	guestUser.Password = ""

	// prepare the response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(guestUser)
}

// UpdateGuestUser updates requested fields of the staff user record from the database
func UpdateGuestUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	guestUser := &model.GuestUser{}

	// decoding the login details from the request to the staff user struct
	err := json.NewDecoder(request.Body).Decode(&guestUser)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("invalid staff user detials: %v", err),
		})
		return
	}

	// getting the staff user id parameter
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// creating the database collection
	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to get db collection: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	// setting up the filter
	filter := bson.M{"_id": id}

	// prepare update object
	update := bson.D{
		{
			"$set", bson.D{
			{"email", guestUser.Email},
			{"password", guestUser.Password},
			{"isAdmin", false},
			{"updated_at", time.Now().Unix()},
		},
		},
	}

	// updating the staff user record in the database
	res := collection.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			ErrorMessage:  fmt.Sprintf("unable to update the record: %v", err),
		})
		return
	}
	// setting staff user id to its initial value
	guestUser.ID = id
	guestUser.Password = ""

	// preparing the user response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(guestUser)
}

// DeleteGuestUser removes requested staff user from the database
func DeleteGuestUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// getting the staff user id parameter
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// creating the database collection
	collection, err := db.GetDBCollection(dbname, guestUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("database error: %v", err),
		})
		return
	}
	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	// setting the filter
	filter := bson.M{"_id": id}

	// deleting the item from the database
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("cannot delete the staff user entry: %v", err),
		})
		return
	}

	// preparing response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}


