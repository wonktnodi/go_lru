package go_lru

import (
    "testing"
    "fmt"
)

func TestBaseLRU(t *testing.T) {
    evictCounter := 0
    onEvicted := func(k string, v interface{}) {
        if k != fmt.Sprintf("%d", v.(int)) {
            t.Fatalf("Evict values not equal (%v!=%v)", k, v)
        }
        evictCounter += 1
    }
    l, err := NewBaseLRU(128, onEvicted, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }

    for i := 0; i < 256; i++ {
        l.Add(fmt.Sprintf("%d", i), i)
    }
    if l.Len() != 128 {
        t.Fatalf("bad len: %v", l.Len())
    }

    if evictCounter != 128 {
        t.Fatalf("bad evict count: %v", evictCounter)
    }

    keys := l.Keys()
    for i, k := range keys {
        if v, ok := l.Get(k); !ok || fmt.Sprintf("%d", v.(int)) != k || v.(int) != i+128 {
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
        ok := l.Remove(fmt.Sprintf("%d", i))
        if !ok {
            t.Fatalf("should be contained")
        }
        ok = l.Remove(fmt.Sprintf("%d", i))
        if ok {
            t.Fatalf("should not be contained")
        }
        _, ok = l.Get(fmt.Sprintf("%d", i))
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

    l.Purge()
    if l.Len() != 0 {
        t.Fatalf("bad len: %v", l.Len())
    }
    if _, ok := l.Get("200"); ok {
        t.Fatalf("should contain nothing")
    }
}

func TestBaseLRU_GetOldest_RemoveOldest(t *testing.T) {
    l, err := NewBaseLRU(128, nil, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }
    for i := 0; i < 256; i++ {
        l.Add(fmt.Sprintf("%d", i), i)
    }
    k, _, ok := l.GetOldest()
    if !ok {
        t.Fatalf("missing")
    }
    if k != "128" {
        t.Fatalf("bad: %v", k)
    }

    k, _, ok = l.RemoveOldest()
    if !ok {
        t.Fatalf("missing")
    }
    if k != "128" {
        t.Fatalf("bad: %v", k)
    }

    k, _, ok = l.RemoveOldest()
    if !ok {
        t.Fatalf("missing")
    }
    if k != "129" {
        t.Fatalf("bad: %v", k)
    }
}

// Test that Add returns true/false if an eviction occurred
func TestBaseLRU_Add(t *testing.T) {
    evictCounter := 0
    onEvicted := func(k string, v interface{}) {
        evictCounter += 1
    }

    l, err := NewBaseLRU(1, onEvicted, NoExpiration)
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

// Test that Contains doesn't update recent-ness
func TestBaseLRU_Contains(t *testing.T) {
    l, err := NewBaseLRU(2, nil, NoExpiration)
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

// Test that Peek doesn't update recent-ness
func TestBaseLRU_Peek(t *testing.T) {
    l, err := NewBaseLRU(2, nil, NoExpiration)
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
