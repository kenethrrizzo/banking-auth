package app

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kenethrrizzo/banking-auth/config"
	"github.com/kenethrrizzo/banking-auth/domain"
	"github.com/kenethrrizzo/banking-auth/service"
	log "github.com/sirupsen/logrus"
)

func Start() {
	router := mux.NewRouter()
	serverConfig := config.NewServerConfig()
	dbclient := getDatabaseClient()

	authRepository := domain.NewAuthRepository(dbclient)

	authService := service.NewAuthService(authRepository)

	authHandler := AuthHandler{authService}

	router.HandleFunc(
		"/auth/login", authHandler.login,
	).Methods(http.MethodPost)

	log.Error(http.ListenAndServe(
		fmt.Sprintf("%s:%s",
			serverConfig.Address,
			serverConfig.Port),
		router).Error())
}

func getDatabaseClient() *sqlx.DB {
	dbconfig := config.NewDatabaseConfig()

	client, err := sqlx.Open(
		dbconfig.Driver,
		fmt.Sprintf("%s:%s@/%s",
			dbconfig.Username,
			dbconfig.Password,
			dbconfig.Name))

	if err != nil {
		log.Error("Error while opening connection:", err.Error())
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
