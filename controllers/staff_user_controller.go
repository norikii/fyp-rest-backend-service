package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"

	"github.com/tatrasoft/fyp-rest-backend-service/db"
	"github.com/tatrasoft/fyp-rest-backend-service/model"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/auth"
	"github.com/tatrasoft/fyp-rest-backend-service/utils/response_models"
)

const (
	dbname = "mydb"
	staffUsersCollection = "staff_users"
)

// CreateStaffUser creates or registers a staff user to the system
func CreateStaffUser(response http.ResponseWriter, request *http.Request) {
	staffUser := &model.StaffUser{}

	// decoding the request body
	err := json.NewDecoder(request.Body).Decode(&staffUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to decode request body: %v", err),
		})
		return
	}

	// getting database collection
	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to connect to db: %v", err),
		})
		return
	}

	// check if user email already in database
	findRes := collection.FindOne(context.Background(), model.StaffUser{Email: staffUser.Email})
	if findRes.Err() == nil {
		response.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusNotAcceptable,
			ErrorMessage:  "email already used",
		})
		return
	}

	// setting up the creation timestamp
	staffUser.CreatedAt = time.Now().Unix()
	staffUser.UpdatedAt = time.Now().Unix()
	// hashing password
	pwd := staffUser.Password
	pwdHash, err := auth.HashAndSaltPwd(pwd)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to create hash: %v", err),
		})
		return
	}
	staffUser.Password = pwdHash

	// inserting the data into the database
	result, err := collection.InsertOne(context.Background(), staffUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("unable to insert entry: %v", err),
		})
		return
	}

	// preparing the response
	response.WriteHeader(http.StatusOK)
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(result)
}

// LoginStaffUser logs in already created user to the system
func LoginStaffUser(response http.ResponseWriter, request *http.Request) {
	staffUser := &model.StaffUser{}
	// decoding the login details from the request to the staff user struct
	err := json.NewDecoder(request.Body).Decode(staffUser)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMessage:  fmt.Sprintf("invalid user detials: %v", err),
		})
		return
	}

	// validating the staff user's credentials
	successResponse, err := checkStaffUser(staffUser.Email, staffUser.Password)
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
func checkStaffUser(email string, password string) (*response_models.SuccessLoginResponse, error) {
	staffUser := &model.StaffUser{}

	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	err = collection.FindOne(context.TODO(), model.StaffUser{Email: email}).Decode(staffUser)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %v", err)
	}

	isValid, err := auth.IsValidPassword(staffUser.Password, password)
	if err != nil && !isValid {
		return nil, fmt.Errorf("invalid password: %v", err)
	}

	token, err := auth.CreateJWTToken(staffUser.ID, staffUser.FirstName, staffUser.Email, staffUser.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("unable to create token: %v", err)
	}

	// not send password in the response
	staffUser.Password = ""

	return &response_models.SuccessLoginResponse{
		StatusCode: http.StatusOK,
		Message:    "user is logged in",
		Token: 		token,
		User:       staffUser,
	}, nil
}

func LogOut(response http.ResponseWriter, request *http.Request) {
	// TODO implement logout
}

// GetStaffUsers retrieves all staff user objects from the database
func GetStaffUsers(response http.ResponseWriter, request *http.Request) {
	var staffUsers []model.StaffUser
	ctx := context.Background()

	// getting the database collection
	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
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
		var staffUser model.StaffUser
		cursor.Decode(&staffUser)
		staffUser.Password = ""
		staffUsers = append(staffUsers, staffUser)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMessage:  fmt.Sprintf("curosor err: %v", err),
		})
		return
	}

	// preparing the response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(staffUsers)
}

// FindUser controller retrieves the user from the
// database or returns error response if entry not present
func FindStaffUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// getting the staff user id parameter
	id, _ := primitive.ObjectIDFromHex(params["id"])

	staffUser := &model.StaffUser{}
	// creating the database collection
	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
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
	err = collection.FindOne(ctx, model.StaffUser{ID: id}).Decode(staffUser)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(response_models.ErrorResponse{
			ErrorCode: http.StatusNotFound,
			ErrorMessage:  fmt.Sprintf("user not found: %v", err),
		})
		return
	}
	staffUser.Password = ""

	// prepare the response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(staffUser)
}

// UpdateStaffUser updates requested fields of the staff user record from the database
func UpdateStaffUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	staffUser := &model.StaffUser{}

	// decoding the login details from the request to the staff user struct
	err := json.NewDecoder(request.Body).Decode(&staffUser)
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
	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
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
			{"first_name", staffUser.FirstName},
			{"last_name", staffUser.LastName},
			{"email", staffUser.Email},
			{"password", staffUser.Password},
			{"is_admin", staffUser.IsAdmin},
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
	staffUser.ID = id
	staffUser.Password = ""

	// preparing the user response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(staffUser)
}

// DeleteStaffUser removes requested staff user from the database
func DeleteStaffUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// getting the staff user id parameter
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// creating the database collection
	collection, err := db.GetDBCollection(dbname, staffUsersCollection)
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

