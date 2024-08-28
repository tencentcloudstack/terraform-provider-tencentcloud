Use this data source to query detailed information of organization services

Example Usage

Query all organization services

```hcl
data "tencentcloud_organization_services" "services" {}
```

Query organization services by filter

```hcl
data "tencentcloud_organization_services" "services" {
  search_key = "KeyWord"
}
```
