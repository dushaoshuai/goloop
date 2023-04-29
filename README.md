# goloop

goloop tries to facilitate looping in Go.
It imitates Go's "for ... range ... {}" looping style.

# Download/Install

```shell
go get -u github.com/dushaoshuai/goloop
```

# Quick Start

If you are tired of writing this trivial code:

```go
for i := 0; i < 10; i++ {
	fmt.Println(i)
}
```

try goloop:

```go
for i := range goloop.Repeat(10) {
	fmt.Println(i)
}
```

if you want to break the loop when certain conditions are met:

```go
for i := range goloop.RepeatWithBreak(10) {
	fmt.Println(i.I)
	if i.I == 5 {
		i.Break()
	}
}
```
