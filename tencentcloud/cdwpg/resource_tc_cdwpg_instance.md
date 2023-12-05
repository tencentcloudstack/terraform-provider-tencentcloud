Provides a resource to create a cdwpg instance

Example Usage

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

Import

cdwpg instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_instance.instance instance_id
```