package go_lru

import (
	"testing"
    "fmt"
    "time"
)

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k string, v interface{}) {
		if k != fmt.Sprint(v) {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter += 1
	}
	l, err := NewWithEvict(128, NoExpiration, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Add(fmt.Sprintf("%d", i), i,)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || fmt.Sprint(v) != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(fmt.Sprintf("%d", i))
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(fmt.Sprintf("%d", i))
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		l.Remove(fmt.Sprintf("%d", i))
		_, ok := l.Get(fmt.Sprintf("%d", i))
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get("192") // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (i < 63 && k != fmt.Sprintf("%d", i+193)) || (i == 63 && k != "192") {
			t.Fatalf("out of order key: %v", k)
		}
	}

	l.AddWithExpire("256", 256, time.Second * 2)
    if v, ok := l.Get("256"); !ok  || fmt.Sprint(v) != "256" || v.(int) != 256 {
        t.Fatalf("bak key: %v", 256)
    }
    time.Sleep(2 * time.Second)
    if  _, ok := l.Get("256"); ok {
        t.Fatalf("sould be evicted by timeout")
    }

	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get("200"); ok {
		t.Fatalf("should contain nothing")
	}
}

// test that Add returns true/false if an eviction occurred
func TestLRUAdd(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k string, v interface{}) {
		evictCounter += 1
	}

	l, err := NewWithEvict(1, NoExpiration, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add("1", 1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add("2", 2) == false || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// test that Contains doesn't update recent-ness
func TestLRUContains(t *testing.T) {
	l, err := New(2, NoExpiration)
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

// test that Contains doesn't update recent-ness
func TestLRUContainsOrAdd(t *testing.T) {
	l, err := New(2, NoExpiration)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add("1", 1)
	l.Add("2", 2)
	contains, evict := l.ContainsOrAdd("1", 1)
	if !contains {
		t.Errorf("1 should be contained")
	}
	if evict {
		t.Errorf("nothing should be evicted here")
	}

	l.Add("3", 3)
	contains, evict = l.ContainsOrAdd("1", 1)
	if contains {
		t.Errorf("1 should not have been contained")
	}
	if !evict {
		t.Errorf("an eviction should have occurred")
	}
	if !l.Contains("1") {
		t.Errorf("now 1 should be contained")
	}
}

// test that Peek doesn't update recent-ness
func TestLRUPeek(t *testing.T) {
	l, err := New(2, NoExpiration)
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