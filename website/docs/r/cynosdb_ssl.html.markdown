---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_ssl"
sidebar_current: "docs-tencentcloud-resource-cynosdb_ssl"
description: |-
  Provides a resource to create a cynosdb ssl
---

# tencentcloud_cynosdb_ssl

Provides a resource to create a cynosdb ssl

## Example Usage

```hcl
resource "tencentcloud_cynosdb_ssl" "example" {
  cluster_id  = "cynosdbmysql-1e0nzayx"
  instance_id = "cynosdbmysql-ins-pfsv6q1e"
  status      = "ON"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster id.
* `instance_id` - (Required, String) instance id.
* `status` - (Required, String) Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `download_url` - Certificate download address.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CloseSSL](https://cloud.tencent.com/document/api/1003/118379)
* Create: [DescribeSSLStatus](https://cloud.tencent.com/document/api/1003/118100)
* Create: [DescribeTasks](https://cloud.tencent.com/document/api/1003/103866)
* Create: [OpenSSL](https://cloud.tencent.com/document/api/1003/118099)
* Read: [DescribeSSLStatus](https://cloud.tencent.com/document/api/1003/118100)
* Update: [CloseSSL](https://cloud.tencent.com/document/api/1003/118379)
* Update: [DescribeSSLStatus](https://cloud.tencent.com/document/api/1003/118100)
* Update: [DescribeTasks](https://cloud.tencent.com/document/api/1003/103866)
* Update: [OpenSSL](https://cloud.tencent.com/document/api/1003/118099)

## Import

cynosdb ssl can be imported using the clusterId#instanceId, e.g.

```
terraform import tencentcloud_cynosdb_ssl.example cynosdbmysql-1e0nzayx#cynosdbmysql-ins-pfsv6q1e
```

