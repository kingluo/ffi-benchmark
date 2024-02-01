package main

/*
#cgo LDFLAGS: -shared
#include <string.h>
void* ngx_http_lua_ffi_task_poll(void *p);
char* ngx_http_lua_ffi_get_req(void *tsk, int *len);
void ngx_http_lua_ffi_respond(void *tsk, int rc, char* rsp, int rsp_len);
*/
import "C"
import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"unsafe"
)

type Request struct {
	Auth string `json:"auth"`
	Body string `json:"body"`
}

type Cfg struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var cfg Cfg

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	username, password, ok = strings.Cut(cs, ":")
	if !ok {
		return "", "", false
	}
	return username, password, true
}

func verify(auth string) (bool, string) {
	if auth == "" {
		return false, "no Authorization"
	}
	username, password, ok := parseBasicAuth(auth)
	if !ok {
		return false, "invalid Authorization format"
	}
	if username == cfg.User && password == cfg.Password {
		return true, ""
	}
	return false, "invalid username or password"
}

//export libffi_init
func libffi_init(cstr *C.char, tq unsafe.Pointer) C.int {
	err := json.Unmarshal([]byte(C.GoString(cstr)), &cfg)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			task := C.ngx_http_lua_ffi_task_poll(tq)
			if task == nil {
				break
			}
			var rlen C.int
			r := C.ngx_http_lua_ffi_get_req(task, &rlen)
			c := unsafe.Pointer(r)
			var i int
			for ; i < int(rlen); i++ {
				b := *(*byte)(unsafe.Add(c, uintptr(i)*1))
				if b == '\n' {
					break
				}
			}
			str1 := unsafe.String((*byte)(unsafe.Pointer(r)), i)
			str2 := unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(r), uintptr(i+1)*1)), int(rlen)-i-1)
			req := &Request{Auth: str1, Body: str2}
			go func() {
				if ok, msg := verify(req.Auth); !ok {
					C.ngx_http_lua_ffi_respond(task, 401, (*C.char)(C.CString(msg)), C.int(len(msg)))
				} else {
					C.ngx_http_lua_ffi_respond(task, 200, (*C.char)(C.CString(req.Body)), C.int(len(req.Body)))
				}
			}()
		}
	}()
	return 0
}

func main() {}
