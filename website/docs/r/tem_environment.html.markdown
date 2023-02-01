---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_environment"
sidebar_current: "docs-tencentcloud-resource-tem_environment"
description: |-
  Provides a resource to create a tem environment
---

# tencentcloud_tem_environment

Provides a resource to create a tem environment

## Example Usage

```hcl
resource "tencentcloud_tem_environment" "environment" {
  environment_name = "demo"
  description      = "demo for test"
  vpc              = "vpc-2hfyray3"
  subnet_ids       = ["subnet-rdkj0agk", "subnet-r1c4pn5m", "subnet-02hcj95c"]
  tags = {
    "created" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `environment_name` - (Required, String) environment name.
* `subnet_ids` - (Required, Set: [`String`]) subnet IDs.
* `vpc` - (Required, String) vpc ID.
* `description` - (Optional, String) environment description.
* `tags` - (Optional, Map) environment tag list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem environment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_environment.environment environment_id
```

