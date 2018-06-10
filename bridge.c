#include <JavaScriptCore/JSObjectRef.h>
#include <JavaScriptCore/JSValueRef.h>

extern JSValueRef jsSend_go(JSContextRef ctx, int channel, JSValueRef value);
extern JSValueRef jsRecv_go(JSContextRef ctx, int channel, JSValueRef listener);
extern JSValueRef jsPrint_go(JSContextRef ctx, int channel, const JSValueRef any[], size_t anyCount);

JSValueRef jsSend(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsSend_go(ctx, (int)JSValueToNumber(ctx, arguments[0], nil), arguments[1]);
}

JSValueRef jsRecv(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsRecv_go(ctx, (int)JSValueToNumber(ctx, arguments[0], nil), arguments[1]);
}

JSValueRef jsPrint(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception) {
	return jsPrint_go(ctx, (int)JSValueToNumber(ctx, arguments[0], nil), &arguments[1], argumentCount - 1);
}