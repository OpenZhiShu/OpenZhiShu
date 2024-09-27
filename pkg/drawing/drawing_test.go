package drawing

import (
	"slices"
	"testing"
)

type person int

func (p person) Key() int {
	return int(p)
}

func TestMakeData(t *testing.T) {
	{
		freshmen := []person{1, 2, 3, 4, 5, 6, 7}
		seniors := []person{1, 2, 3, 4, 5, 6, 7, 8, 9}
		d := MakeData(freshmen, seniors)

		luckyCount := 2
		baseDrawTimes := 1
		seniorsPairedMax := 1

		if len(d.results) != 0 {
			t.Error("results is not empty")
		}
		if !slices.Equal(d.freshmen, freshmen) {
			t.Error("d.freshmen is not equal to freshmen")
		}
		if !slices.Equal(d.seniors, seniors) {
			t.Error("d.seniors is not equal to seniors")
		}
		if !slices.Equal(d.pairableSeniors, seniors) {
			t.Error("d.pairableSeniors is not equal to freshmen")
		}
		if d.waitingFreshmenCount != len(freshmen) {
			t.Errorf("seniorsPairedMax is %v instead of %v", d.waitingFreshmenCount, len(freshmen))
		}
		if d.luckyCount != luckyCount {
			t.Errorf("luckyCount is %v instead of %v", d.luckyCount, luckyCount)
		}
		if d.baseDrawTimes != baseDrawTimes {
			t.Errorf("baseDrawTimes is %v instead of %v", d.baseDrawTimes, baseDrawTimes)
		}
		if d.seniorsPairedMax != seniorsPairedMax {
			t.Errorf("seniorsPairedMax is %v instead of %v", d.seniorsPairedMax, seniorsPairedMax)
		}
		for k, v := range d.seniorsPairedCount {
			if v != 0 {
				t.Errorf("seniorsPairedCount of %v is %v instead of 0", k, v)
			}
		}
	}
	{
		freshmen := []person{1, 2, 3, 4, 5, 6, 7}
		seniors := []person{1, 2, 3}
		d := MakeData(freshmen, seniors)

		luckyCount := 0
		baseDrawTimes := 1
		seniorsPairedMax := 3

		if len(d.results) != 0 {
			t.Error("results is not empty")
		}
		if !slices.Equal(d.freshmen, freshmen) {
			t.Error("d.freshmen is not equal to freshmen")
		}
		if !slices.Equal(d.seniors, seniors) {
			t.Error("d.seniors is not equal to seniors")
		}
		if !slices.Equal(d.pairableSeniors, seniors) {
			t.Error("d.pairableSeniors is not equal to freshmen")
		}
		if d.waitingFreshmenCount != len(freshmen) {
			t.Errorf("seniorsPairedMax is %v instead of %v", d.waitingFreshmenCount, len(freshmen))
		}
		if d.luckyCount != luckyCount {
			t.Errorf("luckyCount is %v instead of %v", d.luckyCount, luckyCount)
		}
		if d.baseDrawTimes != baseDrawTimes {
			t.Errorf("baseDrawTimes is %v instead of %v", d.baseDrawTimes, baseDrawTimes)
		}
		if d.seniorsPairedMax != seniorsPairedMax {
			t.Errorf("seniorsPairedMax is %v instead of %v", d.seniorsPairedMax, seniorsPairedMax)
		}
		for k, v := range d.seniorsPairedCount {
			if v != 0 {
				t.Errorf("seniorsPairedCount of %v is %v instead of 0", k, v)
			}
		}
	}
}
