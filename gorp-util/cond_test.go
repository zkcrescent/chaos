package gorpUtil

import (
	"log"
	"testing"
)

func TestCondition(t *testing.T) {
	k := "aa"
	c := PureCondition(k)
	str, _ := c.String(1)
	log.Println(str)
}
