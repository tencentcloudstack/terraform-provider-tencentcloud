Provide a resource to create a Private Dns Record.

Example Usage

```hcl
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

# create private dns zone
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

# create private dns record
resource "tencentcloud_private_dns_record" "example" {
  zone_id      = tencentcloud_private_dns_zone.example.id
  record_type  = "A"
  record_value = "192.168.1.2"
  sub_domain   = "www"
  ttl          = 300
  weight       = 20
  mx           = 0
  status       = "disabled"
}
```

Import

Private Dns Record can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_record.example zone-iza3a33s#1983030
```