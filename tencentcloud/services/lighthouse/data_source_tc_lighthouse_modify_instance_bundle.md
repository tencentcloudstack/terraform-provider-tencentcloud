Use this data source to query detailed information of lighthouse modify_instance_bundle

Example Usage

```hcl
data "tencentcloud_lighthouse_modify_instance_bundle" "modify_instance_bundle" {
  instance_id = "lhins-xxxxxx"
  filters {
	name = "bundle-id"
	values = ["bundle_gen_mc_med2_02"]

  }
}
```