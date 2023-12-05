The NATs data source lists a number of NATs resource information owned by an TencentCloud account.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_nat_gateways.

Example Usage

```hcl
# Query the NAT gateway by ID
data "tencentcloud_nats" "anat" {
  id = "nat-k6ualnp2"
}

# Query the list of normal NAT gateways
data "tencentcloud_nats" "nat_state" {
  state = 0
}

# Multi conditional query NAT gateway list
data "tencentcloud_nats" "multi_nat" {
  name           = "terraform test"
  vpc_id         = "vpc-ezij4ltv"
  max_concurrent = 3000000
  bandwidth      = 500
}
```