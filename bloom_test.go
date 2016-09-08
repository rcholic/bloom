package bloom

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

//FilterSize  test filter size, for uniformization between tests
const FilterSize = 512

var tests = []string{
	"Bloom",
	"Poney",
	"GitHub",
	"Pwet",
	"Toto",
	"Billy",
	"Jacob",
	"Omelette",
	"De",
	"Fromage",
}

// TestEmptyQuery : Assert that empty bloom filter Match always return false
func TestEmptyQuery(t *testing.T) {
	bf := New(FilterSize, 5)
	for _, v := range tests {
		if bf.Match(v) {
			t.Errorf("Empty Filter return true when matching : %s", v)
		}
	}
}

// Test that inserted elements return true upon Match
func TestMatch(t *testing.T) {
	bf := New(FilterSize, 5)
	for _, v := range tests {
		bf.Feed(v)
		if !bf.Match(v) {
			t.Errorf("Filter match return false on inserted element : %s", v)
		}
	}
}

// Test JSON Export / Import
func TestSerialization(t *testing.T) {
	// TODO
}

// TestMerge : Test the merge functionnality
func TestMerge(t *testing.T) {
	bf := New(FilterSize, 5)
	oth := New(FilterSize, 5)

	bf.Feed("foo")
	for _, v := range tests {
		oth.Feed(v)
	}

	bf.Merge(oth)
	for _, v := range tests {
		if !bf.Match(v) {
			t.Errorf("Element merged in filter was not found : %s", v)
		}
	}
	if !bf.Match("foo") {
		t.Errorf("Original element deleted upon merge")

	}
}

func TestFalseNegative(t *testing.T) {
	bf := New(FilterSize, 5)
	for _, v := range tests {
		bf.Feed(v)
		if !bf.Match(v) {
			t.Errorf("Element inserted in filter was not found : %s", v)
		}
	}
	// fmt.Printf("%v\n", bf.fillRatio())
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func testScalable(t *testing.T, ex float64, c int, f int) {
	file, _ := os.Open("/usr/share/dict/ngerman")
	defer file.Close()

	bf := NewDefaultScalable(ex)
	scanner := bufio.NewScanner(file)
	for i := 0; i < c && scanner.Scan(); i++ {
		v := scanner.Text()
		bf.Feed(v)
	}
	fp := 0
	total := 0
	for i := 0; i < f; i++ {
		total++
		if bf.Match(randString(40)) {
			fp++
		}
	}
	rate := float64(fp) / float64(total)
	fmt.Printf("For %d pass: fp: [%d], rate %f\n", total, fp, rate)
	if rate > ex {
		t.Errorf("Unrespected error rate : %f > %f", rate, ex)
	}
}

func TestScalableLinuxWords(t *testing.T) {
	// t.Log("Testing ScalableFilter")
	// expected := 0.1
	// c := 100000000
	// f := 1000000
	// for expected > 0.001 {
	// 	testScalable(t, expected, c, f)
	// 	expected /= 10.0
	// }
}

/*
 *
 *
 *  BENCHMARKS
 *
 *
 */

// CREATION

func BenchCreation(hashes uint64, b *testing.B) {
	for n := 0; n < b.N; n++ {
		New(FilterSize, 5)
	}
}

func BenchmarkCreation1(b *testing.B)  { BenchCreation(1, b) }
func BenchmarkCreation5(b *testing.B)  { BenchCreation(5, b) }
func BenchmarkCreation10(b *testing.B) { BenchCreation(10, b) }

// INSERTION

func BenchInsert(hashes uint64, b *testing.B) {
	bf := New(FilterSize, hashes)
	for n := 0; n < b.N; n++ {
		bf.Feed(randString(20))
	}
}

func BenchmarkInsert1(b *testing.B)  { BenchInsert(1, b) }
func BenchmarkInsert5(b *testing.B)  { BenchInsert(5, b) }
func BenchmarkInsert10(b *testing.B) { BenchInsert(10, b) }

// MATCHING

func BenchMatch(hashes uint64, b *testing.B) {
	bf := New(FilterSize, hashes)
	bf.Feed("I am a test string")
	for n := 0; n < b.N; n++ {
		bf.Match("I am a test string")
	}
}

func BenchmarkMatch1(b *testing.B)  { BenchMatch(1, b) }
func BenchmarkMatch5(b *testing.B)  { BenchMatch(5, b) }
func BenchmarkMatch10(b *testing.B) { BenchMatch(10, b) }

// INTERNAL ROUTINE : isSet

func BenchIsSet(b *testing.B) {
	bf := New(FilterSize, 5)
	for n := 0; n < b.N; n++ {
		bf.isSet(42)
	}
}

func BenchmarkIsSet(b *testing.B) { BenchIsSet(b) }

// Internal routing : FillRatio

func BenchFillRatio(n uint64, b *testing.B) {
	bf := New(FilterSize*n, 5)
	for n := 0; n < b.N; n++ {
		bf.FillRatio()
	}
}

func BenchmarkFillRatio1(b *testing.B)  { BenchFillRatio(1, b) }
func BenchmarkFillRatio5(b *testing.B)  { BenchFillRatio(5, b) }
func BenchmarkFillRatio10(b *testing.B) { BenchFillRatio(10, b) }