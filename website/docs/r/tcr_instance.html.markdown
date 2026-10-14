---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_instance"
sidebar_current: "docs-tencentcloud-resource-tcr_instance"
description: |-
  Use this resource to create tcr instance.
---

# tencentcloud_tcr_instance

Use this resource to create tcr instance.

~> **NOTE:**If `security_policy` needs to be configured, `open_public_operation` needs to be set to true

## Example Usage

### Create a basic tcr instance.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"

  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create instance with the public network access whitelist.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name                  = "tf-example-tcr"
  instance_type         = "basic"
  open_public_operation = true
  security_policy {
    cidr_block = "10.0.0.1/24"
  }

  security_policy {
    cidr_block = "192.168.1.1"
  }
}
```

### Create instance with Replications.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  replications {
    region_id = var.tcr_region_map["ap-guangzhou"] # 1
  }

  replications {
    region_id = var.tcr_region_map["ap-singapore"] # 9
  }
}

variable "tcr_region_map" {
  default = {
    "ap-guangzhou"     = 1
    "ap-shanghai"      = 4
    "ap-hongkong"      = 5
    "ap-beijing"       = 8
    "ap-singapore"     = 9
    "na-siliconvalley" = 15
    "ap-chengdu"       = 16
    "eu-frankfurt"     = 17
    "ap-seoul"         = 18
    "ap-chongqing"     = 19
    "ap-mumbai"        = 21
    "na-ashburn"       = 22
    "ap-bangkok"       = 23
    "eu-moscow"        = 24
    "ap-tokyo"         = 25
    "ap-nanjing"       = 33
    "ap-taipei"        = 39
    "ap-jakarta"       = 72
  }
}
```

### Create instance with COS bucket configuration.

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name                  = "tf-example-tcr"
  instance_type         = "standard"
  enable_cos_maz        = true
  enable_cos_versioning = true

  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String) TCR types. Valid values are: `standard`, `basic`, `premium`.
* `name` - (Required, String, ForceNew) Name of the TCR instance.
* `delete_bucket` - (Optional, Bool) Indicate to delete the COS bucket which is auto-created with the instance or not.
* `deletion_protection` - (Optional, Bool) Whether to enable Instance Deletion Protection.
* `enable_cos_maz` - (Optional, Bool, ForceNew) Whether to enable COS bucket multi-AZ feature. Default is `false`.
* `enable_cos_versioning` - (Optional, Bool, ForceNew) Whether to enable COS bucket versioning. Advanced Edition Instances: Default is `true` (versioning enabled); Standard / Basic Edition Instances: Default is `false` (disabled).
* `instance_charge_type_prepaid_period` - (Optional, Int) Length of time to purchase an instance (in month). Must set when registry_charge_type is prepaid.
* `instance_charge_type_prepaid_renew_flag` - (Optional, Int) Auto renewal flag. 1: manual renewal, 2: automatic renewal, 3: no renewal and no notification. Must set when registry_charge_type is prepaid.
* `open_public_operation` - (Optional, Bool) Control public network access.
* `registry_charge_type` - (Optional, Int) Charge type of instance. 1: postpaid; 2: prepaid. Default is postpaid.
* `replications` - (Optional, List) Specify List of instance Replications, premium only. The available [source region list](https://www.tencentcloud.com/document/api/1051/41101) is here.
* `security_policy` - (Optional, Set) Public network access allowlist policies of the TCR instance. Only available when `open_public_operation` is `true`.
* `tags` - (Optional, Map) The available tags within this TCR instance.

The `replications` object supports the following:

* `region_id` - (Optional, Int) Replication region ID, check the example at the top of page to find out id of region.
* `region_name` - (Optional, String) Replication region name.
* `syn_tag` - (Optional, Bool) Specify whether to sync TCR cloud tags to COS Bucket. NOTE: You have to specify when adding, modifying will be ignored for now.

The `security_policy` object supports the following:

* `cidr_block` - (Optional, String) The public network IP address of the access source.
* `description` - (Optional, String) Remarks of policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `expired_at` - Instance expiration time (prepaid).
* `internal_end_point` - Internal address for access of the TCR instance.
* `public_domain` - Public address for access of the TCR instance.
* `public_status` - Status of the TCR instance public network access.
* `status` - Status of the TCR instance.

The `replications` object exports the following:

* `id` - Replication registry ID (readonly).

The `security_policy` object exports the following:

* `index` - Index of policy.
* `version` - Version of policy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateInstance](https://cloud.tencent.com/document/api/1141/41572) - This API is used to create an instance.
* Create: [CreateMultipleSecurityPolicy](https://cloud.tencent.com/document/api/1141/56524) - This API is used to create multiple public network access allowlist policies of the TCR instance.
* Create: [DescribeExternalEndpointStatus](https://cloud.tencent.com/document/api/1141/53642) - This API is used to query the public network access entry status of an instance.
* Create: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Create: [ManageExternalEndpoint](https://cloud.tencent.com/document/api/1141/53597) - This API is used to manage the public network access of an instance.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeExternalEndpointStatus](https://cloud.tencent.com/document/api/1141/53642) - This API is used to query the public network access entry status of an instance.
* Read: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Read: [DescribeReplicationInstances](https://cloud.tencent.com/document/api/1141/49201) - This API is used to query the list of replication instances.
* Read: [DescribeSecurityPolicies](https://cloud.tencent.com/document/api/1141/53906) - This API is used to query the public network access allowlist policies of an instance.
* Update: [CreateMultipleSecurityPolicy](https://cloud.tencent.com/document/api/1141/56524) - This API is used to create multiple public network access allowlist policies of the TCR instance.
* Update: [DeleteMultipleSecurityPolicy](https://cloud.tencent.com/document/api/1141/56523) - This API is used to delete multiple public network access allowlist policies of the instance.
* Update: [DescribeExternalEndpointStatus](https://cloud.tencent.com/document/api/1141/53642) - This API is used to query the public network access entry status of an instance.
* Update: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Update: [ManageExternalEndpoint](https://cloud.tencent.com/document/api/1141/53597) - This API is used to manage the public network access of an instance.
* Update: [ModifyInstance](https://cloud.tencent.com/document/api/1141/60996) - This API is used to update instance information.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [RenewInstance](https://cloud.tencent.com/document/api/1141/51792) - This API is used to renew a prepaid instance or change the billing mode from pay-as-you-go billing to monthly subscription billing.
* Delete: [DeleteInstance](https://cloud.tencent.com/document/api/1141/49983) - This API is used to delete a TCR Enterprise Edition instance.
* Delete: [DeleteNamespace](https://cloud.tencent.com/document/api/1141/42728) - This API is used to delete a namespace.
* Delete: [DeleteReplicationInstance](https://cloud.tencent.com/document/api/1141/79581) - This API is used to delete a replica instance.
* Delete: [DescribeInstances](https://cloud.tencent.com/document/api/1141/41569) - This API is used to query the instance information.
* Delete: [DescribeNamespaces](https://cloud.tencent.com/document/api/1141/42765) - This API is used to query the namespace list or the information of the specified namespace.
* Delete: [DescribeReplicationInstances](https://cloud.tencent.com/document/api/1141/49201) - This API is used to query the list of replication instances.

## Import

tcr instance can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_instance.example tcr-4detlt3v
```

