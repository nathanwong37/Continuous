package temp

import (
	"database/sql"
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
	var query string = remov + " FROM timer " + where + " " + timerId + "= UUID_TO_BIN( '" + uuid.String() + "') AND " + nameSpace + " = \"" + namespace + "\" ;"
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
	var query string = "UPDATE timer Se t" + mostRecent + " = \"" + recent + "\", " + amountFired + " = '" + strconv.Itoa(fired) +
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
	_, err := time.Parse("2006-01-02 15:04:05", timerStart)
	if err != nil {
		valid = false
	}
	if !valid {
		now := time.Now()
		timerStart = now.Format("2006-01-02 15:04:05")
	}
	begin = begin + startTime + "," + amountFired
	end = end + timerStart + "' , " + strconv.Itoa(int(timerInfo.GetAmountFired()))
	result := begin + end + ");"
	return result
}
