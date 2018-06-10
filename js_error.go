package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include <JavaScriptCore/JSValueRef.h>
*/
import "C"

type JSError struct {
	ctx C.JSContextRef
	ref C.JSValueRef
}

func NewJSError(ctx C.JSContextRef) *JSError {
	err := new(JSError)
	err.ctx = ctx
	return err
}