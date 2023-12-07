Use this data source to query detailed information of apigateway upstream

Example Usage

```hcl
data "tencentcloud_api_gateway_upstreams" "example" {
  upstream_id = "upstream-4n5bfklc"
}
```

Filtered Queries

```hcl
data "tencentcloud_api_gateway_upstreams" "example" {
  upstream_id = "upstream-4n5bfklc"

  filters {
    name   = "ServiceId"
    values = "service-hvg0uueg"
  }
}
```