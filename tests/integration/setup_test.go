package integration

import (
	"context"
	"crud-task/internal/container"
	"crud-task/internal/handlers"
	"crud-task/internal/repository"
	"crud-task/internal/routes"
	"crud-task/internal/services"
	"crud-task/pkg"
	"crud-task/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testDB *pgxpool.Pool
var testServer *httptest.Server

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")

	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		log.Fatal(err.Error())
	}

	if err := os.Setenv("APP_ENV", "development"); err != nil {
		log.Fatal(err.Error())
	}

	closeLogs, err := logger.Init()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	ctx := context.Background()
	dsn := pkg.MakeDSN()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		closeLogs()
		panic(err)
	}
	testDB = db
	//
	//if _, err := db.Exec(ctx, `TRUNCATE tasks RESTART IDENTITY CASCADE`); err != nil {
	//	closeLogs()
	//	db.Close()
	//	log.Fatalf("failed to truncate tasks: %v", err)
	//}

	app := setupApp(db)
	testServer = httptest.NewServer(app)

	code := m.Run()

	testServer.Close()
	db.Close()
	closeLogs()

	os.Exit(code)
}

func setupApp(dbPool *pgxpool.Pool) http.Handler {
	baseRepository := repository.NewBaseRepository(dbPool)

	taskRepository := repository.NewTaskRepository(baseRepository)
	taskService := services.NewTaskService(taskRepository)
	taskHandler := handlers.NewTaskHandler(taskService)

	return routes.Routes(&container.Container{
		TaskHandler: taskHandler,
	})
}
