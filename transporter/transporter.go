package transporter

import (
	"database/sql"
	"fmt"
	"strconv"

	//need for driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var datab = "mysql"
var connect = "test:test1@tcp(127.0.0.1:3306)/timers"
var insert = "INSERT"
var delete = "DELETE"
var where = "WHERE"
var timerId = "timer_id"
var shardId = "shard_id"
var nameSpace = "namespace"
var interval = "interval_"
var count = "count"
var startTime = "start_time"
var mostRecent = "most_recent"
var amountFired = "amount_fired"
var values = "VALUES"
var into = "INTO"

//Transport struct just to call methods
type Transport struct{}

//TimerInfo holds all the information that a timer has
type TimerInfo struct {
	TimerID     uuid.UUID `json:"timerID"`
	ShardID     int       `json:"shardID"`
	Namespace   string    `json:"namespace"`
	Interval    string    `json:"interval"`
	Count       int       `json:"count"`
	Starttime   string    `json:"startTime"`
	Mostrecent  string    `json:"mostRecent"`
	Amountfired int       `json:"amountFired"`
	Timecreated string    `json:"timeCreated"`
}

func (transporter *Transport) connect() {
	fmt.Println("Go SQL")

	db, err := sql.Open(datab, connect)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
}

//Get connects to database, and gets the matching values
//returns a struct of information
//uuid and namespace should be validated before hand
func (transporter *Transport) Get(uuid uuid.UUID, namespace string) (*TimerInfo, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var query string = "SELECT * FROM timer WHERE timer_id = UUID_TO_BIN(\"" + uuid.String() + "\") AND " + nameSpace + " = \"" + namespace + "\" ;"
	results, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var timerInfo TimerInfo
	for results.Next() {
		err = results.Scan(&timerInfo.TimerID, &timerInfo.ShardID, &timerInfo.Namespace, &timerInfo.Interval, &timerInfo.Count, &timerInfo.Starttime, &timerInfo.Mostrecent,
			&timerInfo.Amountfired, &timerInfo.Timecreated)
		if err != nil {
			return nil, err
		}
	}
	return &timerInfo, nil
}

//Create connects to database and puts in the values
//should take in some TimerInfo object,
func (transporter *Transport) Create(timerinfo *TimerInfo) (bool, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return false, err
	}
	defer db.Close()
	var query string = insert + " " + into + " timer (" + timerId + "," + shardId + "," + nameSpace + "," + interval +
		"," + count + "," + startTime + "," + mostRecent + "," + amountFired + ") " + values +
		"(UUID_TO_BIN( '" + timerinfo.TimerID.String() + "')," + strconv.Itoa(timerinfo.ShardID) + ", \"" + timerinfo.Namespace +
		"\", '" + timerinfo.Interval + "' ," + strconv.Itoa(timerinfo.Count) + ", '" + timerinfo.Starttime + "' , '" + timerinfo.Mostrecent +
		"' , " + strconv.Itoa(timerinfo.Amountfired) + ");"
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil
}

//Delete from the database, if uuid and namespace match
func (transporter *Transport) Delete(uuid uuid.UUID, namespace string) (bool, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return false, err
	}
	defer db.Close()
	var query string = delete + " FROM timer " + where + " " + timerId + "= UUID_TO_BIN( '" + uuid.String() + "') AND " + nameSpace + " = \"" + namespace + "\" ;"
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil
}
