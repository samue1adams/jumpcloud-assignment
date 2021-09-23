package main

import (
	"sync"
	"testing"
)

func TestGetAverage(t *testing.T) {
	intArry := []int64{4, 3, 4, 5}
	average := getAverage(intArry)
	if average != 4 {
		t.Errorf("average was incorrect")
	}
}

func TestHashPassword(t *testing.T) {
	m := make(map[int]string)
	wg := sync.WaitGroup{}
	wg.Add(1)
	p := Password{
		Password: "angryMonkey",
		id:       5}

	hashPassword(p, &wg, m)

	if m[5] == "angryMonkey" {
		t.Errorf("password was not hashed")
	}

}
