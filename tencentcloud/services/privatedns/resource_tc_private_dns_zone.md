Provide a resource to create a Private Dns Zone.

~> **NOTE:** If you want to unbind all VPCs bound to the current private dns zone, simply clearing the declaration will not take effect; you need to set the `region` and `uniq_vpc_id` in `vpc_set` to an empty string.

Example Usage

Create a basic Private Dns Zone

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy = "Terraform"
  }
}
```

Create a Private Dns Zone domain and bind associated accounts'VPC

```hcl
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  account_vpc_set {
    uin         = "123456789"
    uniq_vpc_id = "vpc-adsebmya"
    region      = "ap-guangzhou"
    vpc_name    = "vpc-name"
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy = "Terraform"
  }
}
```

Import

Private Dns Zone can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.example zone-6xg5xgky
```
