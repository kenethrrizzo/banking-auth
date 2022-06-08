package app

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	handl "github.com/kenethrrizzo/banking-auth/app/handlers"
	"github.com/kenethrrizzo/banking-auth/domain"
	repo "github.com/kenethrrizzo/banking-auth/domain/repositories"
	"github.com/kenethrrizzo/banking-auth/service"
	log "github.com/sirupsen/logrus"
)

func Start() {
	router := mux.NewRouter()
	serverConfig := NewServerConfig()
	dbclient := getDatabaseClient()

	authRepository := repo.NewAuthRepository(dbclient)
	rolePermissions := domain.GetRolePermissions()

	authService := service.NewAuthService(authRepository, rolePermissions)

	authHandler := handl.AuthHandler{
		Service: authService,
	}

	router.HandleFunc(
		"/auth/login", authHandler.Login,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/auth/verify", authHandler.Verify,
	).Methods(http.MethodGet)

	log.Error(http.ListenAndServe(
		fmt.Sprintf("%s:%s",
			serverConfig.Address,
			serverConfig.Port),
		router).Error())
}

func getDatabaseClient() *sqlx.DB {
	dbconfig := NewDatabaseConfig()

	client, err := sqlx.Open(
		dbconfig.Driver,
		fmt.Sprintf("%s:%s@/%s",
			dbconfig.Username,
			dbconfig.Password,
			dbconfig.Name))

	if err != nil {
		log.Error("Error while opening connection: ", err.Error())
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
