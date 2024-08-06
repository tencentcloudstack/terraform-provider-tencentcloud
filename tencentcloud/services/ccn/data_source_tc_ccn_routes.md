Use this data source to query detailed information of CCN routes.

Example Usage

Query CCN instance all routes

```hcl
data "tencentcloud_ccn_routes" "routes" {
  ccn_id = "ccn-gr7nynbd"
}
```

Query CCN instance routes by filter

```hcl
data "tencentcloud_ccn_routes" "routes" {
  ccn_id = "ccn-gr7nynbd"
  filters {
    name   = "route-table-id"
    values = ["ccnrtb-jpf7bzn3"]
  }
}
```
