package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"github.com/tatrasoft/fyp-rest-backend-service/controllers"
	"github.com/tatrasoft/fyp-rest-backend-service/db"
	"github.com/tatrasoft/fyp-rest-backend-service/middleware"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var clientErr error

func main() {
	fmt.Println("Starting the application...")

	ctx, _:= context.WithTimeout(context.Background(), 10*time.Second)

	client, clientErr = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	fmt.Println("Connecting to the database...")
	connErr := client.Connect(ctx)
	if connErr != nil {
		log.Fatal(connErr)
	}

	fmt.Println("Checking the db connection...")
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMongoClient(client)
	fmt.Println("Database connection have been establish")

	router := mux.NewRouter().StrictSlash(true)
	//router.Use(middleware.CommonMiddleware)

	// staff users
	router.HandleFunc("/user/staff/register", controllers.CreateStaffUser).Methods("POST")
	router.HandleFunc("/user/staff/login", controllers.LoginStaffUser).Methods("POST")

	// guest users
	router.HandleFunc("/user/guest/register", controllers.CreateGuestUser).Methods("POST")
	router.HandleFunc("/user/guest/login", controllers.LoginGuestUser).Methods("POST")



	// authenticated routes
	apiAuthRoute := router.PathPrefix("/auth").Subrouter()
	apiAuthRoute.Use(middleware.AuthMiddleware)
	apiAuthRoute.HandleFunc("/dashboard", controllers.Dashboard).Methods("GET")

	// items
	apiAuthRoute.HandleFunc("/item", controllers.CreateItem).Methods("POST")
	apiAuthRoute.HandleFunc("/items", controllers.GetItems).Methods("GET")
	apiAuthRoute.HandleFunc("/item/{id}", controllers.GetItem).Methods("GET")
	apiAuthRoute.HandleFunc("/item/{id}", controllers.UpdateItem).Methods("PUT")
	apiAuthRoute.HandleFunc("/item/{id}", controllers.DeleteItem).Methods("DELETE")

	// orders
	apiAuthRoute.HandleFunc("/order", controllers.CreateOrder).Methods("POST")
	apiAuthRoute.HandleFunc("/orders", controllers.GetOrders).Methods("GET")
	apiAuthRoute.HandleFunc("/order/{id}", controllers.GetOrder).Methods("GET")
	apiAuthRoute.HandleFunc("/order/{id}", controllers.UpdateOrder).Methods("PUT")
	apiAuthRoute.HandleFunc("/order/{id}", controllers.DeleteOrder).Methods("DELETE")

	// authenticated routes with admin privileges
	//apiAdminRoute := router.PathPrefix("/admin").Subrouter()
	//apiAdminRoute.Use(middleware.AdminMiddleware)

	// staff users
	apiAuthRoute.HandleFunc("/user/staff/all", controllers.GetStaffUsers).Methods("GET")
	apiAuthRoute.HandleFunc("/user/staff/{id}", controllers.FindStaffUser).Methods("GET")
	apiAuthRoute.HandleFunc("/user/staff/{id}", controllers.UpdateStaffUser).Methods("PUT")
	apiAuthRoute.HandleFunc("/user/staff/{id}", controllers.DeleteStaffUser).Methods("DELETE")

	// guest users
	apiAuthRoute.HandleFunc("/user/guest/all", controllers.GetGuestUsers).Methods("GET")
	apiAuthRoute.HandleFunc("/user/guest/{id}", controllers.FindGuestUser).Methods("GET")
	apiAuthRoute.HandleFunc("/user/guest/{id}", controllers.UpdateGuestUser).Methods("PUT")
	apiAuthRoute.HandleFunc("/user/guest/{id}", controllers.DeleteGuestUser).Methods("DELETE")

	var handler http.Handler
	{
		handler = handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Origin", "Authorization", "Content-Type"}),
			handlers.ExposedHeaders([]string{""}),
			handlers.MaxAge(10),
			handlers.AllowCredentials(),
		)(router)
		handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)
	}

	http.Handle("/", handler)
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is running and listening on port :12345...")
}
