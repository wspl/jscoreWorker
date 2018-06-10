#include <JavaScriptCore/JSObjectRef.h>
#include <JavaScriptCore/JSValueRef.h>

extern JSValueRef jsSend_go(JSContextRef ctx, JSValueRef value, int channel);
extern JSValueRef jsRecv_go(JSContextRef ctx, JSValueRef listener, int channel);
extern JSValueRef jsPrint_go(JSContextRef ctx, JSValueRef any, int channel);

JSValueRef jsSend(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsSend_go(ctx, arguments[0], (int)JSValueToNumber(ctx, arguments[1], nil));
}

JSValueRef jsRecv(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsRecv_go(ctx, arguments[0], (int)JSValueToNumber(ctx, arguments[1], nil));
}

JSValueRef jsPrint(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsPrint_go(ctx, arguments[0], (int)JSValueToNumber(ctx, arguments[1], nil));
}