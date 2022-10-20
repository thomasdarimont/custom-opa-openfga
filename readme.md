Custom Open Policy Agent with prototypical support for OpenFGA
---

This experiment adds support for querying relations from [OpenFGA](https://openfga.dev/) via GRPC to check resource level permissions
as custom builtin commands for [Open Policy Agent](https://www.openpolicyagent.org/).

Currently only one command is supported:
```
openfga.check_permission("SUBJECT","PERMISSION","RESOURCE_ID") -> bool
```

# Build

Note this example uses Go 1.19

```
go get
go build
```

# Demo

> Start openfga demo environment
```
docker compose -f demo/docker-compose.yml up -d --remove-orphans
```

> Run custom Open Policy Agent with openfga plugin enabled
```
./custom-opa-openfga run \
  --set plugins.openfga.apiScheme=http \
  --set plugins.openfga.apiHost=localhost:8080 \
  --set plugins.openfga.apiToken=foobar \
  --set plugins.openfga.storeId=ABCD12345678901234567890XX
```

> Create store and import the example data
```
FGA_API_URL="http://localhost:8080"
FGA_STORE_ID=ABCD12345678901234567890XX
FGA_BEARER_TOKEN=foobar

curl -X POST $FGA_API_URL/stores/$FGA_STORE_ID/authorization-models \
  -H "Authorization: Bearer $FGA_BEARER_TOKEN" \
  -H "content-type: application/json" \
  -d @./demo/schema.json

curl -X POST $FGA_API_URL/stores/$FGA_STORE_ID/write \
  -H "Authorization: Bearer $FGA_BEARER_TOKEN" \
  -H "content-type: application/json" \
  -d '{"writes": { "tuple_keys" : [{"user":"tom","relation":"writer","object":"document:firstdoc"},{"user":"fred","relation":"reader","object":"document:firstdoc"}] }}'
```

> Query relations against openfga
> See the [example schema](./demo/schema.json) for reference.
```
> openfga.check_permission("tom", "view", "document:firstdoc")
true
> openfga.check_permission("tom", "edit", "document:firstdoc")
true
> openfga.check_permission("fred", "edit", "document:firstdoc")
false
> exit
```

> Stop demo environment
```
docker compose -f demo/docker-compose.yml down
```
