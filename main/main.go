package main

import (
	. "github.com/wspl/jscoreWorker"
	"fmt"
)

func main() {
	recvCount := 0
	worker := New(func(msg []byte) []byte {
		println(msg)
		if len(msg) != 5 {
			fmt.Println("bad msg", msg)
		}
		recvCount++
		return nil
	})

	err := worker.Load("codeWithRecv.js", `
		JSCoreWorker.recv(function(msg) {
			JSCoreWorker.print("TestBasic recv byteLength", msg.byteLength);
			if (msg.byteLength !== 3) {
				throw Error("bad message");
			}
		});
	`)
	if err != nil {
		fmt.Println(err)
	}
	err = worker.SendBytes([]byte("hii"))
	if err != nil {
		fmt.Println(err)
	}
	codeWithSend := `
		JSCoreWorker.send(new ArrayBuffer(5));
		JSCoreWorker.send(new ArrayBuffer(5));
	`
	err = worker.Load("codeWithSend.js", codeWithSend)
	if err != nil {
		fmt.Println(err)
	}

	if recvCount != 2 {
		fmt.Println("bad recvCount", recvCount)
	}
}