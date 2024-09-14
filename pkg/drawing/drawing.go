package drawing

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

func test() {
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

type Data[T comparable] struct {
	results              map[T][]T
	freshmen             []T
	seniors              []T
	pairableSeniors      []T
	seniorsPairedCount   map[T]int
	waitingFreshmenCount int
	luckyCount           int
	seniorsPairedMax     int
	baseDrawTimes        int
}

func newData[T comparable](freshmen []T, seniors []T) Data[T] {
	data := Data[T]{
		results:              make(map[T][]T, len(freshmen)),
		freshmen:             slices.Clone(freshmen),
		seniors:              slices.Clone(seniors),
		pairableSeniors:      slices.Clone(seniors),
		seniorsPairedCount:   make(map[T]int, len(seniors)),
		waitingFreshmenCount: len(freshmen),
		luckyCount:           0,
		seniorsPairedMax:     len(freshmen) / len(seniors),
		baseDrawTimes:        1,
	}

	for _, value := range seniors {
		data.seniorsPairedCount[value] = 0
	}

	if len(freshmen) < len(seniors) {
		data.baseDrawTimes = len(seniors) / len(freshmen)
	}

	if len(seniors) > len(freshmen) {
		data.luckyCount = len(seniors) % len(freshmen)
	}

	if len(freshmen)%len(seniors) != 0 {
		data.seniorsPairedMax++
	}

	return data
}

func (d *Data[T]) Results() map[T][]T {
	return maps.Clone(d.results)
}

func (d *Data[T]) Draw(freshman T) ([]T, error) {
	if !slices.Contains(d.freshmen, freshman) {
		return []T{}, fmt.Errorf("freshNumber '%v' not in fresh list", freshman)
	}

	if result, inMap := d.results[freshman]; inMap {
		return result, nil
	}

	drawTimes := d.baseDrawTimes
	if rand.Intn(int(d.waitingFreshmenCount)) < int(d.luckyCount) {
		drawTimes++
		d.luckyCount--
	}

	result := make([]T, 0, drawTimes)
	for range drawTimes {
		randIndex := rand.Intn(len(d.pairableSeniors))
		for slices.Contains(result, d.pairableSeniors[randIndex]) {
			randIndex = rand.Intn(len(d.pairableSeniors))
		}
		paired := d.pairableSeniors[randIndex]

		result = append(result, paired)
		d.seniorsPairedCount[paired]++

		if d.seniorsPairedCount[paired] >= d.seniorsPairedMax {
			d.pairableSeniors[randIndex] = d.pairableSeniors[len(d.pairableSeniors)-1]
			d.pairableSeniors = d.pairableSeniors[:len(d.pairableSeniors)-1]
		}
	}

	d.waitingFreshmenCount--

	return result, nil
}
