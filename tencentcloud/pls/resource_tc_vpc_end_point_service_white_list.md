Provides a resource to create a vpc end_point_service_white_list

Example Usage

```hcl
resource "tencentcloud_vpc_end_point_service_white_list" "end_point_service_white_list" {
  user_uin = "100020512675"
  end_point_service_id = "vpcsvc-69y13tdb"
  description = "terraform for test"
}
```

Import

vpc end_point_service_white_list can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list end_point_service_white_list_id
```