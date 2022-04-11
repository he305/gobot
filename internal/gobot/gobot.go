package gobot

import (
	"fmt"
	"gobot/internal/anime/animefeeder"
	"gobot/internal/anime/releasestorage"
	"gobot/internal/anime/releasestorage/filereleasestorage"
	"gobot/internal/anime/releasestorage/mongodbstorage"
	"gobot/pkg/animeservice"
	"gobot/pkg/animeservice/malv2service"
	"gobot/pkg/animesubs/kitsunekko"
	"gobot/pkg/animeurlfinder/subspleaserss"
	"gobot/pkg/fileio"
	"gobot/pkg/logging"
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

func getInfoForPrinting(animeFeeder animefeeder.AnimeFeeder, storage releasestorage.ReleaseStorage) (st string) {

	missingInCached, missingInNew, err := animeFeeder.UpdateList()
	if err != nil {
		logger.Errorf("Feeder couldn't update list, error %v", err)
	}
	if missingInCached != nil {
		st += "New entries in list\n"
		for _, v := range missingInCached {
			st += v.VerboseOutput()
			st += "\n"
		}
	}

	if missingInNew != nil {
		st += "Entries were deleted\n"
		for _, v := range missingInNew {
			v.ListStatus = animeservice.NotInList
			st += v.VerboseOutput()
			st += "\n"
		}
	}

	latestReleases := animeFeeder.FindLatestReleases()
	newReleases := storage.UpdateStorage(latestReleases)

	for _, v := range newReleases {
		if v.AnimeUrl.Url != "" {
			st += fmt.Sprintf("New release for %s\nanime url: %s\n", v.Anime.Title, v.AnimeUrl.Url)
		}
		if v.SubsUrl.Url != "" {
			st += fmt.Sprintf("New subs for %s\nurl: %s\n", v.Anime.Title, v.SubsUrl.Url)
		}
	}

	return st
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
		logger.Info(".env file was not found")
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

	if err := createPath(kitsunekkoCachePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", kitsunekkoCachePath)
	}
	if err := createPath(releaseStoragePath); err != nil {
		logger.Panicf("Couldn't create path %s, fatal error", releaseStoragePath)
	}

	malserv := malv2service.NewMalv2Service(malv2username, malv2password)
	fileIo := fileio.NewDefaultFileIO()
	kitsunekkoSubService := kitsunekko.NewKitsunekkoScrapper(fileIo, kitsunekkoCachePath, 5*time.Minute)
	subspleaserss := subspleaserss.NewSubsPleaseRss(subspleaserss.Rss1080Url, 5*time.Minute, logger)

	storage, err := mongodbstorage.NewReleaseStorage(os.Getenv("MONGODB_CONNECTION"), "anime_releases", logger)
	if err != nil {
		logger.Error(err)
		logger.Info("Using file storage")
		storage = filereleasestorage.NewFileReleaseStorage(releaseStoragePath)
	}

	animeFeeder := animefeeder.NewAnimeFeeder(malserv, kitsunekkoSubService, subspleaserss, logger)

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
			st := getInfoForPrinting(animeFeeder, storage)

			if st != "" {
				msg := tgbot.NewMessage(telegramChatId, st)
				bot.Send(msg)
			}

			time.Sleep(3 * time.Minute)
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
