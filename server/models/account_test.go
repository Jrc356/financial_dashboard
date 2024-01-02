package models

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestValidateAccount(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		wantErr bool
	}{
		{
			name: "should not error if account is valid",
			account: Account{
				Name:     "test",
				Category: Cash,
				Class:    Asset,
			},
			wantErr: false,
		},
		{
			name: "should error if name is blank",
			account: Account{
				Name:  "",
				Class: Asset,
			},
			wantErr: true,
		},
		{
			name: "should error if class is blank",
			account: Account{
				Name:  "test",
				Class: "",
			},
			wantErr: true,
		},
		{
			name: "should error if category is unknown",
			account: Account{
				Name:     "test",
				Category: "notreal",
				Class:    Asset,
			},
			wantErr: true,
		},
		{
			name: "should error if no tax bucket is set for a retirement account",
			account: Account{
				Name:      "test",
				Category:  Retirement,
				TaxBucket: "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateAccount(test.account)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
		})
	}
}

func TestAccountExists(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name               string
		wantErr            bool
		wantExist          bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:      "should find an account that exists",
			wantErr:   false,
			wantExist: true,
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .+ FROM \"accounts\"",
					args: []driver.Value{
						"test",
					},
					returnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
				},
			},
		},
		{
			name:      "should not find account that does not exist",
			wantErr:   false,
			wantExist: false,
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\"",
					args: []driver.Value{
						"test",
					},
					returnRows: sqlmock.NewRows([]string{"count"}).AddRow(0),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			exists, err := AccountExists(db, "test")
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if !exists && test.wantExist {
				t.Errorf("wanted exist, got nonexistant")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name               string
		wantErr            bool
		wantCreate         bool
		account            Account
		expectedStatements []ExpectedStatement
	}{
		{
			name:       "should create valid account",
			wantErr:    false,
			wantCreate: true,
			account: Account{
				Name:      "test",
				Category:  Retirement,
				Class:     Asset,
				TaxBucket: Taxable,
			},
			expectedStatements: []ExpectedStatement{
				{
					statement: "INSERT INTO \"accounts\" .*",
					args: []driver.Value{
						"test",
						Asset,
						Retirement,
						Taxable,
						AnyTime{},
						AnyTime{},
						nil,
					},
					returnResult: sqlmock.NewResult(1, 1),
				},
			},
		},
		{
			name: "should not insert invalid account",
			account: Account{
				Name: "test",
			},
			wantErr:            true,
			wantCreate:         false,
			expectedStatements: []ExpectedStatement{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			err := CreateAccount(db, test.account)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllAccounts(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	mock.ExpectQuery("SELECT .* \"accounts\" .*").WillReturnRows(sqlmock.NewRows([]string{"test"}))
	_, err = GetAllAccounts(db)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByName(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := Account{
		Name:     "test",
		Category: Cash,
		Class:    Asset,
	}

	tests := []struct {
		name               string
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:    "should retrieve record if name exists",
			wantErr: false,
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						testAccount.Name,
					},
					returnRows: AccountToSQLRow(testAccount),
				},
				{
					statement: "SELECT .* \"account_values\"",
					args: []driver.Value{
						testAccount.Name,
					},
					returnRows: sqlmock.NewRows([]string{"ID", "AccountName", "Value", "CreatedAt"}).AddRow("1", "test", "0", time.Now()),
				},
			},
		},
		{
			name:    "should error if name does not exist",
			wantErr: true,
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						"test",
					},
					returnRows:  nil,
					returnError: gorm.ErrRecordNotFound,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			_, err = GetAccountByNameWithValues(db, "test")
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
		})
	}

}

func TestGetAccountByClass(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name               string
		class              AccountClass
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:    "should get all accounts by valid class",
			wantErr: false,
			class:   Asset,
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE class",
					args: []driver.Value{
						Asset,
					},
					returnRows: sqlmock.NewRows([]string{"Name"}).AddRow("test"),
				},
				{
					statement: "SELECT .* \"account_values\"",
					args: []driver.Value{
						"test",
					},
					returnRows: sqlmock.NewRows([]string{"ID", "AccountName", "Value", "CreatedAt"}).AddRow("1", "test", "0", time.Now()),
				},
			},
		},
		{
			name:               "should error if class is invalid",
			class:              "NotReal",
			wantErr:            true,
			expectedStatements: []ExpectedStatement{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			_, err = GetAllAccountsByClass(db, test.class)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := Account{
		Name:     "test",
		Category: Cash,
		Class:    Asset,
	}

	tests := []struct {
		name               string
		updatedAccount     Account
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:    "happy path - valid update for existing account",
			wantErr: false,
			updatedAccount: Account{
				Name:     "test",
				Category: HSA,
				Class:    Asset,
			},
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						"test",
					},
					returnRows: AccountToSQLRow(testAccount),
				},
				{
					statement: "UPDATE \"accounts\"",
					args: []driver.Value{
						"test",
						Asset,
						HSA,
						AnyTime{},
						"test",
					},
					returnResult: sqlmock.NewResult(1, 1),
				},
			},
		},
		{
			name:    "should fail if update is an invalid account",
			wantErr: true,
			updatedAccount: Account{
				Name:     "test",
				Class:    Asset,
				Category: Retirement,
			},
			expectedStatements: []ExpectedStatement{},
		},
		{
			name:    "should fail if account account does not exist",
			wantErr: true,
			updatedAccount: Account{
				Name:     "test2",
				Category: HSA,
				Class:    Asset,
			},
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						"test2",
					},
					returnError: gorm.ErrRecordNotFound,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			_, err = UpdateAccount(db, test.updatedAccount.Name, test.updatedAccount)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if err := mock.ExpectationsWereMet(); err != nil && !test.wantErr {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := Account{
		Name:     "test",
		Category: Cash,
		Class:    Asset,
	}

	tests := []struct {
		name               string
		accountName        string
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:        "successfully deletes existing account",
			wantErr:     false,
			accountName: "test",
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						"test",
					},
					returnRows: AccountToSQLRow(testAccount),
				},
				{
					statement: "UPDATE \"accounts\" SET \"deleted_at\"",
					args: []driver.Value{
						AnyTime{},
						"test",
					},
					returnResult: sqlmock.NewResult(1, 1),
				},
			},
		},
		{
			name:        "should fail if account account does not exist",
			wantErr:     true,
			accountName: "test2",
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* \"accounts\" WHERE name",
					args: []driver.Value{
						"test2",
					},
					returnError: gorm.ErrRecordNotFound,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			_, err = DeleteAccount(db, test.accountName)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if err := mock.ExpectationsWereMet(); err != nil && !test.wantErr {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
