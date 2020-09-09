---
subcategory: "Ckafka"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_acls"
sidebar_current: "docs-tencentcloud-datasource-ckafka_acls"
description: |-
  Use this data source to query detailed acl information of Ckafka
---

# tencentcloud_ckafka_acls

Use this data source to query detailed acl information of Ckafka

## Example Usage

```hcl
data "tencentcloud_ckafka_acls" "foo" {
  instance_id   = "ckafka-f9ife4zz"
  resource_type = "TOPIC"
  resource_name = "topic-tf-test"
  host          = "2"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Id of the ckafka instance.
* `resource_name` - (Required) ACL resource name, which is related to `resource_type`. For example, if `resource_type` is `TOPIC`, this field indicates the topic name; if `resource_type` is `GROUP`, this field indicates the group name.
* `resource_type` - (Required) ACL resource type. Valid values are `UNKNOWN`, `ANY`, `TOPIC`, `GROUP`, `CLUSTER`, `TRANSACTIONAL_ID`. Currently, only `TOPIC` is available, and other fields will be used for future ACLs compatible with open-source Kafka.
* `host` - (Optional) Host substr used for querying.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `acl_list` - A list of ckafka acls. Each element contains the following attributes:
  * `host` - IP address allowed to access.
  * `operation_type` - ACL operation mode.
  * `permission_type` - ACL permission type, valid values are `UNKNOWN`, `ANY`, `DENY`, `ALLOW`, and `ALLOW` by default. Currently, CKafka supports `ALLOW` (equivalent to allow list), and other fields will be used for future ACLs compatible with open-source Kafka.
  * `principal` - User which can access. `*` means that any user can access.
  * `resource_name` - ACL resource name, which is related to `resource_type`.
  * `resource_type` - ACL resource type.


