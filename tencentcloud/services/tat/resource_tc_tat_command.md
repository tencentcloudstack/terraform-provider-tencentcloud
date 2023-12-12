Provides a resource to create a tat command

Example Usage

```hcl
resource "tencentcloud_tat_command" "command" {
  username          = "root"
  command_name      = "ls"
  content           = "bHM="
  description       = "xxx"
  command_type      = "SHELL"
  working_directory = "/root"
  timeout = 50
  tags {
	key = ""
	value = ""
  }
}

```
Import

tat command can be imported using the id, e.g.
```
$ terraform import tencentcloud_tat_command.command cmd-6fydo27j
```