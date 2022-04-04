package gobot

import (
	"gobot/pkg/animeservice/malv2service"
	"gobot/pkg/logging"
	"os"

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

func Run() {
	initLogger()
	logger = logging.GetLogger()
	if err := initConfig(); err != nil {
		logger.Panic(err)
	}

	malv2username := viper.GetString("malv2username")
	malv2password := viper.GetString("malv2password")

	malserv := malv2service.NewMalv2Service(malv2username, malv2password)
	malserv.GetUserAnimeList()

	// debugMode := viper.GetBool("debugMode")
	// telegramToken := viper.GetString("telegramToken")

	// bot, err := tgbot.NewBotAPI(telegramToken)
	// if err != nil {
	// 	logger.Panic()
	// }

	// bot.Debug = debugMode

	// u := tgbot.NewUpdate(0)
	// u.Timeout = 60

	// //logger.Infow("Telegram bot started")

	// updates := bot.GetUpdatesChan(u)
	// for update := range updates {
	// 	if update.Message != nil {
	// 		msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// 		msg.ReplyToMessageID = update.Message.MessageID
	// 		bot.Send(msg)
	// 	}
	// }
}
