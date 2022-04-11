package mongodbstorage

import (
	"context"
	"gobot/internal/anime/animefeeder"
	"gobot/internal/anime/releasestorage"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlfinder"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoEntry struct {
	AnimeTitle string `bson:"animeTitle,omitempty"`
	AnimeUrl   string `bson:"animeUrl,omitempty"`
	AnimeTime  int64  `bson:"animeTime"`
	SubsTitle  string `bson:"subsTitle,omitempty"`
	SubsUrl    string `bson:"subsUrl,omitempty"`
	SubsTime   int64  `bson:"subsTime"`
}

type mongodbstorage struct {
	client                *mongo.Client
	connectionString      string
	releasesUrlConnection string
	database              string
	logger                *zap.SugaredLogger
	cachedLatestRealeases []animefeeder.LatestReleases
}

var _ releasestorage.ReleaseStorage = (*mongodbstorage)(nil)

func (m *mongodbstorage) UpdateStorage(entries []animefeeder.LatestReleases) (newEntries []animefeeder.LatestReleases) {
	for _, entry := range entries {
		found := false
		for _, cachedEntry := range m.cachedLatestRealeases {
			if entry.Equal(cachedEntry) {
				found = true
			}
		}

		if !found {
			newEntries = append(newEntries, entry)
		}
	}

	m.cachedLatestRealeases = append(m.cachedLatestRealeases, newEntries...)
	if newEntries != nil {
		err := m.saveToStorage(newEntries)
		if err != nil {
			m.logger.Error(err)
		}
	}

	return
}

func (m *mongodbstorage) saveToStorage(entries []animefeeder.LatestReleases) error {
	releasesUrlCollection := m.client.Database(m.database).Collection(m.releasesUrlConnection)
	var mongoEntries []MongoEntry
	for _, entry := range entries {
		mongoEntries = append(mongoEntries, MongoEntry{
			AnimeTitle: entry.AnimeUrl.Title,
			AnimeUrl:   entry.AnimeUrl.Url,
			AnimeTime:  entry.AnimeUrl.TimeUpdated.Unix(),
			SubsTitle:  entry.SubsUrl.Title,
			SubsUrl:    entry.SubsUrl.Url,
			SubsTime:   entry.SubsUrl.TimeUpdated.Unix(),
		})
	}

	var newEntriesInterface []interface{}
	for _, entry := range mongoEntries {
		newEntriesInterface = append(newEntriesInterface, bson.D{

			{Key: "animeTitle", Value: entry.AnimeTitle},
			{Key: "animeUrl", Value: entry.AnimeUrl},
			{Key: "animeTime", Value: entry.AnimeTime},

			{Key: "subsTitle", Value: entry.SubsTitle},
			{Key: "subsUrl", Value: entry.SubsUrl},
			{Key: "subsTime", Value: entry.SubsTime},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := releasesUrlCollection.InsertMany(ctx, newEntriesInterface)
	if err != nil {
		return err
	}

	m.logger.Infof("%d new entries added", len(newEntriesInterface))

	return nil
}

func (m *mongodbstorage) readStorage() error {
	releases, err := m.readCollection(m.releasesUrlConnection)
	if err != nil {
		return err
	}

	m.cachedLatestRealeases = nil

	m.cachedLatestRealeases = append(m.cachedLatestRealeases, releases...)

	m.logger.Debug("Mongo storage has been read")
	return nil
}

func (m *mongodbstorage) readCollection(collectionName string) ([]animefeeder.LatestReleases, error) {
	var animeUrlEntries []animefeeder.LatestReleases

	animeUrlCollection := m.client.Database(m.database).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := animeUrlCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var allEntries []MongoEntry
	for cur.Next(ctx) {
		var en MongoEntry
		err := cur.Decode(&en)
		if err != nil {
			return nil, err
		}
		allEntries = append(allEntries, en)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	for _, entry := range allEntries {
		animeUrlEntries = append(animeUrlEntries, animefeeder.LatestReleases{
			Anime: nil,
			AnimeUrl: animeurlfinder.AnimeUrlInfo{
				Title:       entry.AnimeTitle,
				Url:         entry.AnimeUrl,
				TimeUpdated: time.Unix(entry.AnimeTime, 0),
			},
			SubsUrl: animesubs.SubsInfo{
				Title:       entry.SubsTitle,
				Url:         entry.SubsUrl,
				TimeUpdated: time.Unix(entry.SubsTime, 0),
			},
		})
	}

	return animeUrlEntries, nil
}

func NewReleaseStorage(connectionString string, database string, logger *zap.SugaredLogger) (releasestorage.ReleaseStorage, error) {
	mongoClient, err := connectToDatabase(connectionString)
	if err != nil {
		return nil, err
	}
	m := &mongodbstorage{
		connectionString:      connectionString,
		client:                mongoClient,
		releasesUrlConnection: os.Getenv("releasesUrlConnection"),
		database:              database,
		logger:                logger,
	}

	err = m.readStorage()
	if err != nil {
		return nil, err
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
