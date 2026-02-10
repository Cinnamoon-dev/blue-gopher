package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Cinnamoon-dev/blue-gopher/internal/database"
	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
	"github.com/Cinnamoon-dev/blue-gopher/internal/repositories"
	"github.com/Cinnamoon-dev/blue-gopher/internal/services"
	_ "github.com/mattn/go-sqlite3"
)

var (
	userHandler *handlers.UserHandler
)

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", "../test.db")
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer db.Close()

	database.CreateTables("../internal/database/tables.sql", db)
	database.Populate("../internal/database/rules.sql", db)

	roleRepository := repositories.NewRoleRepository(db)
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, roleRepository)
	h := handlers.NewUserHandler(userService)
	userHandler = &h

	code := m.Run()
	os.Exit(code)
}

// Vou fazer os testes usando os handlers
// Vou fazer um for com cada caso dos métodos de handler

// Create User
// Create user with a valid email, password and role
// Create user with a repeated email
// Create user with an invalid email
// Create user with an invalid password
func TestCreateUser(t *testing.T) {
	var tests = []struct {
		Message    string
		Email      string
		Password   string
		RoleID     int64
		StatusCode int
	}{
		{"valid email and password should work", "valid@email.com", "valid password", 1, http.StatusOK},
		{"repeated email should fail", "valid@email.com", "repeated email password", 1, http.StatusUnprocessableEntity},
		{"invalid email should fail", "", "invalid email password", 1, http.StatusBadRequest},
		{"invalid password should fail", "invalidpassword@mail.com", "", 1, http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.Message, func(t *testing.T) {
			user := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
				RoleID   int64  `json:"role_id"`
			}{
				Email:    test.Email,
				Password: test.Password,
				RoleID:   test.RoleID,
			}

			bodyBytes, err := json.Marshal(user)
			if err != nil {
				t.Fatalf("failed to marshal user; %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/user", io.NopCloser(bytes.NewBuffer(bodyBytes)))
			rec := httptest.NewRecorder()

			userHandler.CreateUser(rec, req)
			res := rec.Result()
			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			if res.StatusCode != test.StatusCode {
				t.Errorf("expected %d, got %d", test.StatusCode, res.StatusCode)
			}

			t.Logf("response: %s", body)
		})
	}

	userHandler.Svc.RoleRepo.DB.Exec("DELETE FROM usuarios WHERE id != 1")
}
