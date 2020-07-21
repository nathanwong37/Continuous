package api

import "testing"

func TestRun(t *testing.T) {
	listen := new(listener)
	listen.run()
}
