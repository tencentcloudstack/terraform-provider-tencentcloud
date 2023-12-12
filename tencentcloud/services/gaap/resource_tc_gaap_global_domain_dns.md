Provides a resource to create a gaap global domain dns

Example Usage

```hcl
resource "tencentcloud_gaap_global_domain_dns" "global_domain_dns" {
	domain_id = "dm-xxxxxx"
	proxy_id_list = ["link-xxxxxx"]
	nation_country_inner_codes = ["101001"]
}
```

Import

gaap global_domain_dns can be imported using the id, e.g.

```
terraform import tencentcloud_gaap_global_domain_dns.global_domain_dns ${domainId}#${dnsRecordId}
```