build shared libraries:

```bash
cd golang-filter
go build -o libffi_go_basic_auth.so -buildmode=c-shared .
cd ..
cd lua-filter
go build -v -o libgo_basic_auth.so -buildmode=c-shared .
```

blog post:

http://luajit.io/posts/envoy-async-http-filter-lua-resty-ffi-vs-golang/
