package go_lru

import (
    "testing"
    "time"
    "fmt"
    "math/rand"
)

func init() {
    rand.Seed(time.Now().Unix())
}

func TestARC_RandomOps(t *testing.T) {
    size := 128
    l, err := NewARC(128, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }

    n := 200000
    for i := 0; i < n; i++ {
        key := fmt.Sprintf("%d", rand.Int63()%512)
        r := rand.Int63()
        switch r % 3 {
        case 0:
            l.Add(key, key)
        case 1:
            l.Get(key)
        case 2:
            l.Remove(key)
        }

        if l.t1.Len()+l.t2.Len() > size {
            t.Fatalf("bad: t1: %d t2: %d b1: %d b2: %d p: %d",
                l.t1.Len(), l.t2.Len(), l.b1.Len(), l.b2.Len(), l.p)
        }
        if l.b1.Len()+l.b2.Len() > size {
            t.Fatalf("bad: t1: %d t2: %d b1: %d b2: %d p: %d",
                l.t1.Len(), l.t2.Len(), l.b1.Len(), l.b2.Len(), l.p)
        }
    }
}

func TestARC_Get_RecentToFrequent(t *testing.T) {
    l, err := NewARC(128, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }

    // Touch all the entries, should be in t1
    for i := 0; i < 128; i++ {
        l.Add(fmt.Sprintf("%d", i), i)
    }
    if n := l.t1.Len(); n != 128 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }

    // Get should upgrade to t2
    for i := 0; i < 128; i++ {
        _, ok := l.Get(fmt.Sprint(i))
        if !ok {
            t.Fatalf("missing: %d", i)
        }
    }
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 128 {
        t.Fatalf("bad: %d", n)
    }

    // Get be from t2
    for i := 0; i < 128; i++ {
        _, ok := l.Get(fmt.Sprint(i))
        if !ok {
            t.Fatalf("missing: %d", i)
        }
    }
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 128 {
        t.Fatalf("bad: %d", n)
    }
}

func TestARC_Add_RecentToFrequent(t *testing.T) {
    l, err := NewARC(128, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }

    // Add initially to t1
    l.Add("1", 1)
    if n := l.t1.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }

    // Add should upgrade to t2
    l.Add("1", 1)
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }

    // Add should remain in t2
    l.Add("1", 1)
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }
}

func TestARC_Adaptive(t *testing.T) {
    l, err := NewARC(4, NoExpiration)
    if err != nil {
        t.Fatalf("err: %v", err)
    }

    // Fill t1
    for i := 0; i < 4; i++ {
        l.Add(fmt.Sprintf("%d", i), i)
    }
    if n := l.t1.Len(); n != 4 {
        t.Fatalf("bad: %d", n)
    }

    // Move to t2
    l.Get("0")
    l.Get("1")
    if n := l.t2.Len(); n != 2 {
        t.Fatalf("bad: %d", n)
    }

    // Evict from t1
    l.Add("4", 4)
    if n := l.b1.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }

    // Current state
    // t1 : (MRU) [4, 3] (LRU)
    // t2 : (MRU) [1, 0] (LRU)
    // b1 : (MRU) [2] (LRU)
    // b2 : (MRU) [] (LRU)

    // Add 2, should cause hit on b1
    l.Add("2", 2)
    if n := l.b1.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }
    if l.p != 1 {
        t.Fatalf("bad: %d", l.p)
    }
    if n := l.t2.Len(); n != 3 {
        t.Fatalf("bad: %d", n)
    }

    // Current state
    // t1 : (MRU) [4] (LRU)
    // t2 : (MRU) [2, 1, 0] (LRU)
    // b1 : (MRU) [3] (LRU)
    // b2 : (MRU) [] (LRU)

    // Add 4, should migrate to t2
    l.Add("4", 4)
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 4 {
        t.Fatalf("bad: %d", n)
    }

    // Current state
    // t1 : (MRU) [] (LRU)
    // t2 : (MRU) [4, 2, 1, 0] (LRU)
    // b1 : (MRU) [3] (LRU)
    // b2 : (MRU) [] (LRU)

    // Add 4, should evict to b2
    l.Add("5", 5)
    if n := l.t1.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 3 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.b2.Len(); n != 1 {
        t.Fatalf("bad: %d", n)
    }

    // Current state
    // t1 : (MRU) [5] (LRU)
    // t2 : (MRU) [4, 2, 1] (LRU)
    // b1 : (MRU) [3] (LRU)
    // b2 : (MRU) [0] (LRU)

    // Add 0, should decrease p
    l.Add("0", 0)
    if n := l.t1.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.t2.Len(); n != 4 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.b1.Len(); n != 2 {
        t.Fatalf("bad: %d", n)
    }
    if n := l.b2.Len(); n != 0 {
        t.Fatalf("bad: %d", n)
    }
    if l.p != 0 {
        t.Fatalf("bad: %d", l.p)
    }

    // Current state
    // t1 : (MRU) [] (LRU)
    // t2 : (MRU) [0, 4, 2, 1] (LRU)
    // b1 : (MRU) [5, 3] (LRU)
    // b2 : (MRU) [0] (LRU)
}

func TestARC(t *testing.T) {
    l, err := NewARC(128, NoExpiration)
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

    l.AddWithExpire(fmt.Sprint(256), 256, time.Second * 2)
    if l.Len() != 128 {
        t.Fatalf("bad len: %v", l.Len())
    }
    if v, ok := l.Get(fmt.Sprint(256)); !ok  || fmt.Sprint(v) != "256" || v.(int) != 256 {
        t.Fatalf("bak key: %v", 256)
    }

    time.Sleep(2 * time.Second)
    if _, ok := l.Get(fmt.Sprint(256)); ok {
        t.Fatalf("bak key: %v", 256)
    }



    for i := 128; i < 192; i++ {
        l.Remove(fmt.Sprintf("%d", i))
        _, ok := l.Get(fmt.Sprintf("%d", i))
        if ok {
            t.Fatalf("should be deleted")
        }
    }

    l.Purge()
    if l.Len() != 0 {
        t.Fatalf("bad len: %v", l.Len())
    }
    if _, ok := l.Get(fmt.Sprint(200)); ok {
        t.Fatalf("should contain nothing")
    }
}

// Test that Contains doesn't update recent-ness
func TestARC_Contains(t *testing.T) {
    l, err := NewARC(2, NoExpiration)
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
func TestARC_Peek(t *testing.T) {
    l, err := NewARC(2, NoExpiration)
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
