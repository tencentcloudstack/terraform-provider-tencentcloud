Provides a resource to manage TEO (TencentCloud EdgeOne) DNS record v2.

~> **NOTE:** This resource uses TEO DNS record API v2, which provides enhanced features compared to the legacy version.

~> **NOTE:** Weight parameter is only applicable for A, AAAA, and CNAME record types.

~> **NOTE:** Priority parameter is only applicable for MX record type.

~> **NOTE:** Location parameter is only applicable for A, AAAA, and CNAME record types.

~> **NOTE:** TTL must be between 60 and 86400 seconds. Default is 300 seconds.

~> **NOTE:** Weight must be between -1 and 100. Default is -1 (no weight).

~> **NOTE:** Priority must be between 0 and 50. Default is 0 (highest priority for MX records).

## Example Usage

### Create an A record

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_a" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "www"
  type    = "A"
  content = "1.2.3.4"
  ttl     = 600
}
```

### Create a CNAME record with weight

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_cname" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "api"
  type    = "CNAME"
  content = "backend.example.com"
  ttl     = 300
  weight  = 10
}
```

### Create an MX record with priority

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_mx" {
  zone_id  = "zone-xxxxxxxxxx"
  name     = "@"
  type     = "MX"
  content  = "mail.example.com"
  ttl      = 600
  priority = 10
}
```

### Create a TXT record

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_txt" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "_dmarc"
  type    = "TXT"
  content = "v=DMARC1; p=none; rua=mailto:dmarc@example.com"
  ttl     = 600
}
```

### Create a record with location (line)

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_location" {
  zone_id  = "zone-xxxxxxxxxx"
  name     = "cdn"
  type     = "CNAME"
  content  = "cdn.example.com"
  location = "电信"  # Telecom line
  ttl      = 600
  weight   = 20
}
```

### Create an AAAA record

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_aaaa" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "ipv6"
  type    = "AAAA"
  content = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  ttl     = 600
}
```

### Update a DNS record

```hcl
resource "tencentcloud_teo_dns_record_v2" "example_update" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "www"
  type    = "A"
  content = "1.2.3.4"
  ttl     = 600
}

# Update the content and ttl
resource "tencentcloud_teo_dns_record_v2" "example_update" {
  zone_id = "zone-xxxxxxxxxx"
  name    = "www"
  type    = "A"
  content = "5.6.7.8"  # Changed
  ttl     = 300        # Changed
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `name` - (Required) DNS record name.
* `type` - (Required) DNS record type. Valid values: `A`, `AAAA`, `CNAME`, `MX`, `TXT`, `NS`, `CAA`, `SRV`.
* `content` - (Required) DNS record content. The format varies by record type.
* `location` - (Optional) DNS record resolution line. Default: `Default`. Only applicable when Type is A, AAAA, or CNAME.
* `ttl` - (Optional) Cache time in seconds. Range: 60-86400. Default: 300.
* `weight` - (Optional) DNS record weight. Range: -1~100. Default: -1. Only applicable when Type is A, AAAA, or CNAME.
* `priority` - (Optional) MX record priority. Range: 0~50. Default: 0. Only applicable when Type is MX.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID, formatted as `zone_id#record_id`.
* `record_id` - DNS record ID.
* `status` - DNS record status. Valid values: `enable`, `disable`.
* `created_on` - Creation time.
* `modified_on` - Modification time.

## Import

TEO DNS record v2 can be imported using the id, e.g.

```
terraform import tencentcloud_teo_dns_record_v2.example zone-xxxxxxxxxx#record-yyyyyyyyyy
```

The ID format is: `zone_id#record_id`

## Timeouts

The `timeouts` block allows you to specify timeouts for certain operations:

* `create` - (Default 5 minutes) Used for creating DNS record.
* `read` - (Default 3 minutes) Used for reading DNS record.
* `update` - (Default 5 minutes) Used for updating DNS record.
* `delete` - (Default 5 minutes) Used for deleting DNS record.
