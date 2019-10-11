package test

import (
	"client"
	"testing"
)

func TestBasic(t *testing.T) {
	client.PerfomanceTest("conf/client.json")
}
