package temp

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	proto "github.com/temp/plugins"
)

func TestConnect(t *testing.T) {
	transporter := new(Transport)
	transporter.connect()
}

func TestGet(t *testing.T) {
	transporter := new(Transport)

	var personalUUID string = "59d3d5c6-bcc0-11ea-a5d3-93973a5f2fcc"
	uuid, err := uuid.Parse(personalUUID)
	testTimer := &RetTimerInfo{
		TimerID:     uuid,
		ShardID:     1,
		Namespace:   "Nathan Wong",
		Interval:    "00:00:10",
		Count:       1,
		Starttime:   "2020-12-24 11:59:50",
		Mostrecent:  "2020-12-24 11:59:50",
		Amountfired: 0,
	}
	assert.NoError(t, err)
	a, err := transporter.Get(uuid, "Nathan Wong")
	assert.NoError(t, err)

	assert.Equal(t, testTimer.TimerID, a.TimerID)
	assert.Equal(t, testTimer.ShardID, a.ShardID)
	assert.Equal(t, testTimer.Namespace, a.Namespace)
	assert.Equal(t, testTimer.Interval, a.Interval)
	assert.Equal(t, testTimer.Count, a.Count)
	assert.Equal(t, testTimer.Starttime, a.Starttime)
	assert.Equal(t, testTimer.Mostrecent, a.Mostrecent)
	assert.Equal(t, testTimer.Amountfired, a.Amountfired)
}

func TestCreate(t *testing.T) {
	transporter := new(Transport)
	uuid := uuid.New()
	testTimer := &proto.TimerInfo{
		TimerID:     uuid.String(),
		ShardID:     1,
		NameSpace:   "Nathan Wong",
		Interval:    "00:00:10",
		Count:       1,
		StartTime:   "2020-12-24 14:59:50",
		MostRecent:  "2020-12-24 14:59:50",
		AmountFired: 0,
	}
	work, err := transporter.Create(testTimer)
	assert.NoError(t, err)
	assert.Equal(t, work, true)

	// validate, err := transporter.Get(uuid.testTimer.GetTimerId(, testTimer.Namespace)
	// assert.NoError(t, err)

	// assert.Equal(t, testTimer.TimerID, validate.TimerID)
	// assert.Equal(t, testTimer.ShardID, validate.ShardID)
	// assert.Equal(t, testTimer.Namespace, validate.Namespace)
	// assert.Equal(t, testTimer.Interval, validate.Interval)
	// assert.Equal(t, testTimer.Count, validate.Count)
	// assert.Equal(t, testTimer.Starttime, validate.Starttime)
	// assert.Equal(t, testTimer.Mostrecent, validate.Mostrecent)
	// assert.Equal(t, testTimer.Amountfired, validate.Amountfired)

	// valid, err := transporter.Remove(testTimer.TimerID, testTimer.Namespace)
	// assert.NoError(t, err)
	// assert.Equal(t, valid, true)

}
