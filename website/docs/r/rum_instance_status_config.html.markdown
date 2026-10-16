---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_instance_status_config"
sidebar_current: "docs-tencentcloud-resource-rum_instance_status_config"
description: |-
  Provides a resource to create a rum instance_status_config
---

# tencentcloud_rum_instance_status_config

Provides a resource to create a rum instance_status_config

## Example Usage

```hcl
resource "tencentcloud_rum_instance_status_config" "instance_status_config" {
  instance_id = "rum-pasZKEI3RLgakj"
  operate     = "stop"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `operate` - (Required, String) `resume`, `stop`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_status` - Instance status (`1`=creating, `2`=running, `3`=abnormal, `4`=restarting, `5`=stopping, `6`=stopped, `7`=deleted).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ResumeInstance](https://www.tencentcloud.com/search/?wd=ResumeInstance&type=document&lang=en)
* Create: [StopInstance](https://www.tencentcloud.com/search/?wd=StopInstance&type=document&lang=en)
* Read: [DescribeTawInstances](https://www.tencentcloud.com/search/?wd=DescribeTawInstances&type=document&lang=en)
* Update: [ResumeInstance](https://www.tencentcloud.com/search/?wd=ResumeInstance&type=document&lang=en)
* Update: [StopInstance](https://www.tencentcloud.com/search/?wd=StopInstance&type=document&lang=en)

## Import

rum instance_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_config.instance_status_config instance_id
```

