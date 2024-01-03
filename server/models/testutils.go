//go:build test

package models

import (
	"database/sql/driver"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func randomDollarAmount() float64 {
	val := rand.Float64() * float64(rand.Int63n(10000))
	ratio := math.Pow(10, 2)
	return math.Round(val*ratio) / ratio
}

func CreateMockDatabase() (*gorm.DB, sqlmock.Sqlmock, error) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		return &gorm.DB{}, nil, err
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 conn,
		PreferSimpleProtocol: true,
		WithoutReturning:     true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, nil, err
	}
	return db, mock, err
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var AccountColumns = []string{
	"Name",
	"Class",
	"Category",
	"TaxBucket",
	"CreatedAt",
	"UpdatedAt",
	"DeletedAt",
}

var AccountValuesColumns = []string{
	"ID",
	"AccountName",
	"Value",
	"CreatedAt",
}

func AccountToSQLRow(account Account) *sqlmock.Rows {
	return sqlmock.NewRows(AccountColumns).AddRow(
		account.Name,
		string(account.Class),
		string(account.Category),
		string(account.TaxBucket),
		time.Now(),
		time.Now(),
		time.Now(),
	)
}

func AddAccountToRows(rows *sqlmock.Rows, account Account) *sqlmock.Rows {
	return rows.AddRow(
		account.Name,
		string(account.Class),
		string(account.Category),
		string(account.TaxBucket),
		time.Now(),
		time.Now(),
		time.Now(),
	)
}

func AddRandomAccountValues(rows *sqlmock.Rows, accountName string, numRows int) *sqlmock.Rows {
	for i := 0; i < numRows; i++ {
		rows.AddRow(i, accountName, randomDollarAmount(), time.Now())
	}
	return rows
}

type ExpectedStatement struct {
	statement    string
	args         []driver.Value
	returnRows   *sqlmock.Rows
	returnResult driver.Result
	returnError  error
}

func CreateStatementsAccountExists(name string) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .+ FROM \"accounts\"",
			args: []driver.Value{
				name,
			},
			returnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
		},
	}
}

func CreateStatementsAccountDoesNotExist(name string) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .+ FROM \"accounts\"",
			args: []driver.Value{
				name,
			},
			returnRows: sqlmock.NewRows([]string{"count"}).AddRow(0),
		},
	}
}

func CreateStatementsCreateAccount(account Account) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "INSERT INTO \"accounts\" .*",
			args: []driver.Value{
				account.Name,
				account.Class,
				account.Category,
				account.TaxBucket,
				AnyTime{},
				AnyTime{},
				nil,
			},
			returnResult: sqlmock.NewResult(1, 1),
		},
	}
}

func CreateStatementsGetAllAccountsWithValues(accounts []Account, valuesPerAccount int) []ExpectedStatement {
	accountRows := sqlmock.NewRows(AccountColumns)
	accountValuesRows := sqlmock.NewRows(AccountValuesColumns)
	accountNames := []driver.Value{}
	for _, account := range accounts {
		AddAccountToRows(accountRows, account)
		accountNames = append(accountNames, account.Name)
		AddRandomAccountValues(accountValuesRows, account.Name, valuesPerAccount)
	}

	return []ExpectedStatement{
		{
			statement:  "SELECT (.+) FROM \"accounts\"",
			returnRows: accountRows,
		},
		{
			statement:  "SELECT (.+) FROM \"account_values\"",
			args:       accountNames,
			returnRows: accountValuesRows,
		},
	}
}

func CreateStatementsGetAccountsByNameWithValues(account Account, numValues int) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .* \"accounts\" WHERE name",
			args: []driver.Value{
				account.Name,
			},
			returnRows: AddAccountToRows(sqlmock.NewRows(AccountColumns), account),
		},
		{
			statement: "SELECT .* \"account_values\" WHERE \"account_values\".\"account_name\"",
			args: []driver.Value{
				account.Name,
			},
			returnRows: AddRandomAccountValues(sqlmock.NewRows(AccountValuesColumns), account.Name, numValues),
		},
	}
}

func CreateStatementsGetAccountsByClassWithValues(class string, accounts []Account, valuesPerAccount int) []ExpectedStatement {
	accountRows := sqlmock.NewRows(AccountColumns)
	accountValuesRows := sqlmock.NewRows(AccountValuesColumns)
	accountNames := []driver.Value{}
	for _, account := range accounts {
		accountNames = append(accountNames, account.Name)
		AddAccountToRows(accountRows, account)
		AddRandomAccountValues(accountValuesRows, account.Name, valuesPerAccount)
	}

	return []ExpectedStatement{
		{
			statement: "SELECT .* \"accounts\" WHERE class",
			args: []driver.Value{
				class,
			},
			returnRows: accountRows,
		},
		{
			statement:  "SELECT (.+) FROM \"account_values\"",
			args:       accountNames,
			returnRows: accountValuesRows,
		},
	}
}

func CreateStatementsAccountCannotBeFound(name string) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .* \"accounts\" WHERE name",
			args: []driver.Value{
				name,
			},
			returnRows:  nil,
			returnError: gorm.ErrRecordNotFound,
		},
	}
}

func CreateStatementsUpdateAccount(existingAccount Account, updateArgs []driver.Value) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .* \"accounts\" WHERE name",
			args: []driver.Value{
				existingAccount.Name,
			},
			returnRows: AddAccountToRows(sqlmock.NewRows(AccountColumns), existingAccount),
		},
		{
			statement:    "UPDATE \"accounts\"",
			args:         updateArgs,
			returnResult: sqlmock.NewResult(1, 1),
		},
	}
}

func CreateStatementsDeleteAccount(account Account) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .* \"accounts\" WHERE name",
			args: []driver.Value{
				"test",
			},
			returnRows: AddAccountToRows(sqlmock.NewRows(AccountColumns), account),
		},
		{
			statement: "UPDATE \"accounts\" SET \"deleted_at\"",
			args: []driver.Value{
				AnyTime{},
				account.Name,
			},
			returnResult: sqlmock.NewResult(1, 1),
		},
	}
}

func CreateStatementsCreateAccountValue(av AccountValue) []ExpectedStatement {
	return []ExpectedStatement{
		{
			statement: "SELECT .* FROM \"accounts\"",
			args: []driver.Value{
				av.AccountName,
			},
			returnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
		},
		{
			statement: "INSERT INTO \"account_values\" .*",
			args: []driver.Value{
				av.AccountName,
				av.Value,
				AnyTime{},
			},
			returnResult: sqlmock.NewResult(1, 1),
		},
	}
}

func LoadStatements(mock sqlmock.Sqlmock, statements []ExpectedStatement) {
	for _, statement := range statements {
		if strings.Contains(statement.statement, "SELECT") {
			eq := mock.ExpectQuery(statement.statement).WithArgs(statement.args...)
			if statement.returnRows != nil {
				eq.WillReturnRows(statement.returnRows)
			}
			if statement.returnError != nil {
				eq.WillReturnError(statement.returnError)
			}
		} else {
			mock.ExpectBegin()
			if statement.returnError != nil {
				mock.ExpectRollback()
			} else {
				mock.ExpectExec(statement.statement).WithArgs(statement.args...).WillReturnResult(statement.returnResult)
				mock.ExpectCommit()
			}
		}
	}
}
