go-lru
==========

This project is forked from github.com/hashicorp/golang-lru. This provides the `lru` package which implements a fixed-size
thread safe LRU cache. It is based on the cache in Groupcache.

Example
=======

Using the LRU is very simple:

```go
l, _ := New(128)
for i := 0; i < 256; i++ {
    l.Add(i, nil)
}
if l.Len() != 128 {
    panic(fmt.Sprintf("bad len: %v", l.Len()))
}
```
