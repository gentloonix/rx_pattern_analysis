package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/juliangruber/go-intersect"
	"golang.org/x/exp/constraints"
)

type Pair struct {
	Key   int
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ReadInts reads whitespace-separated ints from r. If there's an error, it
// returns the ints successfully read so far as well as the error value.
func ReadInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

func main() {
	slices := 16
	chunksTrimmedLimit := (80 << 20 * 8 * 8) / 512

	file, err := os.Open("rnd/read.txt")
	if err != nil {
		log.Println(err)
	}

	ints, err := ReadInts(file)
	if err != nil {
		log.Println(err)
	}

	for i := range ints {
		ints[i] /= 64
	}

	mIntTotalCount := map[int]int{}
	for _, v := range ints {
		mIntTotalCount[v]++
	}

	intTotalCount := make(PairList, len(mIntTotalCount))
	i := 0
	for k, v := range mIntTotalCount {
		intTotalCount[i] = Pair{k, v}
		i++
	}

	var intersections []interface{}
	for s := 0; s < slices; s++ {
		sInts := ints[len(ints)/slices*s : len(ints)/slices*(s+1)]

		msIntCount := map[int]int{}
		for _, v := range sInts {
			msIntCount[v]++
		}

		sIntCount := make(PairList, len(msIntCount))
		i := 0
		for k, v := range msIntCount {
			sIntCount[i] = Pair{k, v}
			i++
		}
		sort.Sort(sort.Reverse(sIntCount))

		chunksTrimmedCount := min(chunksTrimmedLimit, len(sIntCount))
		chunksTrimmed := sIntCount[:chunksTrimmedCount]
		chunksHit := 0
		for _, v := range chunksTrimmed {
			chunksHit += v.Value
		}
		log.Printf("slice %d", s)
		log.Printf("chunks trimmed count %d", chunksTrimmedCount)
		log.Printf("chunks trimmed hit %d of total %d, %f%%", chunksHit, len(sInts), float64(chunksHit)/float64(len(sInts))*100)
		log.Printf("chunks trimmed top 100 %d", chunksTrimmed[:min(100, len(chunksTrimmed))])

		chunksTrimmedRaw := make([]interface{}, len(chunksTrimmed))
		i = 0
		for _, v := range chunksTrimmed {
			chunksTrimmedRaw[i] = v.Key
			i++
		}
		if intersections == nil {
			intersections = chunksTrimmedRaw
		} else {
			intersections = intersect.Hash(intersections, chunksTrimmedRaw)
		}
	}
	log.Printf("intersections count %d of total %d, %f%%", len(intersections), len(intTotalCount), float64(len(intersections))/float64(len(intTotalCount))*100)
	log.Printf("intersections top 100 %d", intersections[:min(100, len(intersections))])
}
