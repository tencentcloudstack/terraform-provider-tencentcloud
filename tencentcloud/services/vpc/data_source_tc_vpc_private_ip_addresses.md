Use this data source to query detailed information of vpc private_ip_addresses

Example Usage

```hcl
data "tencentcloud_vpc_private_ip_addresses" "private_ip_addresses" {
  vpc_id = "vpc-l0dw94uh"
  private_ip_addresses = ["10.0.0.1"]
}

```