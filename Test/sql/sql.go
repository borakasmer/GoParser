package sql

import (
	"Test/shared"
	"context"
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"time"
)

func SqlOpen() *sql.DB {
	db, err := sql.Open("sqlserver", shared.Config.SQLURL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetSqlContent(db *sql.DB) ([]string, []float64, []string, []string, error) {
	var (
		CompanyName []string
		CompanyCode []float64
		CompanyUrl  []string
		ApiUrl      []string
		ctx         context.Context
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select CompanyName,CompanyCode,CompanyUrl,ApiUrl from [dbo].[Companies] where IsDeleted = 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var _companyName string
		var _companyCode float64
		var _companyUrl string
		var _apiUrl string
		err := rows.Scan(&_companyName, &_companyCode, &_companyUrl, &_apiUrl)
		if err != nil {
			return CompanyName, CompanyCode, CompanyUrl, ApiUrl, err
		} else {
			CompanyName = append(CompanyName, _companyName)
			CompanyCode = append(CompanyCode, _companyCode)
			CompanyUrl = append(CompanyUrl, _companyUrl)
			ApiUrl = append(ApiUrl, _apiUrl)
		}
	}
	return CompanyName, CompanyCode, CompanyUrl, ApiUrl, nil
}
func InsertSqlContent(db *sql.DB, company *shared.AddCompany) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO Companies(CompanyName,CompanyCode,CompanyUrl,ApiUrl) VALUES (@p1, @p2,@p3,@p4); select ID = convert(bigint, SCOPE_IDENTITY())")
	if err != nil {
		handleError(err, "Could not insert SqlDB")
		return 0, err
	}
	var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	defer stmt.Close()
	rows := stmt.QueryRowContext(ctx, company.CompanyName, company.CompanyCode, company.CompanyUrl, company.ApiUrl)
	if rows.Err() != nil {
		return 0, err
	}
	var _id int64
	rows.Scan(&_id)
	return _id, nil
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}
