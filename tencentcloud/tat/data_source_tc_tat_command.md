Use this data source to query detailed information of tat command

Example Usage

```hcl
data "tencentcloud_tat_command" "command" {
  # command_id = ""
  # command_name = ""
  command_type = "SHELL"
  created_by = "TAT"
}
```