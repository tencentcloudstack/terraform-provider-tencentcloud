Use this data source to query detailed information of vpc security_group_references

Example Usage

```hcl
data "tencentcloud_vpc_security_group_references" "security_group_references" {
  security_group_ids = ["sg-edmur627"]
}
```