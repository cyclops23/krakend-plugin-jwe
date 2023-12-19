#!/usr/bin/env bash
set -e

KRAKEND_VERSION='2.5.0'
LIBC_VERSION='MUSL-1.2.4_(alpine-3.18.4)'

docker stop krakend || true

docker run -v "$PWD:/app" -w /app devopsfaith/krakend:"$KRAKEND_VERSION" check-plugin --libc "$LIBC_VERSION"

docker run -it -v "$PWD:/app" -w /app --platform 'linux/amd64' \
  -e CGO_ENABLED=1 -e CC=aarch64-linux-musl-gcc -e GOARCH=arm64 -e GOHOSTARCH=amd64 -e EXTRA_LDFLAGS='-extldflags=-fuse-ld=bfd -extld=aarch64-linux-musl-gcc' \
  krakend/builder:"$KRAKEND_VERSION" \
  go build -ldflags="${EXTRA_LDFLAGS}" -buildmode=plugin -o krakend-plugin-jwe.so .

docker run --name krakend --rm -d -p "8079:8079" -v "$PWD:/app" -v "$PWD:/plugins" \
  -e CORE_JWE_SYMMETRIC_KEY="$CORE_JWE_SYMMETRIC_KEY" \
  -e CORE_JWE_SHARED_SECRET="$CORE_JWE_SHARED_SECRET" \
  devopsfaith/krakend:"$KRAKEND_VERSION" \
  run -dc /app/krakend-example-config.json

sleep 2

set e
docker logs krakend

curl -i \
  -H "Authorization: Bearer eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwia2lkIjoiN2YwYTE1MjEtODNiZi00NjRjLWIxYTEtYzU5ZDYyZWQwMzIxIn0..gXyPId-TSqQb0joOAyLQ9Q.lRjnegKHDwct16JHUdoQpBa4FzfEIxicb7tzurh6IEpvt1Ivtc9kcabXqU9VUkfJCtiQLjDfM1jV4VAUBdsVkbDS9w6778BizFmO9CvACDq0c99SNmHYYT-3AWVYwz-KXzREvAkejqxTmUvTG1XEvkMQR_GZYrQSVVL7Wg-yT5EYY8z8AJNmnbKCC3VQyJqwUxQdeB_6QG8ggPh0CXgC3Wk_gfI5LME6zMnZSUWWu3ius6j0hJx-5H8nxPr6DRQnUcZMKwv0unzeHKhhoQE9UpM7f137CMCvX9ar8KA_Mb2q_zNM6Gm2Woc3CHwPsRonz38CLswx4H0hNvEjfka_FPY6yqxPTwR38BfNz02HAfJ7O4ECUt49fqOo_iZqHa7I25NdosRm3nmD-FAQwxNAgg.q7CfpvWijRLbEMFSqtlIl_Djt3uh1BYHMKkRD3pW6rE" \
  -H "Content-Type: application/json" \
  http://localhost:8079/self-service/v1/integrations

#docker stop krakend
