Use this data source to query detailed information of oceanus meta_table

Example Usage

```hcl
data "tencentcloud_oceanus_meta_table" "example" {
	work_space_id = "space-6w8eab6f"
	catalog       = "_dc"
	database      = "_db"
	table         = "tf_table"
}
```