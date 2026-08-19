Use this data source to query the key alias list specified with region supported by the audit.

Example Usage
```hcl
data "tencentcloud_audit_key_alias" "all" {
	region = "ap-hongkong"
}
```