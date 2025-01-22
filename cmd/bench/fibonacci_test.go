package main

import (
	"sort"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func BenchmarkFibo(b *testing.B) {
	count := 20

	b.Run("recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FiboRecursive(count)
		}
	})
	b.Run("optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FiboOptimized(count)
		}
	})

}

func BenchmarkSortSlice(b *testing.B) {
	rand.Seed(uint64(time.Now().UnixNano()))
	b.ResetTimer()
	var cmps int64

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		slice := make([]int, 10000)
		// предзаполняем слайс
		for i := 0; i < len(slice); i++ {
			slice[i] = rand.Intn(1000)
		}
		b.StartTimer()

		// сортируем
		sort.Slice(slice, func(i, j int) bool {
			cmps++
			return slice[i] < slice[j]
		})
	}

	// add metric
	b.ReportMetric(float64(b.N), "iterations")
	// добавляем метрику
	b.ReportMetric(float64(cmps)/float64(b.N), "compares/op")
}
