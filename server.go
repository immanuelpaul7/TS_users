package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func dbConn() (db *sql.DB) {
	dbDriver := "dbDriver"
	dbUser := "dbUser"
	dbPass := "dbPass"
	dbName := "dbName"
  dbHost := "dbHost"
  dbPort := "dbPort"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	checkErr(err)
	return db
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	rows, err := db.Query("SELECT * FROM TS_users")
	defer rows.Close()
	checkErr(err)

	records := getMysqlData(rows)
	defer db.Close()

	var users []UserInfo
	users = FormatRecordsFromDB(records)
	json.NewEncoder(w).Encode(users)
}

// Format Users Data
func FormatRecordsFromDB(records []interface{}) []UserInfo {
	var users []UserInfo
	if records == nil {
		return users
	}

	for _, cols := range records {
		var user UserInfo
		col, _ := cols.(map[string]interface{})

		user.ID = col["id"].(string)
		user.ContactMe = col["contact_me"].(string)
		user.Phone = col["phone"].(string)
		user.FirstName = col["first_name"].(string)
		user.LastName = col["last_name"].(string)
		user.HowLong = col["how_long"].(string)
		user.Email = col["email"].(string)

		users = append(users, user)
	}

	return users
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/user", insertUserInfo).Methods("POST")
	router.HandleFunc("/winner", getAWinner).Methods("GET")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type structInterface interface{}

type mySQLData struct {
	rows *sql.Rows
	structInterface
}

func getMysqlData(rows *sql.Rows) []interface{} {
	columns, err := rows.Columns()
	checkErr(err)
	myMap := []interface{}{}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		checkErr(err)

		var value string
		fields := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fields[columns[i]] = value
		}
		myMap = append(myMap, fields)
	}
	return myMap
}

// UserInfo struct for user
type UserInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	Phone     string `json:"Phone"`
	ContactMe string `json:"ContactMe"`
	HowLong   string `json:"HowLong"`
}

func insertUserInfo(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 400)
	}

	var Person UserInfo
	_ = json.NewDecoder(r.Body).Decode(&Person)

	insertUserRecord, err := db.Prepare("INSERT INTO TS_users (first_name, last_name, email, phone, contact_me, how_long) VALUES(?,?,?,?,?,?)")
	checkErr(err)
	res, err := insertUserRecord.Exec(
		Person.FirstName,
		Person.LastName,
		Person.Email,
		Person.Phone,
		Person.ContactMe,
		Person.HowLong,
	)

	checkErr(err)
	id, _ := res.LastInsertId()
	defer db.Close()
	Person.ID = strconv.FormatInt(id, 16)
	json.NewEncoder(w).Encode(Person)

}

func getAWinner(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	rows, err := db.Query("SELECT * from TS_users WHERE is_winner IS NULL and created_time BETWEEN DATE_SUB(NOW(), INTERVAL 30 MINUTE) AND NOW() ORDER BY RAND() LIMIT 1")
	defer rows.Close()
	checkErr(err)

	records := getMysqlData(rows)
	defer db.Close()

	var users []UserInfo
	users = FormatRecordsFromDB(records)

	updateUserAsWinner(users)
	json.NewEncoder(w).Encode(users)
}

func updateUserAsWinner(user []UserInfo) {
	db := dbConn()
	stmt, err := db.Prepare("UPDATE TS_users set is_winner=1 WHERE id =?")
	checkErr(err)

	row := user
	id := (row)[0].ID
	_, err = stmt.Exec(id)
	checkErr(err)
	defer db.Close()
}
