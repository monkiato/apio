package storage

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"monkiato/apio/internal/data"
	"os"
	"time"
)

const (
	defaultMongodbHost = "localhost:27017"
	defaultMongodbName = "apio"
)

//MongoStorage structure for the storage using a MongoDB
type MongoStorage struct {
	collectionsDefinitions    []data.CollectionDefinition
	collectionsDefinitionsMap map[string]data.CollectionDefinition
	collectionHandlers        map[string]CollectionHandler
	client                    *mongo.Client
	host                      string
	dbName                    string
}

//MongoCollectionHandler  data handler used for a specific collection
type MongoCollectionHandler struct {
	db         *mongo.Database
	collection data.CollectionDefinition
}

//NewMongoStorage create a new MongoStorage instance
//MongoDB connection data con be set by environment variables
//MONGODB_HOST and MONGODB_NAME
func NewMongoStorage() Storage {
	host := defaultMongodbHost
	dbName := defaultMongodbName
	if envHost, exists := os.LookupEnv("MONGODB_HOST"); exists {
		host = envHost
	}
	if envDbName, exists := os.LookupEnv("MONGODB_NAME"); exists {
		dbName = envDbName
	}
	return &MongoStorage{
		collectionHandlers:        map[string]CollectionHandler{},
		collectionsDefinitionsMap: map[string]data.CollectionDefinition{},
		host:                      host,
		dbName:                    dbName,
	}
}

func newMongoStorageCollectionHandler(db *mongo.Database, collection data.CollectionDefinition) CollectionHandler {
	return &MongoCollectionHandler{
		db:         db,
		collection: collection,
	}
}

//GetItem implements storage.CollectionHandler.GetItem
func (msc *MongoCollectionHandler) GetItem(itemID string) (interface{}, bool) {
	objID, _ := primitive.ObjectIDFromHex(itemID)
	// fetch item
	res := msc.db.Collection(msc.collection.Name).
		FindOne(
			createContext(),
			bson.M{"_id": objID},
			options.FindOne().SetProjection(bson.D{{"_id", 0}}))

	// check fetching errors
	if res.Err() != nil {
		fmt.Printf("unable to fetch item id %s. err: %s", itemID, res.Err().Error())
		return nil, false
	}

	// decode data
	var itemBson interface{}
	if err := res.Decode(&itemBson); err != nil {
		fmt.Printf("unable to decode DB data for item id %s. err: %s", itemID, err)
		return nil, false
	}

	// convert bson data to go map
	var item map[string]interface{}
	b, _ := bson.Marshal(itemBson)
	bson.Unmarshal(b, &item)
	return item, true
}

//AddItem implements storage.CollectionHandler.AddItem
func (msc *MongoCollectionHandler) AddItem(item map[string]interface{}) (string, error) {
	res, err := msc.db.Collection(msc.collection.Name).InsertOne(createContext(), item)
	if err != nil {
		fmt.Printf("unable to add new item. err: " + err.Error())
		return "", nil
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	log.Debugf("created new item %s.%s", msc.collection.Name, id)
	return id, nil
}

//UpdateItem implements storage.CollectionHandler.UpdateItem
func (msc *MongoCollectionHandler) UpdateItem(itemID string, newItem map[string]interface{}) error {
	objID, _ := primitive.ObjectIDFromHex(itemID)
	_, err := msc.db.Collection(msc.collection.Name).UpdateOne(createContext(), bson.D{{"_id", objID}}, bson.D{{"$set", newItem}})
	if err != nil {
		fmt.Printf("unable to update item. err: " + err.Error())
		return err
	}
	log.Debugf("updated item %s.%s", msc.collection.Name, itemID)
	return nil
}

//DeleteItem implements storage.CollectionHandler.DeleteItem
func (msc *MongoCollectionHandler) DeleteItem(itemID string) error {
	objID, _ := primitive.ObjectIDFromHex(itemID)
	_, err := msc.db.Collection(msc.collection.Name).DeleteOne(createContext(), bson.M{"_id": objID})
	if err != nil {
		fmt.Printf("unable to delete item. err: " + err.Error())
		return err
	}
	log.Debugf("deleted item %s.%s", msc.collection.Name, itemID)
	return nil
}

//List implements storage.CollectionHandler.List
func (msc *MongoCollectionHandler) List(lastItemID string) []interface{} {
	return nil
}

//Initialize implements storage.Storage.Initialize
func (ms *MongoStorage) Initialize(manifest string) {
	ctx := createContext()
	uri := fmt.Sprintf("mongodb://%s", ms.host)

	log.Debugf("connecting to mongoDB: %s", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic("unable to connect to Mongo db. " + uri)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("Mongo db server not found at " + uri)
	}
	ms.client = client

	ms.initializeCollectionDefinitions(manifest)
	ms.initializeCollections()
}

func createContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

//GetCollectionDefinitions implements storage.Storage.GetCollectionDefinitions
func (ms *MongoStorage) GetCollectionDefinitions() []data.CollectionDefinition {
	return ms.collectionsDefinitions
}

//GetCollection implements storage.Storage.GetCollection
func (ms *MongoStorage) GetCollection(collectionName string) (CollectionHandler, error) {
	if collection, ok := ms.collectionsDefinitionsMap[collectionName]; ok {
		collectionHandler, exists := ms.collectionHandlers[collectionName]
		if !exists {
			collectionHandler = newMongoStorageCollectionHandler(ms.client.Database(ms.dbName), collection)
			ms.collectionHandlers[collectionName] = collectionHandler
		}
		return collectionHandler, nil
	}
	return nil, fmt.Errorf("collection %s not found", collectionName)
}

func (ms *MongoStorage) initializeCollectionDefinitions(manifest string) {
	log.Debugf("parsing manifest...")
	err := json.Unmarshal([]byte(manifest), &ms.collectionsDefinitions)
	if err != nil {
		log.Fatal("Unable to parse manifest")
	}
	log.Debugf("manifest parsed successfully")
}

func (ms *MongoStorage) initializeCollections() {
	for _, collectionDefinition := range ms.collectionsDefinitions {
		ms.collectionsDefinitionsMap[collectionDefinition.Name] = collectionDefinition
	}
}
