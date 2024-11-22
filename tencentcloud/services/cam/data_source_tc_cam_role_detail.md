Use this data source to query detailed information of cam role detail

Example Usage

Query cam role detail by role ID

```hcl
data "tencentcloud_cam_role_detail" "example" {
  role_id = "4611686018441060141"
}
```

Query cam role detail by role name

```hcl
data "tencentcloud_cam_role_detail" "example" {
  role_name = "tf-example"
}
```