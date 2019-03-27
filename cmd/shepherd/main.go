package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/integr8ly/shepherd/pkg/domain/usecases"

	"github.com/integr8ly/shepherd/pkg/chat"

	"github.com/integr8ly/shepherd/pkg/data/bolt"
	shepherdChat "github.com/integr8ly/shepherd/pkg/gchat"
	"github.com/integr8ly/shepherd/pkg/web"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	gchat "google.golang.org/api/chat/v1"
)

var (
	logLevel string
	dbLoc    string
	platform string
)

func main() {
	flag.StringVar(&logLevel, "log-level", "debug", "use this to set log level: error, info, debug")
	flag.StringVar(&dbLoc, "db-loc", "./bot-db", "set the location of the db file")
	flag.StringVar(&platform, "platform", "hangouts", "choose the chat platform to target")
	flag.Parse()
	switch logLevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("log-level set to info")
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Error("log-level set to error")
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("log-level set to debug")
	default:
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Error("log-level set to error")
	}

	_, err := bolt.Connect(dbLoc)
	if err != nil {
		panic(err)
	}
	defer bolt.Disconnect()

	if err := bolt.Setup(); err != nil {
		panic(err)
	}

	chatActionHandler := chat.NewActionHandler()
	router := web.BuildRouter()
	logger := logrus.StandardLogger()

	if platform == "hangouts" {
		// hangout client
		gKey := os.Getenv("GOOGLE_CHAT_KEY")
		gClient, err := google.DefaultClient(context.TODO(), "https://www.googleapis.com/auth/chat.bot")
		if err != nil {
			panic(err)
		}
		gservice, err := gchat.New(gClient)
		if err != nil {
			panic(err)
		}
		spacesService := gchat.NewSpacesService(gservice)
		//TODO change
		handlers := usecases.NewHelpUseCase().Register()

		hangoutChatHandler := shepherdChat.NewActionHandler(spacesService, handlers)

		chatActionHandler.RegisterHandler(hangoutChatHandler)

		handler := web.NewHangoutHandler(chatActionHandler, gKey)
		web.MountHangoutHandler(router, handler)
		// register commands

	}

	httpHandler := web.BuildHTTPHandler(router)
	//sys
	{
		web.MountSystemHandler(router)
	}

	logger.Println("starting api on 8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		logger.Fatal(err)
	}

}
