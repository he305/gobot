package mongodatabase

import (
	"context"
	"gobot/internal/database"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
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

func NewMongoDatabase(connectionString string, database string, animeUrlCollection string, animeSubsCollection string, logger *zap.SugaredLogger) (database.Database, error) {
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

func (md *mongoDatabase) getAnimeEntryByName(collection string, name string) (AnimeEntry, error) {
	filterCursor, err := md.client.Database(md.database).Collection(collection).Find(context.TODO(), bson.M{"title": name})
	if err != nil {
		return AnimeEntry{}, err
	}

	var entries []AnimeEntry
	if err = filterCursor.All(context.TODO(), &entries); err != nil {
		return AnimeEntry{}, err
	}

	if len(entries) == 0 {
		return AnimeEntry{}, nil
	}

	return entries[0], nil
}

func (md *mongoDatabase) addAnimeEntries(collectionString string, entries []AnimeEntry) error {
	if len(entries) == 0 {
		return nil
	}

	collection := md.client.Database(md.database).Collection(collectionString)
	var newEntriesInterface []interface{}
	for _, entry := range entries {
		newEntriesInterface = append(newEntriesInterface, bson.D{
			{Key: "title", Value: entry.Title},
			{Key: "url", Value: entry.Url},
			{Key: "time", Value: entry.Time},
		})
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

func (md *mongoDatabase) AddAnimeUrls(entries ...animeurlservice.AnimeUrlInfo) error {
	var mongoEntries []AnimeEntry
	for _, entry := range entries {
		mongoEntries = append(mongoEntries, AnimeEntry{
			Title: entry.Title,
			Url:   entry.Url,
			Time:  entry.TimeUpdated.Unix(),
		})
	}

	return md.addAnimeEntries(md.animeUrlCollection, mongoEntries)
}

func (md *mongoDatabase) AddSubs(entries ...animesubs.SubsInfo) error {
	var mongoEntries []AnimeEntry
	for _, entry := range entries {
		mongoEntries = append(mongoEntries, AnimeEntry{
			Title: entry.Title,
			Url:   entry.Url,
			Time:  entry.TimeUpdated.Unix(),
		})
	}

	return md.addAnimeEntries(md.animeSubsCollection, mongoEntries)
}

func (md *mongoDatabase) GetAnimeSubByName(name string) (animesubs.SubsInfo, error) {
	entry, err := md.getAnimeEntryByName(md.animeSubsCollection, name)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	return animesubs.SubsInfo{
		Title:       entry.Title,
		Url:         entry.Url,
		TimeUpdated: time.Unix(entry.Time, 0),
	}, nil
}

func (md *mongoDatabase) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	entry, err := md.getAnimeEntryByName(md.animeUrlCollection, name)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	return animeurlservice.AnimeUrlInfo{
		Title:       entry.Title,
		Url:         entry.Url,
		TimeUpdated: time.Unix(entry.Time, 0),
	}, nil
}

var _ database.Database = (*mongoDatabase)(nil)
