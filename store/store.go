package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "10.138.0.5"
	database = "linebot01"
	user     = "linebot01"
	password = "P@ssw0rd"
)

type Kaiin_host struct {
	No        int64
	User_id   string
	User_name string
}

type Bodymanagement struct {
	No       int64
	User_id  string
	Weight   float64
	Height   float64
	Now_date string
}

type Cashflow struct {
	No             int64
	User_id        string
	Label          string
	Money          int64
	Pauchased_item string
	Register_date  string
}

type Messagerecord struct {
	No         int64
	User_id    string
	Rcvmessage string
	Rcvdate    string
}

func errorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Connectdb() *sql.DB {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)

	con, err := sql.Open("mysql", connectionString)
	errorCheck(err)
	//defer con.Close()

	err = con.Ping()
	errorCheck(err)
	fmt.Println("Successfully created connection to database")

	return con
}

func Insert_kaiin_host(user_id string, user_name string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	//データをインサート
	sqlStatement, err := con.Prepare("INSERT INTO kaiin_host (user_id,user_name) VALUES(?,?)")
	errorCheck(err)
	defer sqlStatement.Close()
	result, err := sqlStatement.Exec(user_id, user_name)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
}

func Insert_bodymanagement(user_id string, weight string, height string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	var timezoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timezoneJST
	time.LoadLocation("Asia/Tokyo")
	date := time.Now()
	//データをインサート
	sqlStatement, err := con.Prepare("INSERT INTO bodymanagement (user_id,weight,height,now_date) VALUES(?,?,?,?)")
	errorCheck(err)
	defer sqlStatement.Close()
	result, err := sqlStatement.Exec(user_id, weight, height, date)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
}

func Insert_cashflow(user_id *string, label string, money string, pauchased_item string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	var timezoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timezoneJST
	time.LoadLocation("Asia/Tokyo")
	date := time.Now()
	//データをインサート
	sqlStatement, err := con.Prepare("INSERT INTO cashflow (user_id,label,money,pauchaed_item,register_date) VALUES(?,?,?,?,?)")
	errorCheck(err)
	defer sqlStatement.Close()
	result, err := sqlStatement.Exec(user_id, label, money, pauchased_item, date)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
}

func Insert_messagerecord(user_id string, rcvmessage string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	var timezoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timezoneJST
	time.LoadLocation("Asia/Tokyo")
	date := time.Now()
	//データをインサート
	sqlStatement, err := con.Prepare("INSERT INTO messagerecord (user_id,rcvmessage,rcvdate) VALUES(?,?,?)")
	errorCheck(err)
	defer sqlStatement.Close()
	result, err := sqlStatement.Exec(user_id, rcvmessage, date)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
}
func Select_kaiin_host(user_id string) []Kaiin_host {
	var row *sql.Rows
	var err error

	con := Connectdb()
	defer con.Close()

	query := "SELECT * FROM kaiin_host WHERE user_id = ?"
	sqlStatement, err := con.Prepare(query)
	errorCheck(err)
	defer sqlStatement.Close()
	row, err = sqlStatement.Query(user_id)
	errorCheck(err)

	result := &Kaiin_host{}
	results := make([]Kaiin_host, 0)
	for row.Next() {
		err := row.Scan(&result.No, &result.User_id, &result.User_name)
		results = append(results, Kaiin_host{result.No, result.User_id, result.User_name})
		errorCheck(err)
	}
	return results
}

func Select_cashflow(user_id string) []Cashflow {
	var row *sql.Rows
	var err error

	con := Connectdb()
	defer con.Close()

	query := "SELECT * FROM (SELECT * FROM cashflow ORDER BY register_date DESC ) AS A WHERE user_id = ? LIMIT 1"
	sqlStatement, err := con.Prepare(query)
	errorCheck(err)
	defer sqlStatement.Close()
	row, err = sqlStatement.Query(user_id)
	errorCheck(err)

	result := &Cashflow{}
	results := make([]Cashflow, 0)
	for row.Next() {
		err := row.Scan(&result.No, &result.User_id, &result.Label, &result.Money, &result.Pauchased_item, &result.Register_date)
		results = append(results, Cashflow{result.No, result.User_id, result.Label, result.Money, result.Pauchased_item, result.Register_date})
		errorCheck(err)
	}
	return results
}

