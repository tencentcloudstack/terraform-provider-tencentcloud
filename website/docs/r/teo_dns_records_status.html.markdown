---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_dns_records_status"
sidebar_current: "docs-tencentcloud-resource-teo_dns_records_status"
description: |-
  Provides a resource to manage TEO (EdgeOne) DNS records status config.
---

# tencentcloud_teo_dns_records_status

Provides a resource to manage TEO (EdgeOne) DNS records status config.

~> **NOTE:** This resource is currently used to manage default DNS resolution. The same DNS record cannot be used simultaneously with the tencentcloud_teo_dns_record resource.

## Example Usage

### Enable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id    = "zone-3edjdliiw3he"
  records_id = "record-1234567890"
  status     = "enable"
}
```

### Disable a DNS record

```hcl
resource "tencentcloud_teo_dns_records_status" "example" {
  zone_id    = "zone-3edjdliiw3he"
  records_id = "record-1234567890"
  status     = "disable"
}
```

## Argument Reference

The following arguments are supported:

* `records_id` - (Required, String, ForceNew) DNS record ID, combined with `zone_id` as the unique ID of the resource.
* `status` - (Required, String) DNS record status. Valid values: `enable` (enabled), `disable` (disabled).
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDnsRecords](https://cloud.tencent.com/document/api/1552/80716) - This API is used to query DNS records. Paging, sorting and filtering are supported.
* Create: [ModifyDnsRecordsStatus](https://cloud.tencent.com/document/api/1552/114251) - You can batch modify the status of DNS records through this API, enabling and disabling records in bulk.
* Read: [DescribeDnsRecords](https://cloud.tencent.com/document/api/1552/80716) - This API is used to query DNS records. Paging, sorting and filtering are supported.
* Update: [DescribeDnsRecords](https://cloud.tencent.com/document/api/1552/80716) - This API is used to query DNS records. Paging, sorting and filtering are supported.
* Update: [ModifyDnsRecordsStatus](https://cloud.tencent.com/document/api/1552/114251) - You can batch modify the status of DNS records through this API, enabling and disabling records in bulk.

## Import

TEO (EdgeOne) DNS records status config can be imported using the composite id `zone_id#records_id`, e.g.

```
terraform import tencentcloud_teo_dns_records_status.example zone-3edjdliiw3he#record-1234567890
```

