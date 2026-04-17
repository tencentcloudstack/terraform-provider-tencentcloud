# Design: MongoDB SRV Connection Domain Suffix Truncation

## Affected File

`tencentcloud/services/mongodb/resource_tc_mongodb_instance_srv_connection.go`

## Change Location

`resourceTencentCloudMongodbInstanceSrvConnectionRead` — after calling `DescribeSRVConnectionDomain`, before writing to state.

## Truncation Logic

Use `strings.SplitN(*domain, ".", 2)[0]` to extract the prefix before the first dot.

```
"123asasdfasdf.gz.tencentmdb.com"  →  "123asasdfasdf"
"myprefix.ap-shanghai.tencentmdb.com"  →  "myprefix"
```

If `domain` contains no dot (already a plain prefix or empty), the value is used as-is.

## Code Diff (Read handler)

```go
// Before
if domain != nil {
    _ = d.Set("domain", domain)
}

// After
if domain != nil {
    prefix := strings.SplitN(*domain, ".", 2)[0]
    _ = d.Set("domain", prefix)
}
```

Add `"strings"` to the import block.
