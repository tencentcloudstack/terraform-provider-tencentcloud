Use this data source to query detailed information of DNATs.

Example Usage

```hcl
# query by nat gateway id
data "tencentcloud_dnats" "foo" {
  nat_id = "nat-xfaq1"
}

# query by vpc id
data "tencentcloud_dnats" "foo" {
  vpc_id = "vpc-xfqag"
}

# query by elastic ip
data "tencentcloud_dnats" "foo" {
  elastic_ip = "123.207.115.136"
}
```