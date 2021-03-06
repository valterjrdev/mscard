package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/net/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	_ "ms/card/cmd/api/docs"
	"ms/card/internal/api/handler"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
	"ms/card/pkg/service"
	"ms/card/pkg/telemetry/jaeger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	_ = godotenv.Load()
	server := echo.New()
	server.HideBanner = true

	serviceName := os.Getenv("API_NAME")
	serviceNamespace := os.Getenv("API_NAMESPACE")
	if err := jaeger.Setup(serviceName, serviceNamespace); err != nil {
		server.Logger.Fatalf("telemetry.Jaeger failed with %s\n", err)
	}

	server.Use(middleware.Secure())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.Gzip())

	dsn := os.Getenv("API_DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		server.Logger.Fatalf("gorm.Open() failed with %s\n", err)
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200),
	)

	if err != nil {
		server.Logger.Fatalf("db.Use(dbresolver.Register(dbresolver.Config{}) failed with %s\n", err)
	}

	if err := db.AutoMigrate(&entity.Operation{}, &entity.Account{}, &entity.Transaction{}); err != nil {
		server.Logger.Fatalf("DB.AutoMigrate() failed with %s\n", err)
	}

	accountRepository := repository.NewAccount(server.Logger, db)
	operationRepository := repository.NewOperation(server.Logger, db)
	transactionRepository := repository.NewTransaction(server.Logger, db)

	accountService := service.NewAccount(service.AccountOpts{
		AccountRepository: accountRepository,
	})

	transactionService := service.NewTransaction(service.TransactionOpts{
		Logger:                server.Logger,
		AccountService:        accountService,
		TransactionRepository: transactionRepository,
		AccountRepository:     accountRepository,
		Operation:             operationRepository,
	})

	accountHandler := handler.NewAccount(handler.AccountOpts{
		AccountRepository: accountRepository,
	})

	operationTypeHandler := handler.NewOperation(handler.OperationOpts{
		OperationRepository: operationRepository,
	})

	transactionHandler := handler.NewTransaction(handler.TransactionOpts{
		TransactionService:    transactionService,
		TransactionRepository: transactionRepository,
	})

	server.GET("/docs/*", echoSwagger.WrapHandler)
	server.GET(handler.AccountFindAllPath, accountHandler.FindAll)
	server.GET(handler.AccountFindByIDPath, accountHandler.FindByID)
	server.POST(handler.AccountCreatePath, accountHandler.Create)

	server.GET(handler.OperationFindAllPath, operationTypeHandler.FindAll)
	server.GET(handler.OperationFindByIDPath, operationTypeHandler.FindByID)
	server.POST(handler.OperationCreatePath, operationTypeHandler.Create)

	server.GET(handler.TransactionFindAllPath, transactionHandler.FindAll)
	server.POST(handler.TransactionCreatePath, transactionHandler.Create)

	go func() {
		binding := os.Getenv("API_PORT")
		if err := server.Start(binding); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatalf("server.Start() failed with %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}
