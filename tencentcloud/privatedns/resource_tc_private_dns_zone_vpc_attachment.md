Provides a resource to create a PrivateDns zone_vpc_attachment

~> **NOTE:**  If you need to bind account A to account B's VPC resources, you need to first grant role authorization to account A.

Example Usage

Append VPC associated with private dns zone

```hcl
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone_vpc_attachment" "example" {
  zone_id = tencentcloud_private_dns_zone.example.id

  vpc_set {
    uniq_vpc_id = tencentcloud_vpc.vpc.id
    region      = "ap-guangzhou"
  }
}
```

Add VPC information for associated accounts in the private dns zone

```hcl
resource "tencentcloud_private_dns_zone_vpc_attachment" "example" {
  zone_id = tencentcloud_private_dns_zone.example.id

  account_vpc_set {
    uniq_vpc_id = "vpc-82znjzn3"
    region      = "ap-guangzhou"
    uin         = "100017155920"
  }
}
```

Import

PrivateDns zone_vpc_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_zone_vpc_attachment.example zone-6t11lof0#vpc-jdx11z0t
```