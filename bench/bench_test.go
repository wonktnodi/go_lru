package bench

import (
    "testing"
    "fmt"
    "github.com/wonktnodi/go_lru"
    "math/rand"
)

func BenchmarkLRU_Rand(b *testing.B) {
    l, err := go_lru.New(8192, go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        trace[i] = fmt.Sprintf("%d", rand.Int63()%32768)
    }

    b.ResetTimer()
    b.ReportAllocs()

    var hit, miss int
    for i := 0; i < 2*b.N; i++ {
        if i%2 == 0 {
            l.Add(trace[i], trace[i])
        } else {
            _, ok := l.Get(trace[i])
            if ok {
                hit++
            } else {
                miss++
            }
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkLRU_Freq(b *testing.B) {
    l, err := go_lru.New(8192, go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        if i%2 == 0 {
            trace[i] = fmt.Sprintf("%d", rand.Int63()%16384)
        } else {
            trace[i] = fmt.Sprintf("%d", rand.Int63()%32768)
        }
    }

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        l.Add(trace[i], trace[i])
    }
    var hit, miss int
    for i := 0; i < b.N; i++ {
        _, ok := l.Get(trace[i])
        if ok {
            hit++
        } else {
            miss++
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func Benchmark2Q_Rand(b *testing.B) {
    l, err := go_lru.New2Q(8192, go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        trace[i] = fmt.Sprint(rand.Int63() % 32768)
    }

    b.ResetTimer()
    b.ReportAllocs()

    var hit, miss int
    for i := 0; i < 2*b.N; i++ {
        if i%2 == 0 {
            l.Add(trace[i], trace[i])
        } else {
            _, ok := l.Get(trace[i])
            if ok {
                hit++
            } else {
                miss++
            }
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func Benchmark2Q_Freq(b *testing.B) {
    l, err := go_lru.New2Q(8192, go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        if i%2 == 0 {
            trace[i] = fmt.Sprint(rand.Int63() % 16384)
        } else {
            trace[i] = fmt.Sprint(rand.Int63() % 32768)
        }
    }

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        l.Add(trace[i], trace[i])
    }
    var hit, miss int
    for i := 0; i < b.N; i++ {
        _, ok := l.Get(trace[i])
        if ok {
            hit++
        } else {
            miss++
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkARC_Rand(b *testing.B) {
    l, err := go_lru.NewARC(8192, go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        trace[i] = fmt.Sprintf("%d", rand.Int63()%32768)
    }

    b.ResetTimer()
    b.ReportAllocs()

    var hit, miss int
    for i := 0; i < 2*b.N; i++ {
        if i%2 == 0 {
            l.Add(trace[i], trace[i])
        } else {
            _, ok := l.Get(trace[i])
            if ok {
                hit++
            } else {
                miss++
            }
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkARC_Freq(b *testing.B) {
    l, err := go_lru.NewARC(8192,go_lru.NoExpiration)
    if err != nil {
        b.Fatalf("err: %v", err)
    }

    trace := make([]string, b.N*2)
    for i := 0; i < b.N*2; i++ {
        if i%2 == 0 {
            trace[i] = fmt.Sprintf("%d", rand.Int63()%16384)
        } else {
            trace[i] = fmt.Sprintf("%d", rand.Int63()%32768)
        }
    }

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        l.Add(trace[i], trace[i])
    }
    var hit, miss int
    for i := 0; i < b.N; i++ {
        _, ok := l.Get(trace[i])
        if ok {
            hit++
        } else {
            miss++
        }
    }
    b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}
