package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/shopspring/decimal"
)

func TestNewFinanceController(t *testing.T) {
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

func TestMapToSortedList(t *testing.T) {
	now := time.Now()
	tests := []struct {
		want  []NetWorthPoint
		input map[time.Time]decimal.Decimal
	}{
		{
			want: []NetWorthPoint{
				{
					Date:  now,
					Value: decimal.Zero,
				},
			},
			input: map[time.Time]decimal.Decimal{
				now: decimal.Zero,
			},
		},
		{
			want: []NetWorthPoint{
				{
					Date:  now,
					Value: decimal.Zero,
				},
				{
					Date:  now.Add(1 * time.Hour),
					Value: decimal.Zero,
				},
				{
					Date:  now.Add(2 * time.Hour),
					Value: decimal.Zero,
				},
				{
					Date:  now.Add(3 * time.Hour),
					Value: decimal.Zero,
				},
				{
					Date:  now.Add(4 * time.Hour),
					Value: decimal.Zero,
				},
				{
					Date:  now.Add(5 * time.Hour),
					Value: decimal.Zero,
				},
			},
			input: map[time.Time]decimal.Decimal{
				now.Add(1 * time.Hour): decimal.Zero,
				now.Add(5 * time.Hour): decimal.Zero,
				now:                    decimal.Zero,
				now.Add(2 * time.Hour): decimal.Zero,
				now.Add(4 * time.Hour): decimal.Zero,
				now.Add(3 * time.Hour): decimal.Zero,
			},
		},
	}

	for _, test := range tests {
		sorted := mapToSortedList(test.input)
		for i, e := range sorted {
			assert.Equal(t, e.Date, test.want[i].Date)
			assert.Equal(t, e.Value, test.want[i].Value)
		}
	}
}

func TestCreateTimeBuckets(t *testing.T) {
	now := time.Now()
	interval := 24 * time.Hour
	tests := []struct {
		name     string
		want     []time.Time
		accounts []models.Account
	}{
		{
			name: "create 1 buckets",
			want: []time.Time{
				now.Round(interval),
			},
			accounts: []models.Account{
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name: "test2",
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
		{
			name: "create 2 buckets",
			want: []time.Time{
				now.Round(interval),
				now.Add(interval).Round(interval),
			},
			accounts: []models.Account{
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval),
						},
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name: "test2",
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
		{
			name: "create 4 buckets",
			want: []time.Time{
				now.Round(interval),
				now.Add(interval).Round(interval),
				now.Add(interval * 2).Round(interval),
				now.Add(interval * 3).Round(interval),
			},
			accounts: []models.Account{
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval * 2),
						},
					},
				},
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval * 3),
						},
					},
				},
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval * 2),
						},
					},
				},
				{
					Name: "test",
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval),
						},
					},
				},
				{
					Name: "test2",
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buckets := createTimeBuckets(test.accounts, interval)
			assert.Equal(t, len(buckets), len(test.want))
			for i := range test.want {
				if !test.want[i].Equal(buckets[i]) {
					t.Errorf("wanted: %v, got: %v", test.want[i], buckets[i])
				}
				if !test.want[i].Equal(buckets[i]) {
					t.Errorf("wanted: %v, got: %v", test.want[i], buckets[i])
				}
			}
		})
	}
}

func TestFillBuckets(t *testing.T) {
	now := time.Now()
	interval := 24 * time.Hour
	tests := []struct {
		name     string
		want     map[time.Time]decimal.Decimal
		accounts []models.Account
	}{
		{
			name: "base",
			want: map[time.Time]decimal.Decimal{
				now.Add(interval).Round(interval): decimal.NewFromInt(1),
			},
			accounts: []models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now.Add(interval),
						},
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test2",
					Category: models.Cash,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buckets := createTimeBuckets(test.accounts, interval)
			values := make(map[time.Time]decimal.Decimal, len(buckets))
			fillBuckets(test.accounts, values, buckets, interval)

			assert.Equal(t, len(values), len(test.want))
			for i := range test.want {
				if !test.want[i].Equal(values[i]) {
					t.Errorf("wanted: %v, got: %v", test.want[i], values[i])
				}
				if !test.want[i].Equal(values[i]) {
					t.Errorf("wanted: %v, got: %v", test.want[i], values[i])
				}
			}
		})
	}
}

