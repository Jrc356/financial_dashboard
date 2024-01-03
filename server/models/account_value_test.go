package models

import (
	"testing"
)

func TestCreateAccountValue(t *testing.T) {
	db, mock, err := CreateMockDatabase()
	if err != nil {
		t.Errorf(err.Error())
	}
	d, _ := db.DB()
	defer d.Close()

	testAccountValue := AccountValue{
		AccountName: "test",
		Value:       1.01,
	}

	tests := []struct {
		name               string
		accountValue       AccountValue
		expectedStatements []ExpectedStatement
		wantErr            bool
	}{
		{
			name:               "successfully creates account value",
			wantErr:            false,
			accountValue:       testAccountValue,
			expectedStatements: CreateStatementsCreateAccountValue(testAccountValue),
		},
		{
			name:    "fails to create value for account that doesn't exist",
			wantErr: true,
			accountValue: AccountValue{
				AccountName: "test",
				Value:       1.01,
			},
			expectedStatements: CreateStatementsAccountCannotBeFound(testAccountValue.AccountName),
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
