Use this data source to query detailed information of CAM policies

Example Usage

Query all policies

```hcl
data "tencentcloud_cam_policies" "example" {}
```

Query policies by filter

```hcl
data "tencentcloud_cam_policies" "example" {
  name        = "tf-example"
  policy_id   = "236215899"
  type        = 1
  create_mode = 2
}
```