func Select_bodymanagement(user_id string) []Bodymanagement {
	var row *sql.Rows
	var err error

	con := Connectdb()
	defer con.Close()

	query := "SELECT * FROM (SELECT * FROM bodymanagement ORDER BY now_date DESC ) AS A WHERE user_id = ? LIMIT 1"
	sqlStatement, err := con.Prepare(query)
	errorCheck(err)
	defer sqlStatement.Close()
	row, err = sqlStatement.Query(user_id)
	errorCheck(err)

	result := &Bodymanagement{}
	results := make([]Bodymanagement, 0)
	for row.Next() {
		err := row.Scan(&result.No, &result.User_id, &result.Weight, &result.Height, &result.Now_date)
		results = append(results, Bodymanagement{result.No, result.User_id, result.Weight, result.Height, result.Now_date})
		errorCheck(err)
	}
	return results
}

func Select_messagerecord(user_id string) []Messagerecord {
	var row *sql.Rows
	var err error

	con := Connectdb()
	defer con.Close()

	query := "SELECT * FROM (SELECT * FROM messagerecord ORDER BY rcvdate DESC ) AS A WHERE user_id = ? LIMIT 1"
	sqlStatement, err := con.Prepare(query)
	errorCheck(err)
	defer sqlStatement.Close()
	row, err = sqlStatement.Query(user_id)
	errorCheck(err)

	result := &Messagerecord{}
	results := make([]Messagerecord, 0)
	for row.Next() {
		err := row.Scan(&result.No, &result.User_id, &result.Rcvmessage, &result.Rcvdate)
		results = append(results, Messagerecord{result.No, result.User_id, result.Rcvmessage, result.Rcvdate})
		errorCheck(err)
	}
	return results
}

func Update_data(user_name string, user_id *string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	//アップデートはインサートと同じ
	sqlStatement, err := con.Prepare("UPDATE kaiin_host SET user_name = ? WHERE user_id = ?")
	errorCheck(err)
	result, err := sqlStatement.Exec(user_name, user_id)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("updated %d row(s) of data.\n", rowCount)
}

func Delete_data(user_id string) {
	var rowCount int64
	con := Connectdb()
	defer con.Close()
	//デリートはインサートと同じ
	sqlStatement, err := con.Prepare("DELETE FROM kaiin_host WHERE user_id = ?")
	errorCheck(err)
	result, err := sqlStatement.Exec(user_id)
	errorCheck(err)
	rowCount, err = result.RowsAffected()
	errorCheck(err)
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
}

//デバッグ用
func Select_data(table_name string) {
	var row *sql.Rows
	var err error

	con := Connectdb()
	defer con.Close()

	query := "SELECT count(*) FROM " + table_name
	sqlStatement, err := con.Prepare(query)
	errorCheck(err)
	defer sqlStatement.Close()
	row, err = sqlStatement.Query()
	errorCheck(err)

	var rowcount string
	for row.Next() {
		err := row.Scan(&rowcount)
		errorCheck(err)
	}
}

func Truncate_data() {
	var err error
	con := Connectdb()
	defer con.Close()
	//外部キー制約のため一時的ににオフにする
	con.Exec("set foreign_key_checks = 0")
	//テーブルを空に
	_, err = con.Exec("TRUNCATE cashflow")
	errorCheck(err)
	_, err = con.Exec("TRUNCATE bodymanagement")
	errorCheck(err)
	_, err = con.Exec("TRUNCATE messagerecord")
	errorCheck(err)
	_, err = con.Exec("TRUNCATE kaiin_host")
	errorCheck(err)
	_, err = con.Exec("set foreign_key_checks = 1")
	errorCheck(err)
}
