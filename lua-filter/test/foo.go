package test

import "C"
import (
	"encoding/json"
	"testing"
	"unsafe"
)

type Request struct {
	Auth string `json:"auth"`
	Body string `json:"body"`
}

var rawStr *C.char = C.CString("foobar\nfoobar")
var jsonStr *C.char = C.CString(`{"auth":"foobar","body":"foobar"}`)

func testRAW(t *testing.B) {
	rlen := C.int(13)
	c := unsafe.Pointer(rawStr)
	var i int
	for ; i < int(rlen); i++ {
		b := *(*byte)(unsafe.Add(c, uintptr(i)*1))
		if b == '\n' {
			break
		}
	}
	str1 := unsafe.String((*byte)(unsafe.Pointer(rawStr)), i)
	str2 := unsafe.String((*byte)(unsafe.Add(unsafe.Pointer(rawStr), uintptr(i+1)*1)), int(rlen)-i-1)
	req := &Request{Auth: str1, Body: str2}
	if req.Auth != "foobar" || req.Body != "foobar" {
		t.Errorf("failed %+v", req)
	}
}

func testJSON(t *testing.B) {
	data := C.GoBytes(unsafe.Pointer(jsonStr), 33)
	var req Request
	err := json.Unmarshal(data, &req)
	if err != nil {
		panic(err)
	}
	if req.Auth != "foobar" || req.Body != "foobar" {
		t.Errorf("failed")
	}
}
