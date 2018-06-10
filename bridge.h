#include <JavaScriptCore/JSContextRef.h>

JSValueRef jsSend(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception);
JSValueRef jsRecv(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception);
JSValueRef jsPrint(JSContextRef ctx, JSObjectRef function, JSObjectRef thisObject, size_t argumentCount, const JSValueRef arguments[], JSValueRef *exception);

// for terminate js
typedef bool (*JSShouldTerminateCallback) (JSContextRef ctx, void* context);
void JSContextGroupSetExecutionTimeLimit(JSContextGroupRef group, double limit, JSShouldTerminateCallback callback, void* context);