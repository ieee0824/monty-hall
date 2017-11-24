package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Doors struct {
	Prize   int
	opened  int
	Closed  []int
	Choiced int
}

func NewDoors() *Doors {
	return &Doors{
		rand.Int() % 3,
		-1,
		nil,
		-1,
	}
}

func (d *Doors) choice() {
	d.Choiced = rand.Int() % 3
}

func (d *Doors) open() {
	for i := 0; i < 3; i++ {
		if d.Prize == i || d.Choiced == i {
			continue
		}
		d.opened = i
		return
	}
}

func (d *Doors) closed() []int {
	ret := []int{}
	for i := 0; i < 3; i++ {
		if d.opened == i {
			continue
		}
		ret = append(ret, i)
	}
	return ret
}

func (d *Doors) result() (move, notMove bool) {
	return d.Prize != d.Choiced, d.Prize == d.Choiced
}

func (d Doors) String() string {
	bin, _ := json.MarshalIndent(d, "", "    ")
	return string(bin)
}

func sum(in chan map[string]int64) {
	results := map[string]int64{
		"move":    0,
		"notMove": 0,
	}
	for {
		r, ok := <-in
		if !ok {
			break
		}
		results["move"] += r["move"]
		results["notMove"] += r["notMove"]
	}
	fmt.Println(results)
	if results["move"] > results["notMove"] {
		fmt.Printf("動いたほうが%f倍確率高い\n", float64(results["move"])/float64(results["notMove"]))
	} else {
		fmt.Printf("動かないほうが%f倍確率高い\n", float64(results["notMove"])/float64(results["move"]))
	}
}

func main() {
	var wg sync.WaitGroup

	chResults := make(chan map[string]int64, 16)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			results := map[string]int64{
				"move":    0,
				"notMove": 0,
			}
			defer wg.Done()
			for i := 0; i < 10000000; i++ {
				doors := NewDoors()
				doors.choice()
				doors.open()
				doors.Closed = append(doors.Closed, doors.closed()...)
				move, notMove := doors.result()
				if move {
					results["move"]++
				}

				if notMove {
					results["notMove"]++
				}
			}
			chResults <- results
		}()
	}
	wg.Wait()

	close(chResults)

	sum(chResults)
}
