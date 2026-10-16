---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_instance"
sidebar_current: "docs-tencentcloud-resource-cdwpg_instance"
description: |-
  Provides a resource to create a cdwpg instance
---

# tencentcloud_cdwpg_instance

Provides a resource to create a cdwpg instance

## Example Usage

```hcl
resource "tencentcloud_cdwpg_instance" "instance" {
  instance_name  = "test_cdwpg"
  zone           = "ap-guangzhou-6"
  user_vpc_id    = "vpc-xxxxxx"
  user_subnet_id = "subnet-xxxxxx"
  charge_properties {
    renew_flag  = 0
    time_span   = 1
    time_unit   = "h"
    charge_type = "POSTPAID_BY_HOUR"

  }
  admin_password = "xxxxxx"
  resources {
    spec_name = "S_4_16_H_CN"
    count     = 2
    disk_spec {
      disk_type  = "CLOUD_HSSD"
      disk_size  = 200
      disk_count = 1
    }
    type = "cn"

  }
  resources {
    spec_name = "S_4_16_H_CN"
    count     = 2
    disk_spec {
      disk_type  = "CLOUD_HSSD"
      disk_size  = 20
      disk_count = 10
    }
    type = "dn"

  }
  tags = {
    "tagKey" = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `admin_password` - (Required, String) cluster password.
* `charge_properties` - (Required, List) instance billing mode.
* `instance_name` - (Required, String) instance name.
* `resources` - (Required, List) resource information.
* `user_subnet_id` - (Required, String) subnet.
* `user_vpc_id` - (Required, String) private network.
* `zone` - (Required, String) Availability Zone.
* `product_version` - (Optional, String) Version.
* `tags` - (Optional, Map) Tag description list.

The `charge_properties` object supports the following:

* `renew_flag` - (Required, Int) 0-no automatic renewal,1-automatic renewalNote: This field may return null, indicating that a valid value cannot be obtained.
* `time_span` - (Required, Int) Time RangeNote: This field may return null, indicating that a valid value cannot be obtained.
* `time_unit` - (Required, String) Time Unit,Generally h and mNote: This field may return null, indicating that a valid value cannot be obtained.
* `charge_type` - (Optional, String) Charge type, vaild values: PREPAID, POSTPAID_BY_HOUR.

The `disk_spec` object of `resources` supports the following:

* `disk_count` - (Required, Int) disk count.
* `disk_size` - (Required, Int) disk size.
* `disk_type` - (Required, String) disk type.

The `resources` object supports the following:

* `count` - (Required, Int) resource count.
* `disk_spec` - (Required, List) disk Information.
* `spec_name` - (Required, String) resource name.
* `type` - (Required, String) resource type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateInstanceByApi](https://www.tencentcloud.com/search/?wd=CreateInstanceByApi&type=document&lang=en)
* Create: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read: [DescribeInstanceInfo](https://www.tencentcloud.com/search/?wd=DescribeInstanceInfo&type=document&lang=en)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Update: [ModifyInstance](https://www.tencentcloud.com/search/?wd=ModifyInstance&type=document&lang=en)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Update: [ScaleOutInstance](https://www.tencentcloud.com/search/?wd=ScaleOutInstance&type=document&lang=en)
* Update: [UpgradeInstance](https://www.tencentcloud.com/search/?wd=UpgradeInstance&type=document&lang=en)
* Delete: [DescribeInstanceState](https://www.tencentcloud.com/search/?wd=DescribeInstanceState&type=document&lang=en)
* Delete: [DestroyInstanceByApi](https://www.tencentcloud.com/search/?wd=DestroyInstanceByApi&type=document&lang=en)

## Import

cdwpg instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_instance.instance instance_id
```

