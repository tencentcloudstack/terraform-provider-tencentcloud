---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zones"
sidebar_current: "docs-tencentcloud-datasource-teo_zones"
description: |-
  Use this data source to query detailed information of teo zoneAvailablePlans
---

# tencentcloud_teo_zones

Use this data source to query detailed information of teo zoneAvailablePlans

## Example Usage

```hcl
data "tencentcloud_teo_zones" "teo_zones" {
  filters {
    name   = "zone-id"
    values = ["zone-39quuimqg8r6"]
  }

  filters {
    name   = "tag-key"
    values = ["createdBy"]
  }

  filters {
    name   = "tag-value"
    values = ["terraform"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Optional, String) Sort direction. If the field value is a number, sort by the numeric value. If the field value is text, sort by the ascill code. Values include: `asc`: From the smallest to largest; `desc`: From the largest to smallest. Default value: `desc`.
* `filters` - (Optional, List) Filter criteria. the maximum value of Filters.Values is 20. if this parameter is left empty, all site information authorized under the current appid will be returned. detailed filter criteria are as follows: zone-name: filter by site name; zone-id: filter by site id. the site id is in the format of zone-2noz78a8ev6k; status: filter by site status; tag-key: filter by tag key; tag-value: filter by tag value; alias-zone-name: filter by identical site identifier. when performing a fuzzy query, the fields that support filtering are named zone-name or alias-zone-name.
* `order` - (Optional, String) Sort the returned results according to this field. Values include: `type`: Connection mode; `area`: Acceleration region; `create-time`: Creation time; `zone-name`: Site name; `use-time`: Last used time; `active-status` Effective status. Default value: `create-time`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value of the filtered field.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zones` - Details of sites.
  * `active_status` - Status of the proxy. Values: `active`: Enabled; `inactive`: Not activated; `paused`: Disabled.
  * `alias_zone_name` - The site alias. It can be up to 20 characters consisting of digits, letters, hyphens (-) and underscores (_). Note: This field may return null, indicating that no valid values can be obtained.
  * `area` - The site access region. Values: `global`: Global. `mainland`: Chinese mainland. `overseas`: Outside the Chinese mainland.
  * `cname_speed_up` - Whether CNAME acceleration is enabled. Values: `enabled`: Enabled; `disabled`: Disabled.
  * `cname_status` - CNAME record access status. Values: `finished`: The site is verified.`pending`: The site is being verified.
  * `created_on` - The creation time of the site.
  * `is_fake` - Whether it is a fake site. Valid values: `0`: Non-fake site; `1`: Fake site.
  * `lock_status` - Lock status. Values: `enable`: Normal. Modification is allowed. `disable`: Locked. Modification is not allowed. `plan_migrate`: Adjusting the plan. Modification is not allowed.
  * `modified_on` - The modification date of the site.
  * `name_servers` - The list of name servers assigned by Tencent Cloud.
  * `original_name_servers` - List of name servers used by the site.
  * `ownership_verification` - Ownership verification information. Note: This field may return null, indicating that no valid values can be obtained.
    * `dns_verification` - CNAME, when there is no domain name access, the information required for DNS resolution verification is used. For details, refer to [Site/Domain Ownership Verification
](https://intl.cloud.tencent.com/document/product/1552/70789?from_cn_redirect=1#7af6ecf8-afca-4e35-8811-b5797ed1bde5). Note: This field may return null, which indicates a failure to obtain a valid value.
      * `record_type` - The record type.
      * `record_value` - The record value.
      * `subdomain` - The host record.
    * `file_verification` - CNAME, when there is no domain name access, the information required for file verification is used. For details, refer to [Site/Domain Ownership Verification
](https://intl.cloud.tencent.com/document/product/1552/70789?from_cn_redirect=1#7af6ecf8-afca-4e35-8811-b5797ed1bde5). Note: This field may return null, which indicates a failure to obtain a valid value.
      * `content` - Content of the verification file. The contents of this field need to be filled into the text file returned by `Path`.
      * `path` - EdgeOne obtains the file verification information in the format of "Scheme + Host + URL Path", (e.g. https://www.example.com/.well-known/teo-verification/z12h416twn.txt). This field is the URL path section of the URL you need to create.
    * `ns_verification` - Information required for switching DNS servers. It's applicable to sites connected via NSs. For details, see [Modifying DNS Server](https://intl.cloud.tencent.com/document/product/1552/90452?from_cn_redirect=1).
Note: This field may return null, indicating that no valid values can be obtained.
      * `name_servers` - The DNS server address assigned to the user when connecting a site to EO via NS. You need to switch the NameServer of the domain name to this address.
  * `paused` - Whether the site is disabled.
  * `resources` - The list of billable resources.
    * `area` - Applicable area. Values: `mainland`: Chinese mainland; `overseas`: Regions outside the Chinese mainland; `global`: Global.
    * `auto_renew_flag` - Whether to enable auto-renewal. Values: `0`: Default status. `1`: Enable auto-renewal. `2`: Disable auto-renewal.
    * `create_time` - The creation time.
    * `enable_time` - The effective time.
    * `expire_time` - The expiration time.
    * `group` - The resource type. Values: `plan`: Plan resources; `pay-as-you-go`: Pay-as-you-go resources; `value-added`: Value-added resources. Note: This field may return null, indicating that no valid values can be obtained.
    * `id` - The resource ID.
    * `pay_mode` - Billing mode, `0`: Pay-as-you-go.
    * `plan_id` - ID of the resource associated with the plan.
    * `status` - The plan status. Values: `normal`: Normal; `isolated`: Isolated; `destroyed`: Terminated.
    * `sv` - Pricing query parameter.
      * `instance_id` - ID of the L4 proxy instance. Note: This field may return null, indicating that no valid values can be obtained.
      * `key` - The parameter key.
      * `pack` - Quota for a resource. Values: `zone`: Quota for sites; `custom-rule`: Quota for custom rules; `rate-limiting-rule`: Quota for rate limiting rules; `l4-proxy-instance`: Quota for L4 proxy instances. Note: This field may return null, indicating that no valid values can be obtained.
      * `protection_specs` - The protection specification. Values: `cm_30G`: 30 Gbps base protection bandwidth in **Chinese mainland** service area; `cm_60G`: 60 Gbps base protection bandwidth in **Chinese mainland** service area; `cm_100G`: 100 Gbps base protection bandwidth in **Chinese mainland** service area; `anycast_300G`: 300 Gbps Anycast-based protection in **Global (MLC)** service area; `anycast_unlimited`: Unlimited Anycast-based protection bandwidth in **Global (MLC)** service area; `cm_30G_anycast_300G`: 30 Gbps base protection bandwidth in **Chinese mainland** service area and 300 Gbps Anycast-based protection bandwidth in **Global (MLC)** service area; `cm_30G_anycast_unlimited`: 30 Gbps base protection bandwidth in **Chinese mainland** service area and unlimited Anycast-based protection bandwidth in **Global (MLC)** service area; cm_60G_anycast_300G`: 60 Gbps base protection bandwidth in **Chinese mainland** service area and 300 Gbps Anycast-based protection bandwidth in **Global (MLC)** service area; cm_60G_anycast_unlimited`: 60 Gbps base protection bandwidth in **Chinese mainland** service area and unlimited Anycast-based protection bandwidth in **Global (MLC)** service area</li><li> `cm_100G_anycast_300G`: 100 Gbps base protection bandwidth in **Chinese mainland** service area and 300 Gbps Anycast-based protection bandwidth in **Global (MLC)** service area, cm_100G_anycast_unlimited`: 100 Gbps base protection bandwidth in **Chinese mainland** service area and unlimited Anycast-based protection bandwidth in **Global (MLC)** service area. Note: This field may return null, indicating that no valid values can be obtained.
      * `value` - The parameter value.
    * `type` - Resource tag type. Valid values: vodeo: vodeo resource.
    * `zone_number` - The sites that are associated with the current resources. Note: This field may return null, indicating that no valid values can be obtained.
  * `status` - The site status. Values: `active`: The name server is switched to EdgeOne. `pending`: The name server is not switched. `moved`: The name server is changed to other service providers. `deactivated`: The site is blocked. `initializing`: The site is not bound with any plan.
  * `tags` - The list of resource tags.
    * `tag_key` - The tag key. Note: This field may return null, indicating that no valid values can be obtained.
    * `tag_value` - The tag value. Note: This field may return null, indicating that no valid values can be obtained.
  * `type` - Site access method. Valid values: full: NS access; partial: CNAME access; noDomainAccess: access with no domain name.
  * `vanity_name_servers_ips` - The custom name server IP information. Note: This field may return null, indicating that no valid values can be obtained.
    * `i_pv4` - IPv4 address of the custom name server.
    * `name` - Custom name of the name server.
  * `vanity_name_servers` - The custom name server information.
Note: This field may return null, indicating that no valid values can be obtained.
    * `servers` - List of custom name servers.
    * `switch` - Whether to enable custom name servers. Values: `on`: Enable; `off`: Disable.
  * `zone_id` - Site ID.
  * `zone_name` - The site name.


