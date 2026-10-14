---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_shared_cname"
sidebar_current: "docs-tencentcloud-resource-teo_shared_cname"
description: |-
  Provides a resource to create a TencentCloud EdgeOne (TEO) shared CNAME
---

# tencentcloud_teo_shared_cname

Provides a resource to create a TencentCloud EdgeOne (TEO) shared CNAME

## Example Usage

```hcl
resource "tencentcloud_teo_shared_cname" "example" {
  zone_id             = "zone-39quuimqg8r6"
  shared_cname_prefix = "test-api"
  description         = "example shared cname"
}
```

## Argument Reference

The following arguments are supported:

* `shared_cname_prefix` - (Required, String, ForceNew) The shared CNAME prefix. Please enter a valid domain prefix, for example `test-api` or `test-api.com`, limited to 50 characters.
* `zone_id` - (Required, String, ForceNew) The zone ID of the shared CNAME.
* `description` - (Optional, String) Description. You can enter 1-50 characters.
* `ipssl_setting` - (Optional, List) IP SSL setting for the shared CNAME.

The `ipssl_setting` object supports the following:

* `associated_domain` - (Required, String) The domain associated with IP SSL.
* `status` - (Required, String) Association status. Valid values: `bound` (IP SSL configuration bound), `binding` (IP SSL configuration binding), `unbinding` (IP SSL configuration unbinding), `unbound` (IP SSL configuration unbound).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `shared_cname` - The full shared CNAME returned by the API.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSharedCNAME](https://cloud.tencent.com/document/api/1552/97762) - This API is used to create a shared CNAME. It is now only available to beta users.
* Create: [DescribeSharedCNAME](https://cloud.tencent.com/document/api/1552/130048) - Query the shared CNAME list, support fuzzy search, paginate, and sort.
* Create: [ModifySharedCNAME](https://cloud.tencent.com/document/api/1552/130047) - This API is used to modify a shared CNAME. Currently only support modifying the description of a shared CNAME and setting the associated IP SSL domain name of a shared CNAME with IP SSL type. The shared CNAME itself cannot be modified after creation. This feature is in beta test.
* Read: [DescribeSharedCNAME](https://cloud.tencent.com/document/api/1552/130048) - Query the shared CNAME list, support fuzzy search, paginate, and sort.
* Update: [DescribeSharedCNAME](https://cloud.tencent.com/document/api/1552/130048) - Query the shared CNAME list, support fuzzy search, paginate, and sort.
* Update: [ModifySharedCNAME](https://cloud.tencent.com/document/api/1552/130047) - This API is used to modify a shared CNAME. Currently only support modifying the description of a shared CNAME and setting the associated IP SSL domain name of a shared CNAME with IP SSL type. The shared CNAME itself cannot be modified after creation. This feature is in beta test.
* Delete: [DeleteSharedCNAME](https://cloud.tencent.com/document/api/1552/101362) - This API is used to delete a shared CNAME. It is now only available to beta users.

## Import

TEO shared CNAME can be imported using the composite id (zone_id#shared_cname), e.g.

```
terraform import tencentcloud_teo_shared_cname.example zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com
```

