---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ddos_protection_config"
sidebar_current: "docs-tencentcloud-resource-teo_ddos_protection_config"
description: |-
  Provides a resource to create a TEO ddos protection config
---

# tencentcloud_teo_ddos_protection_config

Provides a resource to create a TEO ddos protection config

~> **NOTE:** If `protection_option` is `protect_specified_domains`, then all domains need to be listed. For domains that do not need protection, just set the `switch` to `off`.

## Example Usage

### Protect all domains

```hcl
resource "tencentcloud_teo_ddos_protection_config" "example" {
  zone_id = "zone-3edjdliiw3he"
  ddos_protection {
    protection_option = "protect_all_domains"
  }
}
```

### Protect designated domains

```hcl
resource "tencentcloud_teo_ddos_protection_config" "example" {
  zone_id = "zone-3edjdliiw3he"
  ddos_protection {
    protection_option = "protect_specified_domains"
    domain_ddos_protections {
      domain = "1.demo.com"
      switch = "on"
    }

    domain_ddos_protections {
      domain = "2.demo.com"
      switch = "on"
    }

    domain_ddos_protections {
      domain = "3.demo.com"
      switch = "off"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ddos_protection` - (Required, List) Specifies the exclusive Anti-DDoS configuration.
* `zone_id` - (Required, String, ForceNew) Zone ID.

The `ddos_protection` object supports the following:

* `protection_option` - (Required, String) Specifies the protection scope of standalone DDoS. valid values:.
<li>protect_all_domains: specifies exclusive Anti-DDoS protection for all domain names in the site. newly added domain names automatically enable exclusive Anti-DDoS protection. when this parameter is specified, DomainDDoSProtections will not be processed.</li>.
<li>protect_specified_domains: only applicable to specified domains. specific scope can be set via DomainDDoSProtection parameter.</li>.
* `domain_ddos_protections` - (Optional, Set) Anti-DDoS configuration of the domain. specifies the exclusive ddos protection settings for the domain in request parameters.
<li>When ProtectionOption remains protect_specified_domains, the domain names not filled in keep their exclusive Anti-DDoS protection configuration unchanged, while explicitly specified domain names are updated according to the input parameters.</li>.
<li>When ProtectionOption switches from protect_all_domains to protect_specified_domains: if DomainDDoSProtections is empty, disable exclusive DDoS protection for all domains under the site; if DomainDDoSProtections is not empty, disable or maintain exclusive DDoS protection for the domain names specified in the parameter, and disable exclusive DDoS protection for other unlisted domain names.</li>.

The `domain_ddos_protections` object of `ddos_protection` supports the following:

* `domain` - (Required, String) Domain name.
* `switch` - (Required, String) Standalone DDoS switch of the domain. valid values:.
<li>on: enabled;</li>.
<li>off: closed.</li>.

The `shared_cname_ddos_protections` object of `ddos_protection` supports the following:


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO ddos protection config can be imported using the id, e.g.

```
terraform import tencentcloud_teo_ddos_protection_config.example zone-3edjdliiw3he
```