func TestRollup(t *testing.T) {
	now := time.Now()
	interval := 24 * time.Hour
	tests := []struct {
		name     string
		want     []NetWorthPoint
		accounts []models.Account
	}{
		{
			name: "rollup into 1 window",
			want: []NetWorthPoint{
				{
					Date:  now.Round(interval),
					Value: decimal.Zero,
				},
			},
			accounts: []models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test2",
					Category: models.HSA,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test3",
					Category: models.Loan,
					Class:    models.Liability,
					Values: []models.AccountValue{
						{
							AccountName: "test3",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test4",
					Category: models.CreditCard,
					Class:    models.Liability,
					Values: []models.AccountValue{
						{
							AccountName: "test4",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
		{
			name: "rollup into 2 windows",
			want: []NetWorthPoint{
				{
					Date:  now.Round(interval),
					Value: decimal.NewFromInt(1),
				},
				{
					Date:  now.Round(interval).Add(interval),
					Value: decimal.NewFromInt(2),
				},
			},
			accounts: []models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(3).Round(2),
							CreatedAt:   now.Add(interval),
						},
						{
							AccountName: "test",
							Value:       decimal.NewFromInt(2).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test2",
					Category: models.HSA,
					Class:    models.Asset,
					Values: []models.AccountValue{
						{
							AccountName: "test2",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test3",
					Category: models.Loan,
					Class:    models.Liability,
					Values: []models.AccountValue{
						{
							AccountName: "test3",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
				{
					Name:     "test4",
					Category: models.CreditCard,
					Class:    models.Liability,
					Values: []models.AccountValue{
						{
							AccountName: "test4",
							Value:       decimal.NewFromInt(1).Round(2),
							CreatedAt:   now,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := rollup(interval, test.accounts)
			assert.Equal(t, len(values), len(test.want))
			for i := range test.want {
				if !test.want[i].Date.Equal(values[i].Date) {
					t.Errorf("wanted: %v, got: %v", test.want[i].Date, values[i].Date)
				}
				if !test.want[i].Value.Equal(values[i].Value) {
					t.Errorf("wanted: %v, got: %v", test.want[i].Value, values[i].Value)
				}
			}
		})
	}
}

func TestGetNetworthOverTime(t *testing.T) {
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
		responseBody       string
		expectedStatements []models.ExpectedStatement
	}{
		{
			name:         "should successfully return a list of networth over time",
			method:       "GET",
			url:          "/api/networth",
			responseCode: http.StatusOK,
			responseBody: `(\[{"date":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2}","value":".*"}\])+`,
			expectedStatements: models.CreateStatementsGetAllAccountsWithValues([]models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
				},
				{
					Name:     "test2",
					Category: models.HSA,
					Class:    models.Asset,
				},
				{
					Name:     "test3",
					Category: models.Loan,
					Class:    models.Liability,
				},
				{
					Name:     "test4",
					Category: models.CreditCard,
					Class:    models.Liability,
				},
			}, 1),
		},
		{
			name:         "should successfully return a list of networth over time",
			method:       "GET",
			url:          "/api/networth",
			responseCode: http.StatusOK,
			responseBody: `(\[{"date":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2}","value":".*"}\])+`,
			expectedStatements: models.CreateStatementsGetAllAccountsWithValues([]models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
				},
				{
					Name:     "test2",
					Category: models.Cash,
					Class:    models.Asset,
				},
				{
					Name:     "test3",
					Category: models.Cash,
					Class:    models.Liability,
				},
				{
					Name:     "test4",
					Category: models.CreditCard,
					Class:    models.Liability,
				},
			}, 1),
		},
		{
			name:         "should successfully return a list of networth over time",
			method:       "GET",
			url:          "/api/networth",
			responseCode: http.StatusOK,
			responseBody: `(\[{"date":"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2}","value":".*"}\])+`,
			expectedStatements: models.CreateStatementsGetAllAccountsWithValues([]models.Account{
				{
					Name:     "test",
					Category: models.Cash,
					Class:    models.Asset,
				},
				{
					Name:     "test2",
					Category: models.Loan,
					Class:    models.Asset,
				},
				{
					Name:     "test3",
					Category: models.Loan,
					Class:    models.Liability,
				},
				{
					Name:     "test4",
					Category: models.CreditCard,
					Class:    models.Liability,
				},
			}, 1),
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	group := router.Group("/api")
	NewFinanceController(db, group)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			models.LoadStatements(mock, test.expectedStatements)
			req, _ := http.NewRequest(test.method, test.url, test.body)
			router.ServeHTTP(w, req)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			assert.Equal(t, test.responseCode, w.Code)
			assert.MatchRegex(t, w.Body.String(), test.responseBody)
		})
	}
}
