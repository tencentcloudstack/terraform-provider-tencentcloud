Provides a resource to create a vpc route table entry config

~> **NOTE:** When setting the route item switch, do not use it together with resource `tencentcloud_route_table_entry`.

Example Usage

Enable route item

```hcl
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = false
}
```

Disable route item

```hcl
resource "tencentcloud_route_table_entry_config" "example" {
  route_table_id = "rtb-8425lgjy"
  route_item_id  = "rti-4f6efqwn"
  disabled       = true
}
```

Import

vpc route table entry config can be imported using the id, e.g.

```
terraform import tencentcloud_route_table_entry_config.example rtb-8425lgjy#rti-4f6efqwn
```
