# goloop

goloop tries to facilitate looping in Go.

# Download/Install

```shell
go get -u github.com/dushaoshuai/goloop
```

# Quick Start

Replace this trivial code:

```go
for i := 0; i < 10; i++ {
	fmt.Println(i)
}
```

with:

```go
for i := range goloop.Repeat(10) {
	fmt.Println(i)
}
```

Break the loop when certain conditions are met:

```go
for i := range goloop.RepeatWithBreak(10) {
	fmt.Println(i.I)
	if i.I == 5 {
		i.Break()
	}
}
```

Range over a sequence of integers:

```go
for i := range goloop.Range(3, 26, 5) {
    fmt.Println(i.I)
    if i.I >= 18 {
        i.Break()
    }
}
```

Range over a sequence of integers from a slice:

```go
for i, n := range goloop.RangeSlice[uint8](250, 255) {
	fmt.Println(i, n)
	if n >= 253 {
		break
	}
}
```
