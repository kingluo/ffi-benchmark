#!/usr/bin/env bash
set -euo pipefail
set -x

OPENRESTY_VERSION=${OPENRESTY_VERSION:-1.21.4.2}
OR_PREFIX=${OR_PREFIX:="/usr/local/openresty"}
debug_args="--with-debug"

wget --no-check-certificate https://openresty.org/download/openresty-${OPENRESTY_VERSION}.tar.gz
tar zxf openresty-${OPENRESTY_VERSION}.tar.gz

luajit_xcflags=${luajit_xcflags:="-DLUAJIT_NUMMODE=2 -DLUAJIT_ENABLE_LUA52COMPAT"}
no_pool_patch=${no_pool_patch:-}

cd openresty-${OPENRESTY_VERSION}

(
cd bundle/
for f in /opt/lua-resty-ffi/patches/*; do
    patch -p0 < $f
done
cp -a /opt/lua-resty-ffi/ngx_http_lua_ffi.c ngx_lua-0.10.25/src/
cp -a /opt/lua-resty-ffi/resty_ffi.lua lua-resty-core-0.1.27/lib/resty/core/
)

./configure --prefix="$OR_PREFIX" \
    $debug_args \
    --with-poll_module \
    --with-pcre-jit \
    --without-http_rds_json_module \
    --without-http_rds_csv_module \
    --without-lua_rds_parser \
    --with-stream \
    --with-stream_ssl_module \
    --with-stream_ssl_preread_module \
    --with-http_v2_module \
    --without-mail_pop3_module \
    --without-mail_imap_module \
    --without-mail_smtp_module \
    --with-http_stub_status_module \
    --with-http_realip_module \
    --with-http_addition_module \
    --with-http_auth_request_module \
    --with-http_secure_link_module \
    --with-http_random_index_module \
    --with-http_gzip_static_module \
    --with-http_sub_module \
    --with-http_dav_module \
    --with-http_flv_module \
    --with-http_mp4_module \
    --with-http_gunzip_module \
    --with-threads \
    --with-compat \
    --with-luajit-xcflags="$luajit_xcflags" \
    $no_pool_patch \
    -j`nproc`

make -j`nproc`
sudo make install
