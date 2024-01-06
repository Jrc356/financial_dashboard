package models

import (
	"database/sql/driver"
	"testing"

	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
)

func TestValidateAccount(t *testing.T) {
	tests := []struct {
		name    string
		account Account
		wantErr bool
	}{
		{
			name: "account is valid",
			account: Account{
				Name:      "test",
				Category:  Retirement,
				Class:     Asset,
				TaxBucket: Taxable,
			},
			wantErr: false,
		},
		{
			name: "should error if name is blank",
			account: Account{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "should error if class is blank",
			account: Account{
				Name:     "test",
				Class:    "",
				Category: Cash,
			},
			wantErr: true,
		},
		{
			name: "should error if class is invalid",
			account: Account{
				Name:     "test",
				Class:    "notreal",
				Category: Cash,
			},
			wantErr: true,
		},
		{
			name: "should error if category is blank",
			account: Account{
				Name:     "test",
				Class:    Asset,
				Category: "",
			},
			wantErr: true,
		},
		{
			name: "should error if category is invalid",
			account: Account{
				Name:     "test",
				Class:    Asset,
				Category: "notreal",
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
		{
			name: "should error if tax bucket is invalid",
			account: Account{
				Name:      "test",
				Category:  Retirement,
				Class:     Asset,
				TaxBucket: "notreal",
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
	db, mock, err := CreateMockDatabase()
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
			name:               "should find an account that exists",
			wantErr:            false,
			wantExist:          true,
			expectedStatements: CreateStatementsAccountExists("test"),
		},
		{
			name:               "should not find account that does not exist",
			wantErr:            false,
			wantExist:          false,
			expectedStatements: CreateStatementsAccountDoesNotExist("test"),
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
	db, mock, err := CreateMockDatabase()
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
			expectedStatements: CreateStatementsCreateAccount(Account{
				Name:      "test",
				Category:  Retirement,
				Class:     Asset,
				TaxBucket: Taxable,
			}),
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

func TestGetAllAccountsWithValues(t *testing.T) {
	db, mock, err := CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccounts := []Account{
		{
			Name:     "test",
			Class:    Asset,
			Category: Cash,
		},
		{
			Name:     "test2",
			Class:    Liability,
			Category: Loan,
		},
	}
	numValues := 10
	LoadStatements(mock, CreateStatementsGetAllAccountsWithValues(testAccounts, numValues))
	resp, err := GetAllAccountsWithValues(db)
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, len(testAccounts), len(resp))

	for _, account := range resp {
		assert.Equal(t, numValues, len(account.Values))
	}
}

func TestGetAccountByNameWithValues(t *testing.T) {
	db, mock, err := CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccount := Account{
		Name:     "test",
		Class:    Asset,
		Category: Cash,
	}
	numValues := 10

	tests := []struct {
		name               string
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:               "should retrieve record if name exists",
			wantErr:            false,
			expectedStatements: CreateStatementsGetAccountByNameWithValues(testAccount, numValues),
		},
		{
			name:               "should error if name does not exist",
			wantErr:            true,
			expectedStatements: CreateStatementsAccountCannotBeFound("test"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			account, err := GetAccountByNameWithValues(db, "test")
			if test.wantErr && (err == nil || err != gorm.ErrRecordNotFound) {
				t.Errorf(err.Error())
			}
			if !test.wantErr {
				assert.Equal(t, testAccount.Name, account.Name)
				assert.Equal(t, numValues, len(account.Values))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAccountByClassWithValues(t *testing.T) {
	db, mock, err := CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	accountClass := Asset
	numValues := 10
	testAccounts := []Account{
		{
			Name:     "test",
			Class:    accountClass,
			Category: Cash,
		},
		{
			Name:     "test2",
			Class:    accountClass,
			Category: Cash,
		},
		{
			Name:     "test3",
			Class:    accountClass,
			Category: Cash,
		},
	}

	tests := []struct {
		name               string
		class              AccountClass
		wantErr            bool
		expectedStatements []ExpectedStatement
	}{
		{
			name:               "should get all accounts by valid class",
			wantErr:            false,
			class:              Asset,
			expectedStatements: CreateStatementsGetAccountsByClassWithValues(accountClass.String(), testAccounts, numValues),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			accounts, err := GetAllAccountsByClassWithValues(db, test.class)
			if test.wantErr && (err == nil || err != gorm.ErrRecordNotFound) {
				t.Errorf(err.Error())
			}

			if !test.wantErr {
				assert.Equal(t, len(testAccounts), len(accounts))
				for _, account := range accounts {
					assert.Equal(t, numValues, len(account.Values))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	db, mock, err := CreateMockDatabase()
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
			name:    "should perform a valid update for existing account",
			wantErr: false,
			updatedAccount: Account{
				Name:     "test",
				Category: HSA,
				Class:    Asset,
			},
			expectedStatements: CreateStatementsUpdateAccount(testAccount, []driver.Value{
				"test",
				Asset,
				HSA,
				AnyTime{},
				"test",
			}),
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
			name:    "should fail if account does not exist",
			wantErr: true,
			updatedAccount: Account{
				Name:     "test2",
				Category: HSA,
				Class:    Asset,
			},
			expectedStatements: CreateStatementsAccountCannotBeFound("test2"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			update, err := UpdateAccount(db, test.updatedAccount.Name, test.updatedAccount)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if !test.wantErr {
				assert.Equal(t, testAccount.Name, update.Name)
			}
			if err := mock.ExpectationsWereMet(); err != nil && !test.wantErr {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	db, mock, err := CreateMockDatabase()
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
			name:               "successfully deletes existing account",
			wantErr:            false,
			accountName:        "test",
			expectedStatements: CreateStatementsDeleteAccount(testAccount),
		},
		{
			name:        "should fail if account does not exist",
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
			account, err := DeleteAccount(db, test.accountName)
			if test.wantErr && (err == nil || err != gorm.ErrRecordNotFound) {
				t.Errorf(err.Error())
			}
			if !test.wantErr {
				assert.NotEqual(t, "", account.DeletedAt)
			}
			if err := mock.ExpectationsWereMet(); err != nil && !test.wantErr {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
