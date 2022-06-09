package id2

import (
	"sync"
	"testing"
)

func TestGenID(t *testing.T) {
	id := GenID()
	_ = id
}

func TestGenUpperID(t *testing.T) {
	id := GenUpperID()
	_ = id
}

func TestIDCollision(t *testing.T) {
	const (
		parallel  = 100
		magnitude = 1000
	)

	wait := sync.WaitGroup{}
	wait.Add(parallel)

	storage := make([]map[string]struct{}, parallel)

	for i := 0; i < parallel; i++ {
		storage[i] = make(map[string]struct{}, magnitude)
		go func(wait *sync.WaitGroup, sg map[string]struct{}) {
			for i := 0; i < magnitude; i++ {
				sg[GenID()] = struct{}{}
			}
			wait.Done()
		}(&wait, storage[i])
	}

	wait.Wait()

	set := make(map[string]struct{}, parallel*magnitude)
	for _, sg := range storage {
		if len(sg) != magnitude {
			t.Fatal("id conflict [", len(sg), "]")
			return
		}
		for id := range sg {
			if _, ok := set[id]; ok {
				t.Fatal("id conflict")
				return
			}
			set[id] = struct{}{}
		}
	}
}

// go test -benchmem  -count=10 -bench Benchmark
func BenchmarkTestGenID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenID()
	}
}
