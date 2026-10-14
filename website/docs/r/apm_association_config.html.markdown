---
subcategory: "Application Performance Management(APM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_apm_association_config"
sidebar_current: "docs-tencentcloud-resource-apm_association_config"
description: |-
  Provides a resource to create a APM association config
---

# tencentcloud_apm_association_config

Provides a resource to create a APM association config

## Example Usage

```hcl
resource "tencentcloud_apm_association_config" "example" {
  instance_id  = tencentcloud_apm_instance.example.id
  product_name = "Prometheus"
  status       = 1
  peer_id      = "prom-kx3eqdby"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Business system ID.
* `product_name` - (Required, String, ForceNew) Associated product name. currently only supports Prometheus.
* `status` - (Required, Int) Status of the association relationship: // association status: 1 (enabled), 2 (disabled).
* `peer_id` - (Optional, String) Associated product instance ID.
* `topic` - (Optional, String) Specifies the CKafka message topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyApmAssociation](https://cloud.tencent.com/document/api/1463/125076) - This API is used to modify the relationship between the apm Business System and other products, including deletion.
* Read: [DescribeApmAssociation](https://cloud.tencent.com/document/api/1463/125079) - This API is used to query the relationship between apm Business System and other product.
* Update: [ModifyApmAssociation](https://cloud.tencent.com/document/api/1463/125076) - This API is used to modify the relationship between the apm Business System and other products, including deletion.

## Import

APM association config can be imported using the id, e.g.

```
terraform import tencentcloud_apm_association_config.example apm-jPr5iQL77#Prometheus
```

