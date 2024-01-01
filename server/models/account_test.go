package models

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
		name      string
		account   Account
		wantExist bool
		wantErr   bool
	}{
		{
			name: "should find an account that exists",
			account: Account{
				Name: "test",
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name: "should not find account that does not exist",
			account: Account{
				Name: "test",
			},
			wantExist: false,
			wantErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var rows *sqlmock.Rows
			if test.wantExist {
				rows = sqlmock.NewRows([]string{"count"}).AddRow(1)
			} else {
				rows = sqlmock.NewRows([]string{"count"}).AddRow(0)
			}
			mock.ExpectQuery("SELECT .+ FROM \"accounts\" .+").WithArgs(test.account.Name).WillReturnRows(rows)

			exists, err := AccountExists(db, test.account.Name)

			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}

			if !exists && test.wantExist {
				t.Errorf("wanted to find %s but did not", test.account.Name)
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
		name       string
		account    Account
		wantInsert bool
		wantErr    bool
	}{
		{
			name: "should insert valid account",
			account: Account{
				Name:      "test",
				Category:  Retirement,
				Class:     Asset,
				TaxBucket: Taxable,
			},
			wantInsert: true,
			wantErr:    false,
		},
		{
			name: "should not insert invalid account",
			account: Account{
				Name: "test",
			},
			wantInsert: false,
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantInsert {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"accounts\" .*").WithArgs(test.account.Name, test.account.Class, test.account.Category, test.account.TaxBucket, AnyTime{}, AnyTime{}, nil).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

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

	mock.ExpectQuery("SELECT .* \"accounts\" WHERE name").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"Name"}).AddRow("test"))
	mock.ExpectQuery("SELECT .* \"account_values\"").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"ID", "AccountName", "Value", "CreatedAt"}).AddRow("1", "test", "0", time.Now()))
	_, err = GetAccountByName(db, "test")
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAccountByClass(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	mock.ExpectQuery("SELECT .* \"accounts\" WHERE class").WithArgs(Asset).WillReturnRows(sqlmock.NewRows([]string{"Name"}).AddRow("test"))
	mock.ExpectQuery("SELECT .* \"account_values\"").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"ID", "AccountName", "Value", "CreatedAt"}).AddRow("1", "test", "0", time.Now()))
	_, err = GetAllAccountsByClass(db, Asset)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateAccount(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name            string
		existingAccount Account
		updatedAccount  Account
		wantErr         bool
	}{
		{
			name: "valid update for existing account",
			existingAccount: Account{
				Name:     "test",
				Category: Cash,
				Class:    Asset,
			},
			updatedAccount: Account{
				Name:     "test",
				Category: HSA,
				Class:    Asset,
			},
			wantErr: false,
		},
		{
			name: "should fail if update is an invalid account",
			existingAccount: Account{
				Name:     "test",
				Category: Cash,
				Class:    Asset,
			},
			updatedAccount: Account{
				Name: "test",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT .* \"accounts\" WHERE name").WithArgs(test.updatedAccount.Name).WillReturnRows(AccountToSQLRow(test.existingAccount))
			mock.ExpectQuery("SELECT .* \"account_values\"").WithArgs(test.updatedAccount.Name).WillReturnRows(sqlmock.NewRows([]string{"ID", "AccountName", "Value", "CreatedAt"}).AddRow("1", test.existingAccount.Name, "0", time.Now()))
			if !test.wantErr {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"accounts\"").WithArgs(test.updatedAccount.Name, test.updatedAccount.Class, test.updatedAccount.Category, AnyTime{}, test.updatedAccount.Name).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}
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
