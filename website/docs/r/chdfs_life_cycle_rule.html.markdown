---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_life_cycle_rule"
sidebar_current: "docs-tencentcloud-resource-chdfs_life_cycle_rule"
description: |-
  Provides a resource to create a chdfs life_cycle_rule
---

# tencentcloud_chdfs_life_cycle_rule

Provides a resource to create a chdfs life_cycle_rule

## Example Usage

```hcl
resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "f14mpfy5lh4e"

  life_cycle_rule {
    life_cycle_rule_name = "terraform-test"
    path                 = "/test"
    status               = 1

    transitions {
      days = 30
      type = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String) file system id.
* `life_cycle_rule` - (Required, List) life cycle rule.

The `life_cycle_rule` object supports the following:

* `life_cycle_rule_name` - (Optional, String) rule name.
* `path` - (Optional, String) rule op path.
* `status` - (Optional, Int) rule status, 1:open, 2:close.
* `transitions` - (Optional, List) life cycle rule transition list.

The `transitions` object supports the following:

* `days` - (Required, Int) trigger days(n day).
* `type` - (Required, Int) transition type, 1: archive, 2: delete, 3: low rate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

chdfs life_cycle_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_life_cycle_rule.life_cycle_rule file_system_id#life_cycle_rule_id
```

