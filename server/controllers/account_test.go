package controllers

import (
	"bytes"
	"database/sql/driver"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestNewAccountController(t *testing.T) {
	db, _, err := models.CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	group := router.Group("/api")
	controller := NewAccountController(db, group)

	assert.Equal(t, controller.DB, db)
}

func TestCreateOrUpdateAccount(t *testing.T) {
	db, mock, err := models.CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name               string
		method             string
		url                string
		body               io.Reader
		responseCode       int
		expectedStatements []models.ExpectedStatement
	}{
		{
			name:         "should create a valid account that does not yet exist successfully",
			method:       "POST",
			url:          "/api/accounts",
			body:         bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":"cash"}`)),
			responseCode: 200,
			expectedStatements: append(
				models.CreateStatementsAccountDoesNotExist("test"),
				models.CreateStatementsCreateAccount(models.Account{Name: "test", Class: models.Asset, Category: models.Cash})...,
			),
		},
		{
			name:               "should not create an invalid account - missing retirement tax bucket",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":"retirement"}`)),
			responseCode:       400,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing name",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"", "class":"asset", "category":"cash"}`)),
			responseCode:       400,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing class",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"test", "class":"", "category":"cash"}`)),
			responseCode:       400,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing category",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":""}`)),
			responseCode:       400,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:         "should update an existing account",
			method:       "POST",
			url:          "/api/accounts",
			body:         bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":"hsa"}`)),
			responseCode: 200,
			expectedStatements: append(
				models.CreateStatementsAccountExists("test"),
				models.CreateStatementsUpdateAccount(
					models.Account{Name: "test", Class: models.Asset, Category: models.Cash},
					[]driver.Value{
						"test",
						models.Asset,
						models.HSA,
						models.AnyTime{},
						"test",
					},
				)...,
			),
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	group := router.Group("/api")
	NewAccountController(db, group)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			models.LoadStatements(mock, test.expectedStatements)
			req, _ := http.NewRequest(test.method, test.url, test.body)
			router.ServeHTTP(w, req)
			assert.Equal(t, test.responseCode, w.Code)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
