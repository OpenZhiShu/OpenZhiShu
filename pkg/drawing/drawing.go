package drawing

import (
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

type Person[T comparable] interface {
	Key() T
}

type Results[P Person[T], T comparable] struct {
	results map[T][]P
}

func (r Results[P, T]) Index(person P) []P {
	return r.results[person.Key()]
}

func (r Results[P, T]) Contains(person P) bool {
	_, contains := r.results[person.Key()]
	return contains
}

func (r *Results[P, T]) Len() int {
	return len(r.results)
}

type Data[P Person[T], T comparable] struct {
	results              Results[P, T]
	freshmen             []P
	seniors              []P
	pairableSeniors      []P
	seniorsPairedCount   map[T]int
	waitingFreshmenCount int
	luckyCount           int
	seniorsPairedMax     int
	baseDrawTimes        int
}

func MakeData[P Person[T], T comparable](freshmen []P, seniors []P) Data[P, T] {
	return Data[P, T]{
		results:              Results[P, T]{make(map[T][]P, len(freshmen))},
		freshmen:             slices.Clone(freshmen),
		seniors:              slices.Clone(seniors),
		pairableSeniors:      slices.Clone(seniors),
		seniorsPairedCount:   make(map[T]int, len(seniors)),
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

func (d *Data[P, T]) Results() Results[P, T] {
	return Results[P, T]{maps.Clone(d.results.results)}
}

func (d *Data[P, T]) ResultsBySenior() Results[P, T] {
	results := Results[P, T]{make(map[T][]P, len(d.seniors))}
	for k, vs := range d.results.results {
		var freshman P
		for _, f := range d.freshmen {
			if f.Key() == k {
				freshman = f
				break
			}
		}
		for _, v := range vs {
			results.results[v.Key()] = append(results.results[v.Key()], freshman)
		}
	}
	return results
}

func (d *Data[P, T]) WaitingFreshmenCount() int {
	return d.waitingFreshmenCount
}

func (d *Data[P, T]) BaseDrawTimes() int {
	return d.baseDrawTimes
}

func (d *Data[P, T]) SeniorsPairedMax() int {
	return d.seniorsPairedMax
}

func (d *Data[P, T]) LuckyCount() int {
	return d.luckyCount
}

func (d *Data[P, T]) Finished() bool {
	return d.waitingFreshmenCount == 0
}

func (d *Data[P, T]) DrawAll() (Results[P, T], error) {
	for _, freshman := range d.freshmen {
		_, err := d.Draw(freshman)
		if err != nil {
			return Results[P, T]{}, nil
		}
	}
	return d.Results(), nil
}

func (d *Data[P, T]) Draw(freshman P) ([]P, error) {
	if !slices.ContainsFunc(d.freshmen, func(p P) bool { return p.Key() == freshman.Key() }) {
		return []P{}, fmt.Errorf("freahman '%v' not in fresh list", freshman)
	}

	if d.results.Contains(freshman) {
		return d.results.Index(freshman), nil
	}

	drawTimes := d.baseDrawTimes
	if rand.Intn(d.waitingFreshmenCount) < int(d.luckyCount) {
		drawTimes++
		d.luckyCount--
	}

	result := make([]P, 0, drawTimes)
	for range drawTimes {
		randIndex := rand.Intn(len(d.pairableSeniors))
		for slices.ContainsFunc(result, func(p P) bool { return p.Key() == d.pairableSeniors[randIndex].Key() }) {
			randIndex = rand.Intn(len(d.pairableSeniors))
		}
		paired := d.pairableSeniors[randIndex]

		result = append(result, paired)
		d.seniorsPairedCount[paired.Key()]++

		if d.seniorsPairedCount[paired.Key()] >= d.seniorsPairedMax {
			d.pairableSeniors[randIndex] = d.pairableSeniors[len(d.pairableSeniors)-1]
			d.pairableSeniors = d.pairableSeniors[:len(d.pairableSeniors)-1]
		}
	}

	d.results.results[freshman.Key()] = result
	d.waitingFreshmenCount--

	return result, nil
}
