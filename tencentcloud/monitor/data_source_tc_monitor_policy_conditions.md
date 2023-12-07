Use this data source to query monitor policy conditions(There is a lot of data and it is recommended to output to a file)

Example Usage

```hcl
data "tencentcloud_monitor_policy_conditions" "monitor_policy_conditions" {
  name               = "Cloud Virtual Machine"
  result_output_file = "./tencentcloud_monitor_policy_conditions.txt"
}
```