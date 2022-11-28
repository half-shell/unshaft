package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const SIMULATIONS_NUM = 9
const ROWS_COUNT = 8
const COL_COUNT = 8
const CELL_NUM = ROWS_COUNT * COL_COUNT
const MINE_NUM = 8

// We can define a shaft only by the index of the mines it contains
type Shaft [MINE_NUM]int
type Shafts map[string]Shaft

type Stats struct {
	Ones   int
	Twos   int
	Threes int
	Fours  int
}

func (s Shaft) IsValid() bool {
	for x := range s {
		for y := range s {
			if x != y && s[x] == s[y] {
				return false
			}
		}
	}

	return true
}

func GenerateShaft() Shaft {
	shaft := *new(Shaft)

	for {
		// FIXME(half-shell): Minor improvement; we can check unicity
		// in this loop instead of looping all around again
		// It might become somewhat impactful for a large number of
		// generated shafts that have a much higher chance of collision
		for i := 0; i < MINE_NUM; i++ {
			shaft[i] = rand.Intn(CELL_NUM)
		}

		sort.Sort(sort.IntSlice(shaft[:]))

		if shaft.IsValid() {
			return shaft
		}
	}
}

func CreateUniqueShaft(shafts *Shafts) Shaft {
	for {
		shaft := GenerateShaft()
		hash := shaft.Hash()

		if _, ok := (*shafts)[hash]; !ok {
			(*shafts)[hash] = shaft
			return shaft
		}
	}
}

func (s Shaft) Hash() string {
	return fmt.Sprintf("%v%v%v%v%v%v%v%v", s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7])
}
func GetNeighbouringIndexes(i int) (results []int) {
	// Is it on the first column?
	if math.Mod(float64(i), 8) == 0 {
		operations := []int{-COL_COUNT, -COL_COUNT + 1, 1, COL_COUNT, COL_COUNT + 1}
		for _, o := range operations {
			n := i + o
			if n > 0 {
				results = append(results, i+o)
			}
		}

		return results
	}

	// Is it on the last column?
	if math.Mod(float64(i), 8) == 7 {
		operations := []int{-COL_COUNT, -COL_COUNT - 1, -1, COL_COUNT, COL_COUNT - 1}
		for _, o := range operations {
			n := i + o
			if n > 0 {
				results = append(results, i+o)
			}
		}

		return results
	}

	// Otherwise
	operations := []int{
		-COL_COUNT - 1,
		-COL_COUNT,
		-COL_COUNT + 1,
		-1,
		1,
		COL_COUNT - 1,
		COL_COUNT,
		COL_COUNT + 1,
	}
	for _, o := range operations {
		n := i + o
		if n > 0 {
			results = append(results, i+o)
		}
	}

	return results

}

func incOnes(stats *Stats, s *Shaft, index int) {
	neighbours := GetNeighbouringIndexes(index)

	for i := range s {
		for n := range neighbours {
			if s[i] == neighbours[n] {
				(*stats).Ones++
			}
		}
	}
}

// TODO(half-shell): A fan-in distribusion might make sense in this case
func ProcessStats(s Shaft) []Stats {
	var wg sync.WaitGroup

	stats := make([]Stats, CELL_NUM)

	for i := 0; i < CELL_NUM; i++ {
		// If index matches a mine index, ignore
		// FIXME(half-shell): This is pretty awkward.
		// We might want to find a better way to process this check
		for y := range s {
			if s[y] == i {
				continue
			}
		}

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			incOnes(&stats[i], &s, i)
		}(i)
	}

	wg.Wait()

	return stats
}

func main() {
	shafts := make(Shafts, SIMULATIONS_NUM)

	for i := 0; i <= SIMULATIONS_NUM; i++ {
		CreateUniqueShaft(&shafts)
	}

	for s := range shafts {
		stats := ProcessStats(shafts[s])
		log.Printf("Stats for shaft %v: %v\n", shafts[s], stats)
	}
}
