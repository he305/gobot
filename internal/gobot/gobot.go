package gobot

import (
	"gobot/internal/anime/animefeeder"
	"gobot/internal/anime/animemessageprovider"
	"gobot/internal/anime/animesubsrepository"
	"gobot/internal/anime/animeurlrepository"
	"gobot/internal/database"
	"gobot/internal/database/filedatabase"
	"gobot/internal/database/mongodatabase"
	"gobot/pkg/animeservice/malv2service"
	"gobot/pkg/animesubs/kitsunekkov2"
	"gobot/pkg/animeurlservice/subspleaserss"
	"gobot/pkg/logging"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	return viper.ReadInConfig()
}

func initLogger() {
	rawJson, err := os.ReadFile("logConfig.json")
	if err != nil {
		logging.InitDefaultLogger()
		return
	}

	logging.InitLoggerConfig(rawJson)
}

func getInfoForPrinting(animeMessageProvider animemessageprovider.AnimeMessageProvider, stChan chan string) {
	st, err := animeMessageProvider.GetMessage()

	if err != nil {
		logger.Errorf("Error getting message from anime message provider, error %v", err)
	}

	stChan <- st
	close(stChan)
}

func createPath(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	_, err = os.Create(path)
	return err
}

func Run() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env file was not found")
	}

	initLogger()
	logger = logging.GetLogger()
	if err := initConfig(); err != nil {
		logger.Panic(err)
	}

	malv2username := os.Getenv("malv2username")
	malv2password := os.Getenv("malv2password")
	telegramChatId, _ := strconv.ParseInt(os.Getenv("telegramChatId"), 10, 64)

	kitsunekkoCachePath := viper.GetString("kitsunekkoCachePath")
	releaseStoragePath := viper.GetString("releaseStoragePath")
	fileStorageFolder := viper.GetString("fileStorageFolder")
	animeUrlCollection := viper.GetString("animeUrlCollection")
	animeSubsCollection := viper.GetString("animeSubsCollection")

	animeUrlStoragePath := fileStorageFolder + animeUrlCollection + ".txt"
	animeSubsStoragePath := fileStorageFolder + animeSubsCollection + ".txt"

	if err := createPath(kitsunekkoCachePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", kitsunekkoCachePath)
	}
	if err := createPath(releaseStoragePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", releaseStoragePath)
	}
	if err := createPath(animeUrlStoragePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", animeUrlStoragePath)
	}
	if err := createPath(animeSubsStoragePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", animeSubsStoragePath)
	}

	malserv := malv2service.NewMalv2Service(malv2username, malv2password)
	//fileIo := fileio.NewDefaultFileIO()
	//kitsunekkoSubService := kitsunekko.NewKitsunekkoScrapper(fileIo, kitsunekkoCachePath, 3*time.Minute)
	kitsunekkoSubService := kitsunekkov2.NewKitsunekkoScrapperV2(3*time.Minute, logger)
	subspleaserss := subspleaserss.NewSubsPleaseRss(subspleaserss.Rss1080Url, 3*time.Minute, logger)

	var database database.Database
	database, err := mongodatabase.NewMongoDatabase(
		os.Getenv("MONGODB_CONNECTION"),
		"anime_releases",
		animeUrlCollection,
		animeSubsCollection,
		logger,
	)

	if err != nil {
		logger.Errorf("Couldn't connect to mongo db, switching to file storage, error: %v", err)
		database = filedatabase.NewFileDatabase(
			animeUrlStoragePath,
			animeSubsStoragePath,
			logger,
		)
	}

	animeUrlRepo := animeurlrepository.NewAnimeUrlRepository(database)
	animeSubsRepo := animesubsrepository.NewAnimeSubsRepository(database)

	animeFeeder := animefeeder.NewAnimeFeeder(malserv, kitsunekkoSubService, subspleaserss, animeUrlRepo, animeSubsRepo, logger)

	animeMessageProvider := animemessageprovider.NewAnimeMessageProvider(animeFeeder)

	debugMode := viper.GetBool("debugMode")
	telegramToken := os.Getenv("telegramToken")

	bot, err := tgbot.NewBotAPI(telegramToken)
	if err != nil {
		logger.Errorf("Couldn't initialize telegram bot")
	}

	bot.Debug = debugMode

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	logger.Infow("Telegram bot started")

	go func() {
		for {
			stChan := make(chan string)
			go getInfoForPrinting(animeMessageProvider, stChan)

			st := <-stChan

			if st != "" {
				msg := tgbot.NewMessage(telegramChatId, st)
				bot.Send(msg)
			}

			time.Sleep(time.Minute * 3)
		}
	}()

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {

			// Just for debug
			charArray := []byte(update.Message.Text)
			if charArray[0] == '/' {
				splittedMessage := strings.Split(update.Message.Text, " ")
				if len(splittedMessage) > 1 && splittedMessage[0] == "/anime" {
					entry, err := malserv.GetAnimeByTitle(strings.Join(splittedMessage[1:], " "))
					var msg tgbot.MessageConfig
					if err != nil {
						msg = tgbot.NewMessage(update.Message.Chat.ID, "Error getting anime")
					} else {
						msg = tgbot.NewMessage(update.Message.Chat.ID, entry.VerboseOutput())
					}
					bot.Send(msg)
				}
			}
		}
	}
}
