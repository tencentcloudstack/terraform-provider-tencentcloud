---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_security_group"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_security_group"
description: |-
  Provides a resource to create a elasticsearch security_group
---

# tencentcloud_elasticsearch_security_group

Provides a resource to create a elasticsearch security_group

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_security_group" "security_group" {
  instance_id = "es-5wn36he6"
  security_group_ids = [
    "sg-mayqdlt1",
    "sg-po2q8cg7",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance Id.
* `security_group_ids` - (Optional, Set: [`String`]) Security group id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyEsVipSecurityGroup](https://cloud.tencent.com/document/api/845/94985)
* Read: [DescribeInstances](https://cloud.tencent.com/document/api/845/30631)
* Update: [ModifyEsVipSecurityGroup](https://cloud.tencent.com/document/api/845/94985)

## Import

elasticsearch security_group can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_security_group.security_group instance_id
```

