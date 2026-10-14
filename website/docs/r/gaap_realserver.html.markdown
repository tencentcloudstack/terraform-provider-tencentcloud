---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_realserver"
sidebar_current: "docs-tencentcloud-resource-gaap_realserver"
description: |-
  Provides a resource to create a GAAP realserver.
---

# tencentcloud_gaap_realserver

Provides a resource to create a GAAP realserver.

## Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the GAAP realserver, the maximum length is 30.
* `domain` - (Optional, String, ForceNew) Domain of the GAAP realserver, conflict with `ip`.
* `ip` - (Optional, String, ForceNew) IP of the GAAP realserver, conflict with `domain`.
* `project_id` - (Optional, Int, ForceNew) ID of the project within the GAAP realserver, '0' means is default project.
* `tags` - (Optional, Map) Tags of the GAAP realserver.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddRealServers](https://www.tencentcloud.com/document/product/608/33061) - This API is used to add the information of the origin server (server), which supports IP or the domain name.
* Create: [DescribeRealServers](https://www.tencentcloud.com/document/product/608/33060) - This API is used to query origin server information. It can query all origin servers under the specified project name, and supports fuzzy query by specified IPs or domain names.
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Read: [DescribeRealServers](https://www.tencentcloud.com/document/product/608/33060) - This API is used to query origin server information. It can query all origin servers under the specified project name, and supports fuzzy query by specified IPs or domain names.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyRealServerName](https://www.tencentcloud.com/document/product/608/33058) - This API (ModifyRealServerName) is used to modify origin server names.
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [RemoveRealServers](https://www.tencentcloud.com/document/product/608/33057) - This API is used to delete the added origin server (server) IP or domain name.

## Import

GAAP realserver can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_realserver.foo rs-4ftghy6
```

