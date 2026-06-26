package integration

import (
	"bytes"
	"context"
	"crud-task/internal/models"
	"crud-task/pkg"
	"encoding/json"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestAllTasks(t *testing.T) {
	dialect := goqu.Dialect("postgres")

	count := 3

	for i := 0; i < count; i++ {
		sql, args, _ := dialect.Insert("tasks").
			Rows(goqu.Record{
				"title":       pkg.RandomString(10),
				"description": pkg.RandomString(100),
				"status":      models.Todo,
				"created_at":  goqu.L("NOW()"),
				"updated_at":  goqu.L("NOW()"),
			}).
			Prepared(true).
			ToSQL()

		_, err := testDB.Exec(context.Background(), sql, args...)

		if err != nil {
			t.Fatalf("failed to insert task: %v", err)
		}
	}

	response, err := http.Get(testServer.URL + "/api/tasks")

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var tasks []models.Task
	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tasks) < count {
		t.Fatalf("expected at least %d tasks, got %d", count, len(tasks))
	}
}

func TestCreateTask(t *testing.T) {
	dialect := goqu.Dialect("postgres")

	task := models.Task{
		Title:       pkg.RandomString(10),
		Description: pkg.RandomString(50),
		Status:      models.Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jsonBody, err := json.Marshal(task)

	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Post(
		testServer.URL+"/api/tasks",
		"application/json",
		bytes.NewReader(jsonBody),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}

	var createdTask models.Task
	json.NewDecoder(response.Body).Decode(&createdTask)

	if createdTask.Id == 0 {
		t.Fatal("id not generated")
	}

	var count int
	sql, args, _ := dialect.From("tasks").
		Select(goqu.COUNT("*").As("count")).
		Prepared(true).
		ToSQL()

	err = testDB.QueryRow(context.Background(), sql, args...).Scan(&count)

	if err != nil {
		t.Fatal("scan error: ", err.Error())
	}

	if count == 0 {
		t.Fatal("db not updated")
	}
}

func TestShowTask(t *testing.T) {
	dialect := goqu.Dialect("postgres")

	var id int

	sql, args, _ := dialect.
		Insert("tasks").
		Rows(goqu.Record{
			"title":       pkg.RandomString(10),
			"description": pkg.RandomString(100),
			"status":      models.Todo,
			"created_at":  goqu.L("NOW()"),
			"updated_at":  goqu.L("NOW()"),
		}).
		Returning("id").
		Prepared(true).
		ToSQL()

	err := testDB.QueryRow(context.Background(), sql, args...).Scan(&id)

	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	response, _ := http.Get(testServer.URL + "/api/tasks/" + strconv.Itoa(id))
	defer response.Body.Close()

	var task models.Task
	json.NewDecoder(response.Body).Decode(&task)

	if task.Id != int64(id) {
		t.Fatal("wrong task")
	}
}

func TestUpdateTask(t *testing.T) {
	dialect := goqu.Dialect("postgres")

	var id int

	sql, args, _ := dialect.
		Insert("tasks").
		Rows(goqu.Record{
			"title":       pkg.RandomString(10),
			"description": pkg.RandomString(100),
			"status":      models.Todo,
			"created_at":  goqu.L("NOW()"),
			"updated_at":  goqu.L("NOW()"),
		}).
		Returning("id").
		Prepared(true).
		ToSQL()

	testDB.QueryRow(context.Background(), sql, args...).Scan(&id)

	body := models.Task{
		Title:       pkg.RandomString(10),
		Description: pkg.RandomString(100),
		Status:      models.InProgress,
	}

	bodyJson, _ := json.Marshal(body)

	request, _ := http.NewRequest(http.MethodPut,
		testServer.URL+"/api/tasks/"+strconv.Itoa(id),
		bytes.NewReader(bodyJson),
	)

	response, _ := http.DefaultClient.Do(request)
	defer response.Body.Close()

	var title string

	sqlFind, argsFind, _ := dialect.From("tasks").
		Select("title").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).
		ToSQL()

	err := testDB.QueryRow(context.Background(), sqlFind, argsFind...).Scan(&title)

	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	if title != body.Title {
		t.Fatal("not updated")
	}
}

func TestDeleteTask(t *testing.T) {
	dialect := goqu.Dialect("postgres")

	var id int

	sql, args, _ := dialect.
		Insert("tasks").
		Rows(goqu.Record{
			"title":       pkg.RandomString(10),
			"description": pkg.RandomString(100),
			"status":      models.Todo,
			"created_at":  goqu.L("NOW()"),
			"updated_at":  goqu.L("NOW()"),
		}).
		Returning("id").
		Prepared(true).
		ToSQL()

	testDB.QueryRow(context.Background(), sql, args...).Scan(&id)

	request, _ := http.NewRequest(http.MethodDelete,
		testServer.URL+"/api/tasks/"+strconv.Itoa(id),
		nil,
	)

	http.DefaultClient.Do(request)

	var count int
	sqlDelete, argsDelete, _ := dialect.
		From("tasks").
		Select(goqu.COUNT("*").As("count")).
		Where(goqu.C("id").Eq(id)).
		Prepared(true).
		ToSQL()

	err := testDB.QueryRow(context.Background(), sqlDelete, argsDelete...).Scan(&count)

	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	if count != 0 {
		t.Fatal("not deleted")
	}
}
