Use this data source to query SSL certificates.

Example Usage

Query all SSL certificates

```hcl
data "tencentcloud_ssl_certificates" "example" {}
```

Query SSL certificates by filter

```hcl
data "tencentcloud_ssl_certificates" "example" {
  name = "tf-example"
  type = "CA"
  id   = "LCYouprI"
}
```