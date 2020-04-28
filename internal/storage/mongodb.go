package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"rodrigocollavo/apio/internal/data"
	"time"
)

const (
	mongodbHost     = "mongodb://localhost:32771"
	mongodbName		= "default"
)

type MongoStorage struct {
	collectionsDefinitions []data.CollectionDefinition
	collectionsDefinitionsMap map[string]data.CollectionDefinition
	collectionHandlers     map[string]CollectionHandler
	client                     *mongo.Client
}

type MongoCollectionHandler struct {
	db         		*mongo.Database
	collection 		data.CollectionDefinition
}

func NewMongoStorage() Storage {
	return &MongoStorage{
		collectionHandlers: map[string]CollectionHandler{},
		collectionsDefinitionsMap: map[string]data.CollectionDefinition{},
	}
}

func newMongoStorageCollectionHandler(db *mongo.Database, collection data.CollectionDefinition) CollectionHandler {
	return &MongoCollectionHandler{
		db: db,
		collection: collection,
	}
}

func (msc *MongoCollectionHandler) GetItem(itemId string) (interface{}, bool) {
	objId, _ := primitive.ObjectIDFromHex(itemId)
	// fetch item
	res := msc.db.Collection(msc.collection.Name).
		FindOne(
		createContext(),
			bson.M{"_id": objId},
			options.FindOne().SetProjection(bson.D{{"_id", 0}}))

	// check fetching errors
	if res.Err() != nil {
		fmt.Printf("unable to fetch item id %s. err: %s", itemId, res.Err().Error())
		return nil, false
	}

	// decode data
	var itemBson interface{}
	if err := res.Decode(&itemBson); err != nil {
		fmt.Printf("unable to decode DB data for item id %s. err: %s", itemId, err)
		return nil, false
	}

	// convert bson data to go map
	var item map[string]interface{}
	b, _ := bson.Marshal(itemBson)
	bson.Unmarshal(b, &item)
	return item, true
}

func (msc *MongoCollectionHandler) AddItem(item map[string]interface{}) (string, error) {
	res, err := msc.db.Collection(msc.collection.Name).InsertOne(createContext(), item)
	if err != nil {
		fmt.Printf("unable to add new item. err: " + err.Error())
		return "", nil
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (msc *MongoCollectionHandler) UpdateItem(itemId string, newItem map[string]interface{}) error {
	objId, _ := primitive.ObjectIDFromHex(itemId)
	_, err := msc.db.Collection(msc.collection.Name).UpdateOne(createContext(), bson.D{{"_id", objId}}, bson.D{{"$set", newItem}})
	if err != nil {
		fmt.Printf("unable to update item. err: " + err.Error())
		return err
	}
	return nil
}

func (msc *MongoCollectionHandler) DeleteItem(itemId string) error {
	objId, _ := primitive.ObjectIDFromHex(itemId)
	_, err := msc.db.Collection(msc.collection.Name).DeleteOne(createContext(), bson.M{"_id": objId})
	if err != nil {
		fmt.Printf("unable to delete item. err: " + err.Error())
		return err
	}
	return nil
}

func (ms *MongoStorage) Initialize(manifest string) {
	ctx := createContext()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbHost))
	if err != nil {
		panic("unable to connect to Mongo db. " + mongodbHost)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic("Mongo db server not found at " + mongodbHost)
	}
	ms.client = client

	ms.initializeCollectionDefinitions(manifest)
	ms.initializeCollections()
}

func createContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (ms *MongoStorage) GetCollectionDefinitions() []data.CollectionDefinition {
	return ms.collectionsDefinitions
}

func (ms *MongoStorage) GetCollection(collectionName string) (CollectionHandler, error) {
	if collection, ok := ms.collectionsDefinitionsMap[collectionName]; ok {
		collectionHandler, exists := ms.collectionHandlers[collectionName]
		if !exists {
			collectionHandler = newMongoStorageCollectionHandler(ms.client.Database(mongodbName), collection)
			ms.collectionHandlers[collectionName] = collectionHandler
		}
		return collectionHandler, nil
	}
	return nil, fmt.Errorf("collection %s not found", collectionName)
}

func (ms *MongoStorage) initializeCollectionDefinitions(manifest string) {
	err := json.Unmarshal([]byte(manifest), &ms.collectionsDefinitions)
	if err != nil {
		log.Fatal("Unable to parse manifest")
	}
}

func (ms *MongoStorage) initializeCollections() {
	for _, collectionDefinition := range ms.collectionsDefinitions {
		ms.collectionsDefinitionsMap[collectionDefinition.Name] = collectionDefinition
	}
}