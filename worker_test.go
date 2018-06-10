package jscoreWorker

import (
	"strings"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	println(Version())
}

func TestSetFlags(t *testing.T) {
	// One of the V8 flags to use as a test:
	//   --lazy (use lazy compilation)
	//      type: bool  default: true
	//args := []string{"hello", "--lazy", "foobar"}
	//modified := SetFlags(args)
	//if len(modified) != 2 || modified[0] != "hello" || modified[1] != "foobar" {
	//	t.Fatalf("unexpected %v", modified)
	//}

	println("Not applicable to JavaScriptCore")
}

func TestPrint(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		t.Fatal("shouldn't recieve Message")
		return nil
	})
	err := worker.Load("code.js", `JSCoreWorker.print("ready");`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSyntaxError(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		t.Fatal("shouldn't recieve Message")
		return nil
	})

	code := `JSCoreWorker.print(hello world");`
	err := worker.Load("codeWithSyntaxError.js", code)
	errorContains(t, err, "codeWithSyntaxError.js")
	errorContains(t, err, "hello")
}

func TestSendRecv(t *testing.T) {
	recvCount := 0
	worker := New(func(msg []byte) []byte {
		println(msg)
		if len(msg) != 5 {
			t.Fatal("bad msg", msg)
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
		t.Fatal(err)
	}
	err = worker.SendBytes([]byte("hii"))
	if err != nil {
		t.Fatal(err)
	}
	codeWithSend := `
		JSCoreWorker.send(new ArrayBuffer(5));
		JSCoreWorker.send(new ArrayBuffer(5));
	`
	err = worker.Load("codeWithSend.js", codeWithSend)
	if err != nil {
		t.Fatal(err)
	}

	if recvCount != 2 {
		t.Fatal("bad recvCount", recvCount)
	}
}

func TestSendWithReturnArrayBuffer(t *testing.T) {
	recvCount := 0
	worker := New(func(msg []byte) []byte {
		if len(msg) != 123 {
			t.Fatal("unexpected message")
		}
		recvCount++
		return []byte{1, 2, 3}
	})
	err := worker.Load("TestSendWithReturnArrayBuffer.js", `
		var ret = JSCoreWorker.send(new ArrayBuffer(123));
		if (!(ret instanceof ArrayBuffer)) throw Error("bad");
		if (ret.byteLength !== 3) throw Error("bad");
		ret = new Uint8Array(ret);
		if (ret[0] !== 1) throw Error("bad");
		if (ret[1] !== 2) throw Error("bad");
		if (ret[2] !== 3) throw Error("bad");
	`)
	if err != nil {
		t.Fatal(err)
	}
	if recvCount != 1 {
		t.Fatal("bad recvCount", recvCount)
	}
}

func TestThrowDuringLoad(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		return nil
	})
	err := worker.Load("TestThrowDuringLoad.js", `throw Error("bad");`)
	errorContains(t, err, "TestThrowDuringLoad.js")
	errorContains(t, err, "bad")
}

func TestThrowInRecvCallback(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		return nil
	})
	err := worker.Load("TestThrowInRecvCallback.js", `
		JSCoreWorker.recv(function(msg) {
			throw Error("bad");
		});
	`)
	if err != nil {
		t.Fatal(err)
	}
	err = worker.SendBytes([]byte("abcd"))
	errorContains(t, err, "TestThrowInRecvCallback.js")
	errorContains(t, err, "bad")
}

func TestPrintUint8Array(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		return nil
	})
	err := worker.Load("buffer.js", `
		var uint8 = new Uint8Array(16);
		JSCoreWorker.print(uint8);
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultipleWorkers(t *testing.T) {
	recvCount := 0
	worker1 := New(func(msg []byte) []byte {
		if len(msg) != 5 {
			t.Fatal("bad message")
		}
		recvCount++
		return nil
	})
	worker2 := New(func(msg []byte) []byte {
		if len(msg) != 3 {
			t.Fatal("bad message")
		}
		recvCount++
		return nil
	})

	err := worker1.Load("1.js", `JSCoreWorker.send(new ArrayBuffer(5))`)
	if err != nil {
		t.Fatal(err)
	}

	err = worker2.Load("2.js", `JSCoreWorker.send(new ArrayBuffer(3))`)
	if err != nil {
		t.Fatal(err)
	}

	if recvCount != 2 {
		t.Fatal("bad recvCount", recvCount)
	}
}

func TestRequestFromJS(t *testing.T) {
	var captured []byte
	worker := New(func(msg []byte) []byte {
		captured = msg
		return nil
	})
	code := ` JSCoreWorker.send(new ArrayBuffer(4)); `
	err := worker.Load("code.js", code)
	if err != nil {
		t.Fatal(err)
	}
	if len(captured) != 4 {
		t.Fail()
	}
}

// Test breaking script execution
func TestWorkerBreaking(t *testing.T) {
	worker := New(func(msg []byte) []byte {
		return nil
	})

	go func(w *Worker) {
		time.Sleep(time.Second)
		w.TerminateExecution()
	}(worker)

	worker.Load("forever.js", ` while (true) { ; } `)
}

func errorContains(t *testing.T, err error, substr string) {
	if err == nil {
		t.Fatal("Expected to get error")
	}
	if !strings.Contains(err.Error(), substr) {
		t.Fatalf("Expected error to have '%s' in it.", substr)
	}
}