Use this data source to query detailed information of vpc sg_snapshot_file_content

Example Usage

```hcl
data "tencentcloud_vpc_sg_snapshot_file_content" "sg_snapshot_file_content" {
  snapshot_policy_id = "sspolicy-ebjofe71"
  snapshot_file_id   = "ssfile-017gepjxpr"
  security_group_id  = "sg-ntrgm89v"
}
```