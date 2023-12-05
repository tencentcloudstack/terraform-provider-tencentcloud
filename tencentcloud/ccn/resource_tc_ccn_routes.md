Provides a resource to create a vpc ccn_routes

Example Usage

```hcl
resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = "ccn-39lqkygf"
  route_id = "ccnr-3o0dfyuw"
  switch = "on"
}
```

Import

vpc ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_ccn_routes.ccn_routes ccnId#routesId
```