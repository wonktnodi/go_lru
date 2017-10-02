package main

import (
    "github.com/wonktnodi/go_lru"
    "fmt"
    "time"
)

func main() {
    test_2q()
}

func test_2q() {
    l, err := go_lru.New2Q(128, go_lru.NoExpiration)
    if err != nil {
        fmt.Errorf("err: %v", err)
    }

    l.AddWithExpire(fmt.Sprint(256), 256, time.Second * 2)

    if v, ok := l.Get(fmt.Sprint(256)); !ok  || fmt.Sprint(v) != "256" || v.(int) != 256 {
        fmt.Errorf("bak key: %v", 256)
    }

    time.Sleep(2 * time.Second)
    if _, ok := l.Get(fmt.Sprint(256)); ok {
        fmt.Errorf("bak key: %v", 256)
    }
}

func test_arc() {
    l, err := go_lru.NewARC(128, go_lru.NoExpiration)
    if err != nil {
        fmt.Errorf("err: %v", err)
    }

    l.AddWithExpire(fmt.Sprint(256), 256, time.Second*2)
    if l.Len() != 128 {
        fmt.Errorf("bad len: %v", l.Len())
    }
    if v, ok := l.Get(fmt.Sprint(256)); !ok || fmt.Sprint(v) != "256" || v.(int) != 256 {
        fmt.Errorf("bak key: %v", 256)
    }
    time.Sleep(3 * time.Second)
    if _, ok := l.Get(fmt.Sprint(256)); ok {
        fmt.Errorf("bak key: %v", 256)
    }
}