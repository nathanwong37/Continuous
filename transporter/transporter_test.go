package Continuous

import (
	"testing"

	proto "github.com/Continuous/plugins"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
		ShardID:     344,
		NameSpace:   "Nathan Wong",
		Interval:    "00:00:10",
		Count:       1,
		StartTime:   "2020-08-03 18:18:50",
		AmountFired: 0,
	}
	work, err := transporter.Create(testTimer)
	assert.NoError(t, err)
	assert.Equal(t, work, true)
	transporter.GetRows(344)

}
