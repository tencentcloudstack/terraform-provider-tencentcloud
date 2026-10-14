---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-resource-security_group"
description: |-
  Provides a resource to create Security group.
---

# tencentcloud_security_group

Provides a resource to create Security group.

## Example Usage

### Create a basic security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg test"
}
```

### Create a complete security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg test"
  project_id  = 0

  tags = {
    "createdBy" = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the security group to be queried.
* `description` - (Optional, String) Description of the security group.
* `project_id` - (Optional, Int, ForceNew) Project ID of the security group.
* `tags` - (Optional, Map) Tags of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `delete` - (Defaults to `3m`) Used when deleting the resource.
## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityGroup](https://www.tencentcloud.com/document/product/215/15806) - This API is used to create a security group (SecurityGroup).
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeSecurityGroups](https://www.tencentcloud.com/document/product/215/15808) - This API (DescribeSecurityGroups) is used to query security groups.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Update: [ModifySecurityGroupAttribute](https://www.tencentcloud.com/document/product/215/15805) - This API (ModifySecurityGroupAttribute) is used to modify the attributes of a security group (SecurityGroupPolicy).
* Delete: [DeleteSecurityGroup](https://www.tencentcloud.com/document/product/215/15803) - This API (DeleteSecurityGroup) is used to delete security groups (SecurityGroup).
* Delete: [DescribeSecurityGroupAssociationStatistics](https://www.tencentcloud.com/document/product/215/17799) - This API (DescribeSecurityGroupAssociationStatistics) is used to query statistics on the instances associated with a security group.

## Import

Security group can be imported using the id, e.g.

```
terraform import tencentcloud_security_group.example sg-ey3wmiz1
```

