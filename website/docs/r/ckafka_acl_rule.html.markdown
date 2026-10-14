---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_acl_rule"
sidebar_current: "docs-tencentcloud-resource-ckafka_acl_rule"
description: |-
  Provides a resource to create a ckafka acl_rule
---

# tencentcloud_ckafka_acl_rule

Provides a resource to create a ckafka acl_rule

## Example Usage

```hcl
resource "tencentcloud_ckafka_acl_rule" "acl_rule" {
  instance_id   = "ckafka-xxx"
  resource_type = "Topic"
  pattern_type  = "PREFIXED"
  rule_name     = "RuleName"
  rule_list {
    operation       = "All"
    permission_type = "Deny"
    host            = "*"
    principal       = "User:*"

  }
  pattern    = "prefix"
  is_applied = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) instance id.
* `pattern_type` - (Required, String, ForceNew) Match type, currently supports prefix matching and preset strategy, enumeration value list{PREFIXED/PRESET}.
* `resource_type` - (Required, String, ForceNew) Acl resource type, currently only supports Topic, enumeration value list{Topic}.
* `rule_list` - (Required, List, ForceNew) List of configured ACL rules.
* `rule_name` - (Required, String, ForceNew) rule name.
* `is_applied` - (Optional, Int) Whether the preset ACL rule is applied to the newly added topic.
* `pattern` - (Optional, String, ForceNew) A value representing the prefix that the prefix matches.

The `rule_list` object supports the following:

* `host` - (Required, String, ForceNew) The default is *, which means that any host can be accessed. Currently, ckafka does not support host and ip network segment.
* `operation` - (Required, String, ForceNew) Acl operation mode, enumeration value (all operations All, read Read, write Write).
* `permission_type` - (Required, String, ForceNew) permission type, (Deny|Allow).
* `principal` - (Required, String, ForceNew) User list, the default is User:, which means that any user can access, and the current user can only be the user included in the user list. The input format needs to be prefixed with [User:]. For example, user A is passed in as User:A.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAclRule](https://www.tencentcloud.com/document/product/597/54585) - This API shows you how to create an ACL rule.
* Read: [DescribeAclRule](https://www.tencentcloud.com/document/product/597/54584) - This API is used to query the ACL rule list.
* Update: [ModifyAclRule](https://www.tencentcloud.com/document/product/597/54583) - This API is used to modify ACL policy, currently only support whether to apply preset rules to newly-added topics.
* Delete: [DeleteAclRule](https://www.tencentcloud.com/document/product/597/74428) - This API is used to delete an ACL rule.

## Import

ckafka acl_rule can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_acl_rule.acl_rule acl_rule_id
```

