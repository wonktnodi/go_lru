package go_lru

import (
	"testing"
    "fmt"
    "math/rand"
    "time"
)

// Test that Peek doesn't update recent-ness
func Test2Q_Peek(t *testing.T) {
	l, err := New2Q(2, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add("1", 1)
	l.Add("2", 2)
	if v, ok := l.Peek("1"); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Add("3", 3)
	if l.Contains("1") {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

func Test2Q_RandomOps(t *testing.T) {
	size := 128
	l, err := New2Q(128, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	n := 200000
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", rand.Int63() % 512)
		r := rand.Int63()
		switch r % 3 {
		case 0:
			l.Add(key, key)
		case 1:
			l.Get(key)
		case 2:
			l.Remove(key)
		}

		if l.recent.Len()+l.frequent.Len() > size {
			t.Fatalf("bad: recent: %d freq: %d",
				l.recent.Len(), l.frequent.Len())
		}
	}
}

func Test2Q_Get_RecentToFrequent(t *testing.T) {
	l, err := New2Q(128, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Touch all the entries, should be in t1
	for i := 0; i < 128; i++ {
		l.Add(fmt.Sprintf("%d", i), i)
	}
	if n := l.recent.Len(); n != 128 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}

	// Get should upgrade to t2
	for i := 0; i < 128; i++ {
		_, ok := l.Get(fmt.Sprintf("%d", i))
		if !ok {
			t.Fatalf("missing: %d", i)
		}
	}
	if n := l.recent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 128 {
		t.Fatalf("bad: %d", n)
	}

	// Get be from t2
	for i := 0; i < 128; i++ {
		_, ok := l.Get(fmt.Sprintf("%d", i))
		if !ok {
			t.Fatalf("missing: %d", i)
		}
	}
	if n := l.recent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 128 {
		t.Fatalf("bad: %d", n)
	}
}

func Test2Q_Add_RecentToFrequent(t *testing.T) {
	l, err := New2Q(128, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Add initially to recent
	l.Add("1", 1)
	if n := l.recent.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}

	// Add should upgrade to frequent
	l.Add("1", 1)
	if n := l.recent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}

	// Add should remain in frequent
	l.Add("1", 1)
	if n := l.recent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}
}

func Test2Q_Add_RecentEvict(t *testing.T) {
	l, err := New2Q(4, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Add 1,2,3,4,5 -> Evict 1
	l.Add("1", 1)
	l.Add("2", 2)
	l.Add("3", 3)
	l.Add("4", 4)
	l.Add("5", 5)
	if n := l.recent.Len(); n != 4 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.recentEvict.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 0 {
		t.Fatalf("bad: %d", n)
	}

	// Pull in the recently evicted
	l.Add("1", 1)
	if n := l.recent.Len(); n != 3 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.recentEvict.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}

	// Add 6, should cause another recent evict
	l.Add("6", 6)
	if n := l.recent.Len(); n != 3 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.recentEvict.Len(); n != 2 {
		t.Fatalf("bad: %d", n)
	}
	if n := l.frequent.Len(); n != 1 {
		t.Fatalf("bad: %d", n)
	}
}

func Test2Q(t *testing.T) {
	l, err := New2Q(128, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Add(fmt.Sprintf("%d", i), i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || fmt.Sprint(v) != k || v.(int) != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(fmt.Sprint(i))
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(fmt.Sprint(i))
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		l.Remove(fmt.Sprint(i))
		_, ok := l.Get(fmt.Sprint(i))
		if ok {
			t.Fatalf("should be deleted")
		}
	}

    l.AddWithExpire(fmt.Sprint(256), 256, time.Second * 2)

    if v, ok := l.Get(fmt.Sprint(256)); !ok  || fmt.Sprint(v) != "256" || v.(int) != 256 {
        t.Fatalf("bak key: %v", 256)
    }

    time.Sleep(2 * time.Second)
    if _, ok := l.Get(fmt.Sprint(256)); ok {
        t.Fatalf("bak key: %v", 256)
    }

	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get("200"); ok {
		t.Fatalf("should contain nothing")
	}
}

// Test that Contains doesn't update recent-ness
func Test2Q_Contains(t *testing.T) {
	l, err := New2Q(2, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add("1", 1)
	l.Add("2", 2)
	if !l.Contains("1") {
		t.Errorf("1 should be contained")
	}

	l.Add("3", 3)
	if l.Contains("1") {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}


