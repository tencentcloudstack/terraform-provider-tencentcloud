Provides a resource to create a vpc route_table

Example Usage

```hcl
resource "tencentcloud_route_table_association" "route_table_association" {
  route_table_id = "rtb-5toos5sy"
  subnet_id      = "subnet-2y2omd4k"
}
```

Import

vpc route_table can be imported using the id, e.g.

```
terraform import tencentcloud_route_table_association.route_table_association subnet_id
```