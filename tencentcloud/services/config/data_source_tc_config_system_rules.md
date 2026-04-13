Use this data source to query detailed information of Config system preset rules.

Example Usage

Query all system preset rules

```hcl
data "tencentcloud_config_system_rules" "example" {}
```

Query system rules by keyword

```hcl
data "tencentcloud_config_system_rules" "example" {
  keyword = "cam"
}
```

Query system rules by risk level

```hcl
data "tencentcloud_config_system_rules" "example" {
  risk_level = 1
}
```
