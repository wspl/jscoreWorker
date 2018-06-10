package main

import (
	"github.com/wspl/jscoreWorker"
	"fmt"
)

func main() {
	worker := jscoreWorker.New(func(msg []byte) []byte {
		fmt.Println(msg)
		return nil
	})
	worker.Load("codeWithRecv.js", `
		JSCoreWorker.recv(function(msg) {
			JSCoreWorker.print('CJK 中文')
			JSCoreWorker.send(msg)
		});
	`)
	worker.SendBytes([]byte("abcd"))
}