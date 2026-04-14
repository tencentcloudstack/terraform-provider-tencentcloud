Use this data source to query detailed information of Config compliance packs.

Example Usage

Query all compliance packs

```hcl
data "tencentcloud_config_compliance_packs" "example" {}
```

Query compliance packs by name

```hcl
data "tencentcloud_config_compliance_packs" "example" {
  compliance_pack_name = "tf-example"
}
```

Query compliance packs by filters

```hcl
data "tencentcloud_config_compliance_packs" "example" {
  compliance_pack_name = "tf-example"
  risk_level           = [1, 2]
  status               = "ACTIVE"
  compliance_result    = ["NON_COMPLIANT"]
  order_type           = "desc"
}
```
