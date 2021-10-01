package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	core "github.com/Qalifah/aboki-africa-assessment"
	"github.com/Qalifah/aboki-africa-assessment/routes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Qalifah/aboki-africa-assessment/config"
	"github.com/Qalifah/aboki-africa-assessment/database/postgres"
	"github.com/Qalifah/aboki-africa-assessment/handler"
	"github.com/dimfeld/httptreemux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var url = "http://localhost:%s"

type TestHandler struct {
	userRepository         		core.UserRepository
	userRefCodeRepository  		core.ReferralCodeRepository
	userReferralRepository 		core.ReferralRepository
	userPointRepository    		core.PointRepository
	userTransactionRepository	core.TransactionRepository
	client                 		*postgres.Client
}

var testHandler *TestHandler

func TestMain(m *testing.M) {
	file, err := os.Open("../config/config.yml")
	if err != nil {
		log.Fatalf("unable to open config file: %v", err)
	}

	cfg := &config.BaseConfig{}
	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	postgresClient, err := postgres.New(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to create postgre client: %v", err)
	}

	userRepo := postgres.NewUserRepository(postgresClient)
	referralCodeRepo := postgres.NewReferralCodeRepository(postgresClient)
	referralRepo := postgres.NewReferralRepository(postgresClient)
	pointsRepo := postgres.NewPointRepository(postgresClient)
	transactionRepo := postgres.NewTransactionRepository(postgresClient)

	h := handler.New(userRepo, referralCodeRepo, referralRepo, pointsRepo, transactionRepo, postgresClient.BeginTx)

	router := httptreemux.New()

	routes.SetupRoutes(router, h)

	url = fmt.Sprintf(url, cfg.ServePort)
	srv := &http.Server{
		Addr:    ":" + cfg.ServePort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("unable to listen: %s", err)
		}
	}()

	// allow the goroutine above start the server
	time.Sleep(time.Second)

	testHandler = &TestHandler{
		userRepository:         userRepo,
		userRefCodeRepository: referralCodeRepo,
		userReferralRepository: referralRepo,
		userPointRepository:    pointsRepo,
		userTransactionRepository: transactionRepo,
		client:                 postgresClient,
	}
	// run the tests
	code := m.Run()

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("unable to shutdown server gracefully: %v", err)
	}

	os.Exit(code)
}

// serialize obj into json bytes
func serialize(obj interface{}) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(obj); err != nil {
		log.Fatalf("unable to serialize obj: %v", err)
	}
	return buf
}

func deleteAllFromTable(name string) error {
	_, err := testHandler.client.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", name))
	return err
}