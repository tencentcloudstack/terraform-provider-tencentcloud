Use this data source to query detailed information of vpc limits

Example Usage

```hcl
data "tencentcloud_vpc_limits" "limits" {
  limit_types = ["appid-max-vpcs", "vpc-max-subnets"]
}
```