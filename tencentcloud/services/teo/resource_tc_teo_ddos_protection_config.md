Provides a resource to create a TEO ddos protection config

~> **NOTE:** If `protection_option` is `protect_specified_domains`, then all domains need to be listed. For domains that do not need protection, just set the `switch` to `off`.

Example Usage

Protect all domains

```hcl
resource "tencentcloud_teo_ddos_protection_config" "example" {
  zone_id = "zone-3edjdliiw3he"
  ddos_protection {
    protection_option = "protect_all_domains"
  }
}
```

Protect designated domains

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

Import

TEO ddos protection config can be imported using the id, e.g.

```
terraform import tencentcloud_teo_ddos_protection_config.example zone-3edjdliiw3he
```
