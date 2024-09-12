package main

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

func main() {
	freshmen := []int{1, 2, 3, 4, 5, 6, 7}
	seniors := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// seniors := []int{1, 2, 3, 4, 5, 6}
	d := newData(freshmen, seniors)
	fmt.Printf("waitingFreshmenCount: %v\n", d.waitingFreshmenCount)
	fmt.Printf("luckyCount: %v\n", d.luckyCount)
	fmt.Printf("seniorsPairedMax: %v\n", d.seniorsPairedMax)
	fmt.Printf("baseDrawTimes: %v\n", d.baseDrawTimes)
	fmt.Println("")
	for _, v := range freshmen {
		result, err := d.Draw(v)
		fmt.Printf("%v: %v (%v)\n", v, result, err)
	}
}

type Data struct {
	results              map[int][]int
	freshmen             []int
	seniors              []int
	pairableSeniors      []int
	seniorsPairedCount   map[int]uint
	waitingFreshmenCount uint
	luckyCount           uint
	seniorsPairedMax     uint
	baseDrawTimes        uint
}

func newData(freshmen []int, seniors []int) Data {
	data := Data{
		results:              make(map[int][]int, len(freshmen)),
		freshmen:             slices.Clone(freshmen),
		seniors:              slices.Clone(seniors),
		pairableSeniors:      slices.Clone(seniors),
		seniorsPairedCount:   make(map[int]uint, len(seniors)),
		waitingFreshmenCount: uint(len(freshmen)),
		luckyCount:           0,
		seniorsPairedMax:     uint(len(freshmen) / len(seniors)),
		baseDrawTimes:        1,
	}

	for k := range seniors {
		data.seniorsPairedCount[k] = 0
	}

	if len(freshmen) < len(seniors) {
		data.baseDrawTimes = uint(len(seniors) / len(freshmen))
	}

	if len(seniors) > len(freshmen) {
		data.luckyCount = uint(len(seniors) % len(freshmen))
	}

	if len(freshmen)%len(seniors) != 0 {
		data.seniorsPairedMax++
	}

	return data
}

func (d *Data) Results() map[int][]int {
	return maps.Clone(d.results)
}

func (d *Data) Draw(freshNumber int) ([]int, error) {
	if !slices.Contains(d.freshmen, freshNumber) {
		return []int{}, fmt.Errorf("freshNumber '%v' not in fresh list", freshNumber)
	}

	if result, inMap := d.results[freshNumber]; inMap {
		return result, nil
	}

	drawTimes := d.baseDrawTimes
	if rand.Intn(int(d.waitingFreshmenCount)) < int(d.luckyCount) {
		drawTimes++
		d.luckyCount--
	}

	result := make([]int, 0, drawTimes)
	for range drawTimes {
		randIndex := rand.Intn(len(d.pairableSeniors))
		for slices.Contains(result, d.pairableSeniors[randIndex]) {
			randIndex = rand.Intn(len(d.pairableSeniors))
		}

		d.seniorsPairedCount[randIndex]++
		result = append(result, d.pairableSeniors[randIndex])

		if d.seniorsPairedCount[randIndex] >= d.seniorsPairedMax {
			d.pairableSeniors[randIndex] = d.pairableSeniors[len(d.pairableSeniors)-1]
			d.pairableSeniors = d.pairableSeniors[:len(d.pairableSeniors)-1]
		}
	}

	d.waitingFreshmenCount--

	return result, nil
}
