package database

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/icrxz/crm-api-core/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func getClient(config config.Database) (*mongo.Client, error) {

	//const uri = "mongodb+srv://<YOUR USERNAME HERE>:<YOUR PASSWORD HERE>@cluster0.e5akf.mongodb.net/myFirstDatabese?retryWrites=true&w=majority"
	// usrname, password
	var uri = config.ConnStr
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

//var mongoDatabase *mongo.Database

func NewDBInstance(config config.Database) (*mongo.Client, error) {
	if mongoClient == nil {
		client, err := getClient(config)
		return client, err
		//if err == nil {
		//	mongoDatabase = client.Database("go-crm")
		//}
		//return mongoDatabase, err
	}
	return mongoClient, nil
	//return mongoDatabase, nil
}
