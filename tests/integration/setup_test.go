package integration

import (
	"context"
	"crud-task/internal/container"
	"crud-task/internal/handlers"
	"crud-task/internal/repository"
	"crud-task/internal/routes"
	"crud-task/internal/services"
	"crud-task/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDB *pgxpool.Pool
var testServer *httptest.Server

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")

	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()

	dsn := pkg.MakeDSN()

	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		panic(err)
	}

	testDB = db

	db.Exec(ctx, `TRUNCATE tasks RESTART IDENTITY CASCADE`)

	app := setupApp(db)

	testServer = httptest.NewServer(app)

	m.Run()

	testServer.Close()
	db.Close()
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
