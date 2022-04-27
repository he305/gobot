package mongodatabase

import (
	"context"
	"gobot/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoDatabase struct {
	connectionString    string
	client              *mongo.Client
	database            string
	animeUrlCollection  string
	animeSubsCollection string
	logger              *zap.SugaredLogger
}

func makeEntryInterface(entry map[string]interface{}) bson.D {
	var bsonD bson.D
	for k, v := range entry {
		bsonD = append(bsonD, bson.E{
			Key: k, Value: v,
		})
	}
	return bsonD
}

// AddEntries implements database.Database
func (md *mongoDatabase) AddEntries(collectionName string, entries []map[string]interface{}) error {
	if len(entries) == 0 {
		return nil
	}

	collection := md.client.Database(md.database).Collection(collectionName)
	newEntriesInterface := make([]interface{}, 0, len(entries))
	for _, entry := range entries {
		newEntriesInterface = append(newEntriesInterface, makeEntryInterface(entry))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, newEntriesInterface)
	if err != nil {
		return err
	}

	md.logger.Infof("%d new entries added", len(newEntriesInterface))

	return nil
}

// AddEntry implements database.Database
func (md *mongoDatabase) AddEntry(collectionName string, entry map[string]interface{}) error {
	return md.AddEntries(collectionName, []map[string]interface{}{entry})
}

// GetEntryByName implements database.Database
func (md *mongoDatabase) GetEntryByName(collectionName string, key string, name string) (map[string]interface{}, error) {
	filterCursor, err := md.client.Database(md.database).Collection(collectionName).Find(context.TODO(), bson.M{key: name})
	if err != nil {
		return nil, err
	}

	var entries []map[string]interface{}
	if err = filterCursor.All(context.TODO(), &entries); err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, nil
	}

	return entries[0], nil
}

func NewMongoDatabase(connectionString string, database string, animeUrlCollection string, animeSubsCollection string, logger *zap.SugaredLogger) (*mongoDatabase, error) {
	mongoClient, err := connectToDatabase(connectionString)
	if err != nil {
		return nil, err
	}
	m := &mongoDatabase{
		connectionString:    connectionString,
		client:              mongoClient,
		animeUrlCollection:  animeUrlCollection,
		animeSubsCollection: animeSubsCollection,
		database:            database,
		logger:              logger,
	}

	return m, nil
}

func connectToDatabase(connectionString string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

var _ database.Database = (*mongoDatabase)(nil)
