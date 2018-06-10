package main

import (
	. "github.com/wspl/jscoreWorker"
)

func main() {
	recvCount := 0
	worker1 := New(func(msg []byte) []byte {
		if len(msg) != 5 {
			panic("bad message")
		}
		recvCount++
		return nil
	})
	worker2 := New(func(msg []byte) []byte {
		if len(msg) != 3 {
			panic("bad message")
		}
		recvCount++
		return nil
	})

	err := worker1.Load("1.js", `JSCoreWorker.send(new ArrayBuffer(5))`)
	if err != nil {
		panic(err)
	}

	err = worker2.Load("2.js", `JSCoreWorker.send(new ArrayBuffer(3))`)
	if err != nil {
		panic(err)
	}

	if recvCount != 2 {
		panic("bad recvCount")
	}
}