Use this data source to query detailed information of lighthouse reset_instance_blueprint

Example Usage

```hcl
data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = "lhins-123456"
  offset = 0
  limit = 20
}
```