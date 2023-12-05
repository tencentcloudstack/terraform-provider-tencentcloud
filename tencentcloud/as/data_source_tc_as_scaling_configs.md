Use this data source to query scaling configuration information.

Example Usage

```hcl
data "tencentcloud_as_scaling_configs" "as_configs" {
  configuration_id   = "asc-oqio4yyj"
  result_output_file = "my_test_path"
}
```