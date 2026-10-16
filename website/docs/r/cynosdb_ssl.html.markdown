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

* Create: [CloseSSL](https://www.tencentcloud.com/search/?wd=CloseSSL&type=document&lang=en)
* Create: [DescribeSSLStatus](https://www.tencentcloud.com/search/?wd=DescribeSSLStatus&type=document&lang=en)
* Create: [DescribeTasks](https://www.tencentcloud.com/search/?wd=DescribeTasks&type=document&lang=en)
* Create: [OpenSSL](https://www.tencentcloud.com/search/?wd=OpenSSL&type=document&lang=en)
* Read: [DescribeSSLStatus](https://www.tencentcloud.com/search/?wd=DescribeSSLStatus&type=document&lang=en)
* Update: [CloseSSL](https://www.tencentcloud.com/search/?wd=CloseSSL&type=document&lang=en)
* Update: [DescribeSSLStatus](https://www.tencentcloud.com/search/?wd=DescribeSSLStatus&type=document&lang=en)
* Update: [DescribeTasks](https://www.tencentcloud.com/search/?wd=DescribeTasks&type=document&lang=en)
* Update: [OpenSSL](https://www.tencentcloud.com/search/?wd=OpenSSL&type=document&lang=en)

## Import

cynosdb ssl can be imported using the clusterId#instanceId, e.g.

```
terraform import tencentcloud_cynosdb_ssl.example cynosdbmysql-1e0nzayx#cynosdbmysql-ins-pfsv6q1e
```

