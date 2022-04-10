package mongodbstorage

import (
	"context"
	"fmt"
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
	Title string `bson:"title"`
	Url   string `bson:"url"`
	Time  int64  `bson:"time"`
}

type mongodbstorage struct {
	client                *mongo.Client
	connectionString      string
	animeUrlCollection    string
	subsInfoCollection    string
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
		fmt.Println(newEntries)
		err := m.saveToStorage(newEntries)
		if err != nil {
			m.logger.Error(err)
		}
	}

	return
}

func (m *mongodbstorage) saveToStorage(entries []animefeeder.LatestReleases) error {
	animeUrlCollection := m.client.Database(m.database).Collection(m.animeUrlCollection)
	subsUrlCollection := m.client.Database(m.database).Collection(m.subsInfoCollection)

	var newAnimeUrlEntries []MongoEntry
	var newSubsUrlEntries []MongoEntry

	for _, entry := range entries {
		if entry.AnimeUrl.Url != "" {
			newAnimeUrlEntries = append(newAnimeUrlEntries, MongoEntry{
				Title: entry.AnimeUrl.Title,
				Url:   entry.AnimeUrl.Url,
				Time:  entry.AnimeUrl.TimeUpdated.Unix(),
			})
		}

		if entry.SubsUrl.Url != "" {
			newSubsUrlEntries = append(newSubsUrlEntries, MongoEntry{
				Title: entry.SubsUrl.Title,
				Url:   entry.SubsUrl.Url,
				Time:  entry.SubsUrl.TimeUpdated.Unix(),
			})
		}
	}

	// TOOD
	var newAnimeUrlInterface []interface{}
	for _, v := range newAnimeUrlEntries {
		newAnimeUrlInterface = append(newAnimeUrlInterface, bson.D{
			{Key: "title", Value: v.Title},
			{Key: "url", Value: v.Url},
			{Key: "time", Value: v.Time},
		})
	}
	var newSubsUrlInterface []interface{}
	for _, v := range newSubsUrlEntries {
		newSubsUrlInterface = append(newSubsUrlInterface, bson.D{
			{Key: "title", Value: v.Title},
			{Key: "url", Value: v.Url},
			{Key: "time", Value: v.Time},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := animeUrlCollection.InsertMany(ctx, newAnimeUrlInterface)
	if err != nil {
		return err
	}

	_, err = subsUrlCollection.InsertMany(ctx, newSubsUrlInterface)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongodbstorage) readStorage() error {
	animeUrlReleases, err := m.readCollection(m.animeUrlCollection)
	if err != nil {
		return err
	}

	subsUrlReleases, err := m.readCollection(m.subsInfoCollection)
	if err != nil {
		return err
	}

	m.cachedLatestRealeases = nil

	m.cachedLatestRealeases = append(m.cachedLatestRealeases, animeUrlReleases...)
	m.cachedLatestRealeases = append(m.cachedLatestRealeases, subsUrlReleases...)

	fmt.Println(m.cachedLatestRealeases)
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
		animeUrlEntries = append(animeUrlEntries, m.createLatestRelease(collectionName, entry))
	}

	return animeUrlEntries, nil
}

func (m *mongodbstorage) createLatestRelease(collectionName string, entry MongoEntry) animefeeder.LatestReleases {
	switch collectionName {
	case m.animeUrlCollection:
		return animefeeder.LatestReleases{
			Anime: nil,
			AnimeUrl: animeurlfinder.AnimeUrlInfo{
				Title:       entry.Title,
				Url:         entry.Url,
				TimeUpdated: time.Unix(entry.Time, 0),
			},
		}
	case m.subsInfoCollection:
		return animefeeder.LatestReleases{
			Anime: nil,
			SubsUrl: animesubs.SubsInfo{
				Title:       entry.Title,
				Url:         entry.Url,
				TimeUpdated: time.Unix(entry.Time, 0),
			},
		}
	}

	return animefeeder.LatestReleases{}
}

func NewReleaseStorage(connectionString string, database string, logger *zap.SugaredLogger) (releasestorage.ReleaseStorage, error) {
	mongoClient, err := connectToDatabase(connectionString)
	if err != nil {
		return nil, err
	}
	m := &mongodbstorage{
		connectionString:   connectionString,
		client:             mongoClient,
		animeUrlCollection: os.Getenv("animeUrlCollection"),
		subsInfoCollection: os.Getenv("subsInfoCollection"),
		database:           database,
		logger:             logger,
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
