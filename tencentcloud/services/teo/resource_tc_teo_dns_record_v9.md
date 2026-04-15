# tencentcloud_teo_dns_record_v9

Provides a resource to manage a DNS record for a teo (EdgeOne) zone.

## Example Usage

### Create A record

```hcl
resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
  ttl      = 300
}
```

### Create CNAME record with location and weight

```hcl
resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id  = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "CNAME"
  content  = "alias.example.com"
  location = "Asia"
  weight   = 10
  ttl      = 600
}
```

### Create MX record with priority

```hcl
resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id  = tencentcloud_teo_zone.test.id
  name     = "@"
  type     = "MX"
  content  = "mail.example.com"
  priority = 10
}
```

### Update record fields

```hcl
resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
  ttl      = 300

  # Later update
  content  = "5.6.7.8"
  ttl      = 600
}
```

### Delete record

```hcl
# The DNS record will be deleted when it's removed from the configuration

resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
}

# To delete the record, simply remove it from the configuration
```

### Timeout configuration

```hcl
resource "tencentcloud_teo_zone" "test" {
  zone_name = "example.com"
  area      = "mainland"
  type      = "full"
}

resource "tencentcloud_teo_dns_record_v9" "test" {
  zone_id = tencentcloud_teo_zone.test.id
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
  ttl      = 300

  timeouts {
    create = 10m
    update = 10m
    delete = 10m
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `name` - (Required) DNS record name. For Chinese, Korean, or Japanese domain names, convert to punycode before input.
* `type` - (Required) DNS record type. Valid values: `A`, `AAAA`, `MX`, `CNAME`, `TXT`, `NS`, `CAA`, `SRV`.
* `content` - (Required) DNS record content. Fill in content corresponding to the Type value. For Chinese, Korean, or Japanese domain names, convert to punycode before input.
* `location` - (Optional) DNS record resolution line. Default is `Default`, meaning default resolution line and effective for all regions. Only applicable when Type is `A`, `AAAA`, or `CNAME`.
* `ttl` - (Optional) Cache time, range 60-86400 seconds. Smaller value means faster propagation of changes. Default is 300.
* `weight` - (Optional) DNS record weight, range -1 to 100. -1 means no weight set, 0 means no resolution. Only applicable when Type is `A`, `AAAA`, or `CNAME`.
* `priority` - (Optional) MX record priority, range 0-50. Lower value means higher priority. Only applicable when Type is `MX`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID in the format `zone_id#record_id`.
* `status` - DNS record resolution status. Valid values: `enable` (active), `disable` (stopped).
* `created_on` - Creation time of the DNS record.
* `modified_on` - Last modification time of the DNS record.

## Import

DNS record can be imported using the resource ID in the format `zone_id#record_id`, e.g.:

```shell
terraform import tencentcloud_teo_dns_record_v9.test zone123#record456
```
