package models

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestCreateAccountValue(t *testing.T) {
	db, mock, err := CreateDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	tests := []struct {
		name               string
		accountValue       AccountValue
		expectedStatements []ExpectedStatement
		wantErr            bool
	}{
		{
			name:    "successfully creates account value",
			wantErr: false,
			accountValue: AccountValue{
				AccountName: "test",
				Value:       1.01,
			},
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* FROM \"accounts\"",
					args: []driver.Value{
						"test",
					},
					returnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
				},
				{
					statement: "INSERT INTO \"account_values\" .*",
					args: []driver.Value{
						"test",
						1.01,
						AnyTime{},
					},
					returnResult: sqlmock.NewResult(1, 1),
				},
			},
		},
		{
			name:    "fails to create value for account that doesn't exist",
			wantErr: true,
			accountValue: AccountValue{
				AccountName: "test",
				Value:       1.01,
			},
			expectedStatements: []ExpectedStatement{
				{
					statement: "SELECT .* FROM \"accounts\"",
					args: []driver.Value{
						"test",
					},
					returnError: gorm.ErrRecordNotFound,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			LoadStatements(mock, test.expectedStatements)
			err := CreateAccountValue(db, test.accountValue)
			if err != nil && !test.wantErr {
				t.Errorf(err.Error())
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
