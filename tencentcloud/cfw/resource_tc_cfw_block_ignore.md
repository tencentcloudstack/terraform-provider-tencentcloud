Provides a resource to create a cfw block_ignore

~> **NOTE:** If create domain rule, `RuleType` not support set 2.

Example Usage

If create ip rule

```hcl
resource "tencentcloud_cfw_block_ignore" "example" {
  ip         = "1.1.1.1"
  direction  = 0
  comment    = "remark."
  start_time = "2023-09-01 00:00:00"
  end_time   = "2023-10-01 00:00:00"
  rule_type  = 1
}
```

If create domain rule

```hcl
resource "tencentcloud_cfw_block_ignore" "example" {
  domain     = "domain.com"
  direction  = 1
  comment    = "remark."
  start_time = "2023-09-01 00:00:00"
  end_time   = "2023-10-01 00:00:00"
  rule_type  = 1
}
```

Import

cfw block_ignore_list can be imported using the id, e.g.

If import ip rule

```
terraform import tencentcloud_cfw_block_ignore.example 1.1.1.1##0#1
```

If import domain rule

```
terraform import tencentcloud_cfw_block_ignore.example domain.com##0#1
```