Use this data source to query API gateway domain list.

Example Usage

```hcl
resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-ohxqslqe"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test","/root#release"]
}

data "tencentcloud_api_gateway_customer_domains" "id" {
	service_id = tencentcloud_api_gateway_custom_domain.foo.service_id
}
```