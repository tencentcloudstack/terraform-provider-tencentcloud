# tencentcloud_teo_dns_record_v13

Provides a TEO DNS Record v13 resource.

## Example Usage

### Create basic A record

```hcl
resource "tencentcloud_teo_dns_record_v13" "example" {
  zone_id  = "zone-xxxxxx"
  type     = "A"
  content  = "1.2.3.4"
  name     = "www"
}
```

### Create CNAME record with TTL and location

```hcl
resource "tencentcloud_teo_dns_record_v13" "cname_example" {
  zone_id  = "zone-xxxxxx"
  type     = "CNAME"
  content  = "example.com"
  name     = "api"
  ttl      = 600
  location = "Mainland"
}
```

### Create MX record with priority

```hcl
resource "tencentcloud_teo_dns_record_v13" "mx_example" {
  zone_id  = "zone-xxxxxx"
  type     = "MX"
  content  = "mail.example.com"
  name     = "@"
  priority = 10
}
```

### Create AAAA record with weight

```hcl
resource "tencentcloud_teo_dns_record_v13" "aaaa_example" {
  zone_id  = "zone-xxxxxx"
  type     = "AAAA"
  content  = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  name     = "ipv6"
  weight   = 10
}
```

### Create TXT record

```hcl
resource "tencentcloud_teo_dns_record_v13" "txt_example" {
  zone_id  = "zone-xxxxxx"
  type     = "TXT"
  content  = "v=DMARC1; p=none"
  name     = "_dmarc"
}
```

### Create SRV record

```hcl
resource "tencentcloud_teo_dns_record_v13" "srv_example" {
  zone_id  = "zone-xxxxxx"
  type     = "SRV"
  content  = "10 60 5060 sipserver.example.com"
  name     = "_sip._tcp"
}
```

### Create CAA record

```hcl
resource "tencentcloud_teo_dns_record_v13" "caa_example" {
  zone_id  = "zone-xxxxxx"
  type     = "CAA"
  content  = "0 issue \"letsencrypt.org\""
  name     = "@"
}
```

### Create NS record

```hcl
resource "tencentcloud_teo_dns_record_v13" "ns_example" {
  zone_id  = "zone-xxxxxx"
  type     = "NS"
  content  = "ns.example.com"
  name     = "sub"
}
```

### Update record parameters

```hcl
resource "tencentcloud_teo_dns_record_v13" "example" {
  zone_id  = "zone-xxxxxx"
  type     = "A"
  content  = "1.2.3.4"
  name     = "www"
}

# Update the record content
resource "tencentcloud_teo_dns_record_v13" "example" {
  zone_id  = "zone-xxxxxx"
  type     = "A"
  content  = "5.6.7.8"
  name     = "www"
  ttl      = 600
  weight   = 20
}
```

### Configure Timeouts

```hcl
resource "tencentcloud_teo_dns_record_v13" "example" {
  zone_id  = "zone-xxxxxx"
  type     = "A"
  content  = "1.2.3.4"
  name     = "www"

  timeouts {
    create = "10m"
    update = "5m"
    delete = "5m"
  }
}
```

### Import existing DNS record

```bash
terraform import tencentcloud_teo_dns_record_v13.example zone-xxxxxx#record-yyyyyy
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID.
* `name` - (Required) DNS record name. If domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `type` - (Required) DNS record type. Valid values are: `A`, `AAAA`, `CNAME`, `TXT`, `MX`, `NS`, `CAA`, `SRV`.
* `content` - (Required) DNS record content. Fill in the corresponding content according to the `type` value. If domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `location` - (Optional) DNS record resolution route. If not specified, defaults to `DEFAULT`, which means the default resolution route and is effective in all regions. Resolution route configuration is only applicable when `type` is `A`, `AAAA`, or `CNAME`. Resolution route configuration is only applicable to standard version and enterprise edition packages.
* `ttl` - (Optional) Cache time. Users can specify a value range of 60-86400. The smaller the value, the faster the modification records will take effect in all regions. Default value: 300. Unit: seconds.
* `weight` - (Optional) DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, defaults to -1, which means no weight is set. Weight configuration is only applicable when `type` is `A`, `AAAA`, or `CNAME`. Note: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.
* `priority` - (Optional) MX record priority, which takes effect only when `type` is `MX`. The smaller the value, the higher the priority. Users can specify a value range of 0-50. Default value is 0 if not specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID, formatted as `zone_id#record_id`.
* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: `enable` (has taken effect), `disable` (has been disabled).
* `created_on` - Creation time.
* `modified_on` - Modification time.

## Import

DNS record can be imported using the `zone_id#record_id`, e.g.

```bash
terraform import tencentcloud_teo_dns_record_v13.example zone-xxxxxx#record-yyyyyy
```
