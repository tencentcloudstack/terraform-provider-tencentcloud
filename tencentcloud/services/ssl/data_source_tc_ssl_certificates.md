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
}

data "tencentcloud_ssl_certificates" "example" {
  type = "CA"
}

data "tencentcloud_ssl_certificates" "example" {
  id = "LCYouprI"
}

data "tencentcloud_ssl_certificates" "example" {
  domain = "example.com"
}
```

Query the latest SSL certificate of a domain

```hcl
data "tencentcloud_ssl_certificates" "example" {
  domain  = "example.com"
  sort_by = "end_time"
}
```
