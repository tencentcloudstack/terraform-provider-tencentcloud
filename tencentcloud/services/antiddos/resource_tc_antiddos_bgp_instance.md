Provides a resource to create a AntiDDoS bgp instance

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud AntiDDoS through a ticket.

Example Usage

Create standard bgp instance(POSTPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "POSTPAID_BY_MONTH"
  package_type         = "Standard"
  standard_package_config {
    region                 = "ap-guangzhou"
    protect_ip_count       = 1
    bandwidth              = 100
    elastic_bandwidth_flag = true
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Create standard edition 2.0 bgp instance(PREPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "PREPAID"
  package_type         = "StandardPlus"
  instance_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }

  standard_plus_package_config {
    region                 = "ap-guangzhou"
    protect_count          = "TWO_TIMES"
    protect_ip_count       = 1
    bandwidth              = 100
    elastic_bandwidth_flag = true
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Create enterprise bgp instance(POSTPAID)

```hcl
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "POSTPAID_BY_MONTH"
  package_type         = "Enterprise"

  enterprise_package_config {
    region                  = "ap-guangzhou"
    protect_ip_count        = 10
    basic_protect_bandwidth = 300
    bandwidth               = 100
    elastic_bandwidth_flag  = false
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Import

AntiDDoS bgp instance can be imported using the resourceId#packageRegion, e.g.

```
terraform import tencentcloud_antiddos_bgp_instance.example bgp-00000fyi#ap-guangzhou
```
