Use this data source to query detailed information of mongodb instance_current_op

Example Usage

```hcl
data "tencentcloud_mongodb_instance_current_op" "instance_current_op" {
  instance_id = "cmgo-b43i3wkj"
  op = "command"
  order_by_type = "desc"
}
```