#!/usr/bin/env bash

FGA_API_URL="http://localhost:8080"
FGA_STORE_ID=ABCD12345678901234567890XX
FGA_BEARER_TOKEN=foobar

curl -X POST $FGA_API_URL/stores/$FGA_STORE_ID/authorization-models \
  -H "Authorization: Bearer $FGA_BEARER_TOKEN" \
  -H "content-type: application/json" \
  -d @./schema.json

curl -X POST $FGA_API_URL/stores/$FGA_STORE_ID/write \
  -H "Authorization: Bearer $FGA_BEARER_TOKEN" \
  -H "content-type: application/json" \
  -d '{"writes": { "tuple_keys" : [{"user":"tom","relation":"writer","object":"document:firstdoc"},{"user":"fred","relation":"reader","object":"document:firstdoc"}] }}'