Provides a resource to create a vpc ccn_routes switch

Example Usage

```hcl
resource "tencentcloud_ccn_routes" "example" {
  ccn_id   = "ccn-gr7nynbd"
  route_id = "ccnrtb-jpf7bzn3"
  switch   = "off"
}
```

Import

vpc ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_ccn_routes.ccn_routes ccn-gr7nynbd#ccnr-5uhewx1s
```
