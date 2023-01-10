---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_application_service"
sidebar_current: "docs-tencentcloud-resource-tem_application_service"
description: |-
  Provides a resource to create a tem application_service
---

# tencentcloud_tem_application_service

Provides a resource to create a tem application_service

## Example Usage

```hcl
resource "tencentcloud_tem_application_service" "application_service" {
  environment_id = "en-dpxyydl5"
  application_id = "app-jrl3346j"
  service {
    type         = "CLUSTER"
    service_name = "test0-1"
    port_mapping_item_list {
      port        = 80
      target_port = 80
      protocol    = "tcp"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String) application ID.
* `environment_id` - (Required, String) environment ID.
* `service` - (Optional, List) service detail list.

The `port_mapping_item_list` object supports the following:

* `port` - (Optional, Int) container port.
* `protocol` - (Optional, String) UDP or TCP.
* `target_port` - (Optional, Int) application listen port.

The `service` object supports the following:

* `port_mapping_item_list` - (Optional, List) port mapping item list.
* `service_name` - (Optional, String) application service name.
* `type` - (Optional, String) application service type: EXTERNAL | VPC | CLUSTER.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem application_service can be imported using the environmentId#applicationId#serviceName, e.g.

```
terraform import tencentcloud_tem_application_service.application_service en-dpxyydl5#app-jrl3346j#test0-1
```

