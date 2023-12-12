Use this data source to query the detail information of an existing autoscaling group.

Example Usage

```hcl
data "tencentcloud_as_scaling_groups" "as_scaling_groups" {
  scaling_group_name = "myasgroup"
  configuration_id   = "asc-oqio4yyj"
  result_output_file = "my_test_path"
}
```