package server

import (
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	mongo2 "github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/repository/mongo"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"log"
	"net/http"
	"time"
)

type App struct {
	httpServer *http.Server

	authUseCase auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := mongo2.NewUserRepository(db, viper.GetString("mongo.collection"))
	authUseCase := usecase.NewAuthorizer(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl")*time.Second,
	)

	return &App{
		authUseCase: authUseCase,
	}
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
