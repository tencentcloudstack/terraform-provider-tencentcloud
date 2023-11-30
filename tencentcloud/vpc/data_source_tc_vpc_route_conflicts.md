Use this data source to query detailed information of vpc route_conflicts

Example Usage

```hcl
data "tencentcloud_vpc_route_conflicts" "route_conflicts" {
  route_table_id = "rtb-6xypllqe"
  destination_cidr_blocks = ["172.18.111.0/24"]
}
```