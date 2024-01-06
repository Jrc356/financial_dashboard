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
	"github.com/shopspring/decimal"
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
			responseCode: http.StatusCreated,
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
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing name",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"", "class":"asset", "category":"cash"}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing class",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"test", "class":"", "category":"cash"}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an invalid account - missing category",
			method:             "POST",
			url:                "/api/accounts",
			body:               bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":""}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:         "should update an existing account",
			method:       "POST",
			url:          "/api/accounts",
			body:         bytes.NewReader([]byte(`{"name":"test", "class":"asset", "category":"hsa"}`)),
			responseCode: http.StatusOK,
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

func TestGetAccounts(t *testing.T) {
	db, mock, err := models.CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := models.Account{
		Name:     "test",
		Class:    models.Asset,
		Category: models.Cash,
	}

	tests := []struct {
		name               string
		method             string
		url                string
		body               io.Reader
		responseCode       int
		expectedStatements []models.ExpectedStatement
	}{
		{
			name:               "should get all accounts if no queries",
			method:             "GET",
			url:                "/api/accounts",
			responseCode:       http.StatusOK,
			expectedStatements: models.CreateStatementsGetAllAccountsWithValues([]models.Account{testAccount}, 10),
		},
		{
			name:               "should get account by name if name query is specified",
			method:             "GET",
			url:                "/api/accounts?name=test",
			responseCode:       http.StatusOK,
			expectedStatements: models.CreateStatementsGetAccountByNameWithValues(testAccount, 10),
		},
		{
			name:               "should get accounts by class if class query is specified",
			method:             "GET",
			url:                "/api/accounts?class=asset",
			responseCode:       http.StatusOK,
			expectedStatements: models.CreateStatementsGetAccountsByClassWithValues(models.Asset.String(), []models.Account{testAccount}, 10),
		},
		{
			name:               "should return an error if named account does not exist",
			method:             "GET",
			url:                "/api/accounts?name=test",
			responseCode:       http.StatusNotFound,
			expectedStatements: models.CreateStatementsAccountCannotBeFound("test"),
		},
		{
			name:               "should return an error if class does not exist",
			method:             "GET",
			url:                "/api/accounts?class=test",
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
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

func TestDeleteAccounts(t *testing.T) {
	db, mock, err := models.CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := models.Account{
		Name:     "test",
		Class:    models.Asset,
		Category: models.Cash,
	}

	tests := []struct {
		name               string
		method             string
		url                string
		body               io.Reader
		responseCode       int
		expectedStatements []models.ExpectedStatement
	}{
		{
			name:         "should delete an existing account by name",
			method:       "DELETE",
			url:          "/api/accounts?name=test",
			responseCode: http.StatusOK,
			expectedStatements: append(
				models.CreateStatementsAccountExists("test"),
				models.CreateStatementsDeleteAccount(testAccount)...,
			),
		},
		{
			name:               "should return an error if named account does not exist",
			method:             "DELETE",
			url:                "/api/accounts?name=test",
			responseCode:       http.StatusNotFound,
			expectedStatements: models.CreateStatementsAccountCannotBeFound("test"),
		},
		{
			name:               "should return an error if class does not exist",
			method:             "DELETE",
			url:                "/api/accounts?class=test",
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
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

func TestCreateAccountValue(t *testing.T) {
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
			name:         "should create an account value",
			method:       "POST",
			url:          "/api/accounts/value",
			responseCode: http.StatusOK,
			body:         bytes.NewReader([]byte(`{"account_name": "test", "value": 532.01}`)),
			expectedStatements: append(
				models.CreateStatementsAccountExists("test"),
				models.CreateStatementsCreateAccountValue(models.AccountValue{AccountName: "test", Value: decimal.NewFromFloat(532.01)})...,
			),
		},
		{
			name:               "should not create an account value for an account that does not exist",
			method:             "POST",
			url:                "/api/accounts/value",
			body:               bytes.NewReader([]byte(`{"account_name": "test", "value": 8791.43}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: models.CreateStatementsAccountDoesNotExist("test"),
		},
		{
			name:               "should not create an account value with no value provided",
			method:             "POST",
			url:                "/api/accounts/value",
			body:               bytes.NewReader([]byte(`{"account_name": "test"}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
		},
		{
			name:               "should not create an account value if no account_name provided",
			method:             "POST",
			url:                "/api/accounts/value",
			body:               bytes.NewReader([]byte(`{"value": 532.23}`)),
			responseCode:       http.StatusBadRequest,
			expectedStatements: []models.ExpectedStatement{},
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
			println(w.Body.String())
			assert.Equal(t, test.responseCode, w.Code)
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
