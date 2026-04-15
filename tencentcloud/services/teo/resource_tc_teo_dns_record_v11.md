# terraform-provider-tencentcloud

## Resource: tencentcloud_teo_dns_record_v11

Provides a TEO DNS record v11 resource.

## Example Usage

### Basic Example

```hcl
resource "tencentcloud_teo_dns_record_v11" "example" {
  zone_id      = "zone-xxxxxx"
  domain       = "example.com"
  record_type  = "A"
  record_value = "1.2.3.4"
  ttl          = 300
}
```

### A Record with Weight

```hcl
resource "tencentcloud_teo_dns_record_v11" "a_record_weight" {
  zone_id      = "zone-xxxxxx"
  domain       = "www.example.com"
  record_type  = "A"
  record_value = "1.2.3.4"
  ttl          = 300
  weight       = 10
}
```

### CNAME Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "cname_record" {
  zone_id      = "zone-xxxxxx"
  domain       = "www.example.com"
  record_type  = "CNAME"
  record_value = "cname.example.com"
  ttl          = 300
}
```

### MX Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "mx_record" {
  zone_id      = "zone-xxxxxx"
  domain       = "example.com"
  record_type  = "MX"
  record_value = "mail.example.com"
  priority     = 10
  ttl          = 300
}
```

### TXT Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "txt_record" {
  zone_id      = "zone-xxxxxx"
  domain       = "_dmarc.example.com"
  record_type  = "TXT"
  record_value = "v=DMARC1; p=none; rua=mailto:dmarc@example.com"
  ttl          = 300
}
```

### AAAA Record

```hcl
resource "tencentcloud_teo_dns_record_v11" "aaaa_record" {
  zone_id      = "zone-xxxxxx"
  domain       = "www.example.com"
  record_type  = "AAAA"
  record_value = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  ttl          = 300
}
```

### A Record with Location

```hcl
resource "tencentcloud_teo_dns_record_v11" "a_record_location" {
  zone_id      = "zone-xxxxxx"
  domain       = "www.example.com"
  record_type  = "A"
  record_value = "1.2.3.4"
  ttl          = 300
  location     = "Asia"
}
```

### Complete Example with All Parameters

```hcl
resource "tencentcloud_teo_dns_record_v11" "complete" {
  zone_id      = "zone-xxxxxx"
  domain       = "www.example.com"
  record_type  = "A"
  record_value = "1.2.3.4"
  ttl          = 600
  priority     = 0
  weight       = 10
  location     = "Default"
}
```

## Update Example

```hcl
resource "tencentcloud_teo_dns_record_v11" "example" {
  zone_id      = "zone-xxxxxx"
  domain       = "example.com"
  record_type  = "A"
  record_value = "1.2.3.4"
  ttl          = 300
}

# Update the record value
resource "tencentcloud_teo_dns_record_v11" "example" {
  zone_id      = "zone-xxxxxx"
  domain       = "example.com"
  record_type  = "A"
  record_value = "5.6.7.8"
  ttl          = 600
}
```

## Delete Example

To delete a DNS record, simply remove the resource from your Terraform configuration:

```hcl
# resource "tencentcloud_teo_dns_record_v11" "example" {
#   zone_id      = "zone-xxxxxx"
#   domain       = "example.com"
#   record_type  = "A"
#   record_value = "1.2.3.4"
# }
```

Or use `terraform destroy` command.

## Import Example

You can import an existing DNS record using its ID:

```bash
terraform import tencentcloud_teo_dns_record_v11.example zone-xxxxxx#record-yyyyyy
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID.
* `domain` - (Required, ForceNew) DNS record name. If domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.
* `record_type` - (Required, ForceNew) DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.
* `record_value` - (Required) DNS record content. Fill in corresponding content according to Type value.
* `ttl` - (Optional, Default: 300) Cache time. The value range is 60~86400, in seconds. The smaller value, faster modification takes effect globally.
* `priority` - (Optional, Default: 0) MX record priority. This parameter is only valid when Type (DNS record type) is MX. The value range is 0~50. The smaller value, higher the priority.
* `weight` - (Optional, Default: -1) DNS record weight. The value range is -1~100. -1 means no weight is set, and 0 means no resolution. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.
* `location` - (Optional, Default: "Default") DNS record resolution line. Default is Default, indicating default resolution line, which takes effect in all regions. Resolution line configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_id` - DNS record ID.
* `status` - DNS record resolution status. Valid values: enable (active), disable (inactive).
* `created_at` - Creation time.
* `updated_at` - Modification time.

## Timeouts

The `timeouts` block allows you to specify timeouts for certain operations:

* `create` - (Default: 10 minutes) Timeout for creating DNS record.
* `read` - (Default: 3 minutes) Timeout for reading DNS record.
* `update` - (Default: 10 minutes) Timeout for updating DNS record.
* `delete` - (Default: 10 minutes) Timeout for deleting DNS record.

## Notes

1. **Record Type Constraints**:
   - `weight` parameter is only applicable when `record_type` is A, AAAA, or CNAME.
   - `priority` parameter is only applicable when `record_type` is MX.
   - `location` parameter is only applicable when `record_type` is A, AAAA, or CNAME.

2. **ID Format**:
   The resource ID is composed of `zone_id` and `record_id`, separated by `#`. For example: `zone-xxxxxx#record-yyyyyy`.

3. **Domain Name Encoding**:
   If the domain name is in Chinese, Korean, or Japanese, it must be converted to punycode before being used as the `domain` parameter.

4. **Weight Configuration**:
   For the same subdomain and resolution line, different DNS records should maintain the same weight configuration (either all have weights set or all have no weights set).

5. **Asynchronous Operations**:
   DNS record operations are asynchronous. The resource waits for the operation to complete before returning. You can adjust the timeout values using the `timeouts` block.

6. **Importing**:
   When importing an existing DNS record, you need to provide both the zone ID and record ID in the format `zone_id#record_id`.
