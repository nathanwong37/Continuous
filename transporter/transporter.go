package transporter

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	proto "github.com/Continuous/plugins"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	//Need to run sqldriver
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var datab = "mysql"
var connect = "test1:test@tcp(127.0.0.1:3306)/timer"
var databName = "timer"

//Transporter interface is an interface to what transporter should have
type Transporter interface {
	Get(uuid.UUID, string) (*RetTimerInfo, error)
	Create(proto.TimerInfo) (bool, error)
	Remove(uuid.UUID, string) (bool, error)
	Update(string, string, string, int) (bool, error)
	GetRows(int) ([]*proto.TimerInfo, error)
	BuildQuery(string, string, *proto.TimerInfo) (string, error)
}

//Transport struct just to call methods
type Transport struct{}

//RetTimerInfo holds all the information that a timer has
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

//Get connects to database, and gets the matching values
//returns a struct of information
//uuid and namespace should be validated before hand
func (transporter *Transport) Get(uuid uuid.UUID, namespace string) (*RetTimerInfo, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	uuidstr := uuid.String()
	parameter := &proto.TimerInfo{
		TimerID:   uuidstr,
		NameSpace: namespace,
	}
	query, err := transporter.BuildQuery(databName, "get", parameter)
	if err != nil {
		return nil, err
	}
	results := db.QueryRow(query)

	var timerInfo RetTimerInfo
	err = results.Scan(&timerInfo.TimerID, &timerInfo.ShardID, &timerInfo.Namespace, &timerInfo.Interval, &timerInfo.Count, &timerInfo.Starttime, &timerInfo.Mostrecent,
		&timerInfo.Amountfired, &timerInfo.Timecreated)
	if err != nil {
		return nil, err
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
	timerInfo.MostRecent = time.Now().Format("2006-01-02 15:04:05")
	timerStart := timerInfo.GetStartTime()
	valid := true
	_, err = time.ParseInLocation("2006-01-02 15:04:05", timerStart, time.Local)
	if err != nil {
		valid = false
	}
	if !valid {
		timerInfo.StartTime = time.Now().Format("2006-01-02 15:04:05")
	}
	query, err := transporter.BuildQuery(databName, "create", timerInfo)
	if err != nil {
		return false, err
	}
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
	uuidstr := uuid.String()
	parameter := &proto.TimerInfo{
		TimerID:   uuidstr,
		NameSpace: namespace,
	}
	query, err := transporter.BuildQuery(databName, "remove", parameter)
	if err != nil {
		return false, err
	}
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
	parameter := &proto.TimerInfo{
		TimerID:     uuid,
		NameSpace:   namespace,
		MostRecent:  recent,
		AmountFired: int32(fired),
	}
	query, err := transporter.BuildQuery(databName, "update", parameter)
	result, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return true, nil

}

//GetRows returns the rows of data that is read in from the database
func (transporter *Transport) GetRows(shardID int) ([]*proto.TimerInfo, error) {
	db, err := sql.Open(datab, connect)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	parameter := &proto.TimerInfo{
		ShardID: int32(shardID),
	}
	query, err := transporter.BuildQuery(databName, "getRow", parameter)
	results, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
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

//BuildQuery is used to build queries
func (transporter *Transport) BuildQuery(database, command string, parameter *proto.TimerInfo) (string, error) {
	sql := ""
	dialect := goqu.Dialect("mysql")
	switch command {
	case "get":
		sql, _, _ = dialect.From(database).Where(goqu.Ex{
			"timer_id":  "",
			"namespace": parameter.GetNameSpace(),
		}).ToSQL()
		sql = addUUIDGetRem(sql, parameter.GetTimerID())
	case "getRow":
		sql, _, _ = dialect.From(database).Where(goqu.Ex{
			"shard_id": int(parameter.GetShardID()),
		}).ToSQL()
	case "create":
		create := dialect.Insert(database).Rows(
			goqu.Record{"timer_id": "", "shard_id": int(parameter.GetShardID()), "namespace": parameter.GetNameSpace(), "interval_": parameter.GetInterval(), "count": int(parameter.GetCount()), "start_time": parameter.GetStartTime(),
				"most_recent": parameter.GetMostRecent(), "amount_fired": parameter.GetAmountFired()},
		)
		sql, _, _ = create.ToSQL()
		sql = addUUID(sql, parameter.GetTimerID())

	case "remove":
		ds := dialect.Delete(database).Where(goqu.C("namespace").Eq(parameter.GetNameSpace()), goqu.C("timer_id").Eq(""))
		sql, _, _ = ds.ToSQL()
		sql = addUUIDGetRem(sql, parameter.GetTimerID())
	case "update":
		ds := dialect.Update(database).Where(goqu.C("namespace").Eq(parameter.GetNameSpace()), goqu.C("timer_id").Eq("")).Set(
			goqu.Record{"most_recent": parameter.GetMostRecent(), "amount_fired": parameter.GetAmountFired()},
		)
		sql, _, _ = ds.ToSQL()
		sql = addUUIDGetRem(sql, parameter.GetTimerID())
	default:
		return sql, errors.New("Improper Command")
	}
	return sql, nil
}

//should be safe, since we generate uuid
func addUUID(sql string, uuidstr string) string {
	sql = (string)([]rune(sql)[:len(sql)-3])
	sql = sql + "UUID_TO_BIN( '" + uuidstr + "'))"
	return sql
}

func addUUIDGetRem(sql string, uuidstr string) string {
	sql = (string)([]rune(sql)[:len(sql)-4])
	sql = sql + "UUID_TO_BIN( '" + uuidstr + "')))"
	return sql
}
