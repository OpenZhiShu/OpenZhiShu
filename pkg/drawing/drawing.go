package drawing

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

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

func MakeData[T comparable](freshmen []T, seniors []T) Data[T] {
	return Data[T]{
		results:         make(map[T][]T, len(freshmen)),
		freshmen:        slices.Clone(freshmen),
		seniors:         slices.Clone(seniors),
		pairableSeniors: slices.Clone(seniors),
		seniorsPairedCount: maps.Collect(func(yield func(T, int) bool) {
			for _, k := range seniors {
				if !yield(k, 0) {
					return
				}
			}
		}),
		waitingFreshmenCount: len(freshmen),
		luckyCount: func() int {
			if len(seniors) > len(freshmen) {
				return len(seniors) % len(freshmen)
			}
			return 0
		}(),
		seniorsPairedMax: len(freshmen)/len(seniors) + min(1, len(freshmen)%len(seniors)),
		baseDrawTimes:    max(1, len(seniors)/len(freshmen)),
	}
}

func (d *Data[T]) Results() map[T][]T {
	return maps.Clone(d.results)
}

func (d *Data[T]) ResultsBySenior() map[T][]T {
	results := make(map[T][]T, len(d.seniors))
	for k, vs := range d.results {
		for _, v := range vs {
			results[v] = append(results[v], k)
		}
	}
	return results
}

func (d *Data[T]) WaitingFreshmenCount() int {
	return d.waitingFreshmenCount
}

func (d *Data[T]) BaseDrawTimes() int {
	return d.baseDrawTimes
}

func (d *Data[T]) SeniorsPairedMax() int {
	return d.seniorsPairedMax
}

func (d *Data[T]) LuckyCount() int {
	return d.luckyCount
}

func (d *Data[T]) Finished() bool {
	return d.waitingFreshmenCount == 0
}

func (d *Data[T]) DrawAll() (map[T][]T, error) {
	for _, freshman := range d.freshmen {
		_, err := d.Draw(freshman)
		if err != nil {
			return map[T][]T{}, nil
		}
	}
	return d.Results(), nil
}

func (d *Data[T]) Draw(freshman T) ([]T, error) {
	if !slices.Contains(d.freshmen, freshman) {
		return []T{}, fmt.Errorf("freshNumber '%v' not in fresh list", freshman)
	}

	if result, inMap := d.results[freshman]; inMap {
		return result, nil
	}

	drawTimes := d.baseDrawTimes
	if rand.Intn(d.waitingFreshmenCount) < int(d.luckyCount) {
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

	d.results[freshman] = result
	d.waitingFreshmenCount--

	return result, nil
}
