package jscoreWorker

/*
#cgo LDFLAGS: -framework JavaScriptCore

#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include "bridge.h"
typedef void (*closure)();
*/
import "C"
import "strconv"

type ReceiveMessageCallback func(msg []byte) []byte

func SetFlags(args []string) []string {
	return args
}

func Version() string {
	return "JavaScriptCore"
}

type Worker struct {
	// contains filtered or unexported fields
	cb ReceiveMessageCallback
	ctx *JSContext
	id int
}

func New(cb ReceiveMessageCallback) *Worker {
	worker := new(Worker)
	worker.id = uniqueInt()
	worker.cb = cb
	worker.ctx = NewJSContext()

	global := worker.ctx.GetGlobal()

	rawSendFn := C.JSObjectMakeFunctionWithCallback(worker.ctx.ref, NewJSString("send").ref, C.closure(C.jsSend))
	rawRecvFn := C.JSObjectMakeFunctionWithCallback(worker.ctx.ref, NewJSString("recv").ref, C.closure(C.jsRecv))
	rawPrintFn := C.JSObjectMakeFunctionWithCallback(worker.ctx.ref, NewJSString("print").ref, C.closure(C.jsPrint))

	rawWorker := NewJSObject(worker.ctx.Convert())
	rawWorker.SetProperty("send", NewJSValueFromRef(worker.ctx.Convert(), C.JSValueRef(rawSendFn)))
	rawWorker.SetProperty("recv", NewJSValueFromRef(worker.ctx.Convert(), C.JSValueRef(rawRecvFn)))
	rawWorker.SetProperty("print", NewJSValueFromRef(worker.ctx.Convert(), C.JSValueRef(rawPrintFn)))
	global.SetProperty("rawWorker", rawWorker.Value())

	script := `
const JSCoreWorker = {
	send (buf) {
		const ret = rawWorker.send(` + strconv.Itoa(worker.id) + `, buf)
		return ret ? ret.buffer : ret
	},
	recv (cb) { return rawWorker.recv(` + strconv.Itoa(worker.id) + `, cb) },
	print (...any) { return rawWorker.print(` + strconv.Itoa(worker.id) + `, ...any) }
}
	`
	worker.ctx.EvaluateScript(script, "__jscoreWorker.js")

	GoRecv(worker.ctx.Convert(), func(buf []byte) []byte {
		return worker.cb(buf)
	}, worker.id)

	return worker
}

func (w *Worker) Dispose() {

}

func (w *Worker) Load(scriptName string, code string) error {
	_, err := w.ctx.EvaluateScript(code, scriptName)
	return err
}

func (w *Worker) SendBytes(msg []byte) error {
	return GoSend(w.ctx.Convert(), msg, w.id)
}

func (w *Worker) TerminateExecution() {

}