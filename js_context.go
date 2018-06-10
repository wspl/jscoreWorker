package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
*/
import "C"

type JSContext struct {
	ref C.JSGlobalContextRef
}

func NewJSContext() *JSContext {
	ctx := new(JSContext)
	ctxGroup := C.JSContextGroupCreate()
	ctxRef := C.JSGlobalContextCreateInGroup(ctxGroup, nil)
	ctx.ref = ctxRef
	return ctx
}

func NewJSContextFromRef(ref C.JSGlobalContextRef) *JSContext {
	ctx := new(JSContext)
	ctx.ref = ref
	return ctx
}

func (ctx *JSContext) Convert() C.JSContextRef {
	return C.JSContextRef(ctx.ref)
}

func (ctx *JSContext) GetGlobal() *JSObject {
	return NewJSObjectFromRef(ctx.Convert(), C.JSContextGetGlobalObject(ctx.ref))
}

func (ctx *JSContext) EvaluateScript(script string, sourceUrl string) (*JSValue, error) {
	jsErr := NewJSError(ctx.Convert())
	ret := C.JSEvaluateScript(
		ctx.ref,
		NewJSString(script).ref,
		NewJSObject(ctx.Convert()).ref,
		NewJSString(sourceUrl).ref,
		C.int(0),
		&jsErr.ref)

	if jsErr.ref != nil {
		return nil, jsErr.Error(script, sourceUrl)
	}

	return NewJSValueFromRef(ctx.Convert(), ret), nil
}