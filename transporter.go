package temp

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	//need for driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	proto "github.com/temp/plugins"
)

var datab = "mysql"
var connect = "test:test1@tcp(127.0.0.1:3306)/timers"
var insert = "INSERT"
var remov = "DELETE"
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
type RetTimerInfo struct {
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

	db, err := sql.Open(datab, connect)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
}

//Get connects to database, and gets the matching values
//returns a struct of information
//uuid and namespace should be validated before hand
func (transporter *Transport) Get(uuid uuid.UUID, namespace string) (*RetTimerInfo, error) {
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

	var timerInfo RetTimerInfo
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
func (transporter *Transport) Create(timerInfo *proto.TimerInfo) (bool, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return false, err
	}
	defer db.Close()
	query := createQueryBuilder(timerInfo)
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil
}

//Remove from the database, if uuid and namespace match
func (transporter *Transport) Remove(uuid uuid.UUID, namespace string) (bool, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return false, err
	}
	defer db.Close()
	var query string = "DELETE FROM timer WHERE timer_id = UUID_TO_BIN(\"" + uuid.String() + "\") AND " + nameSpace + " = \"" + namespace + "\""
	fmt.Println(query)
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil
}

//Update should only be used to update RecentTime and AmountFired
func (transporter *Transport) Update(uuid, recent, namespace string, fired int) (bool, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return false, err
	}
	defer db.Close()
	var query string = "UPDATE timer Set " + mostRecent + " = \"" + recent + "\", " + amountFired + " = '" + strconv.Itoa(fired) +
		"' WHERE timer_id = UUID_TO_BIN(\"" + uuid + "\") AND " + nameSpace + " = \"" + namespace + "\" ;"
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil

}

//createQueryString helps create the query
func createQueryBuilder(timerInfo *proto.TimerInfo) string {
	var begin, end string
	//know for sure every timer has a uuid, shardId,namesapce, and count
	begin = insert + " " + into + " timer (" + timerId + "," + shardId + "," + nameSpace + "," + interval + "," + count + ","
	end = ") " + values + "(UUID_TO_BIN( '" + timerInfo.GetTimerID() + "')," + strconv.Itoa(int(timerInfo.GetShardID())) + ", \"" + timerInfo.GetNameSpace() +
		"\", '" + timerInfo.GetInterval() + "' ," + strconv.Itoa(int(timerInfo.GetCount())) + ", '"
	var valid bool = true
	timerStart := timerInfo.GetStartTime()
	now := time.Now()
	_, err := time.ParseInLocation("2006-01-02 15:04:05", timerStart, time.Local)
	if err != nil {
		valid = false
	}
	if !valid {
		timerStart = now.Format("2006-01-02 15:04:05")
	}
	begin = begin + startTime + "," + mostRecent + "," + amountFired
	end = end + timerStart + "' , '" + now.Format("2006-01-02 15:04:05") + " ' ," + strconv.Itoa(int(timerInfo.GetAmountFired()))
	result := begin + end + ");"
	return result
}

//GetRows returns the rows of data that is read in from the database
func (transporter *Transport) GetRows(shardId int) ([]*proto.TimerInfo, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var query string = "Select * FROM timer WHERE shard_id = " + strconv.Itoa(shardId) + ";"
	results, err := db.Query(query)
	timers := make([]*proto.TimerInfo, 0)
	for results.Next() {
		info := new(proto.TimerInfo)
		a := uuid.New()
		err = results.Scan(&a, &info.ShardID, &info.NameSpace, &info.Interval, &info.Count, &info.StartTime, &info.MostRecent,
			&info.AmountFired, &info.TimeCreated)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		info.TimerID = a.String()
		timers = append(timers, info)
	}
	return timers, nil
}
