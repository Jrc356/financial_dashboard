package models

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDatabase() (*gorm.DB, sqlmock.Sqlmock, error) {
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

func AccountToSQLRow(account Account) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"Name",
		"Class",
		"Category",
		"TaxBucket",
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
	}).AddRow(account.Name,
		string(account.Class),
		string(account.Category),
		string(account.TaxBucket),
		time.Now(),
		time.Now(),
		time.Now(),
	)
}

type ExpectedStatement struct {
	statement    string
	args         []driver.Value
	returnRows   *sqlmock.Rows
	returnResult driver.Result
	returnError  error
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
