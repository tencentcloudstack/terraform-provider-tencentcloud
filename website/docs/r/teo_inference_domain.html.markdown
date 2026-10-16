---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_domain"
sidebar_current: "docs-tencentcloud-resource-teo_inference_domain"
description: |-
  Provides a resource to create a TEO inference domain.
---

# tencentcloud_teo_inference_domain

Provides a resource to create a TEO inference domain.

~> **NOTE:** Before creating an inference domain, please verify the domain ownership of the site.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  service_id  = "sid-xxxxxx"
  domain      = "www.demo.com"
  auth_switch = "On"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Inference service domain.
* `service_id` - (Required, String, ForceNew) Inference service ID.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `auth_switch` - (Optional, String, ForceNew) Inference task request authentication switch. Valid values: Off: disable authentication; On: enable authentication. Default value: On.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domains` - Inference service domain list.
  * `auth_switch` - Inference task request authentication switch. Valid values: Off: disable authentication; On: enable authentication. Default value: On.
  * `certificate` - Domain certificate information.
    * `apply_type` - Application type. Valid values: apply: managed by EdgeOne; none: not managed by EdgeOne. If not specified, the default value is none.
    * `host_cert_info` - Server certificate configuration.
      * `alias` - Certificate alias name.
      * `cert_id` - Server certificate ID.
      * `deploy_time` - Certificate deployment time.
      * `expire_time` - Certificate expiration time.
      * `sign_algo` - Signature algorithm.
      * `status` - Certificate status. Valid values: deployed: deployed; processing: deploying; applying: applying; failed: application failed; issued: binding failed.
      * `type` - Certificate type. Valid values: default: default certificate; upload: user uploaded; managed: Tencent Cloud managed.
    * `host` - Domain.
    * `mode` - Certificate configuration mode. Valid values: disable: no certificate; eofreecert: EdgeOne free certificate; sslcert: SSL certificate.
  * `cname` - CNAME address.
  * `create_time` - Creation time.
  * `domain` - Domain name.
  * `ownership_verification` - When the domain needs to perform ownership verification to continue providing services, this object will carry the information required by the corresponding verification method.
    * `dns_verification` - Information required when using DNS resolution verification for CNAME or no-domain access.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `subdomain` - Host record.
    * `file_verification` - Information required when using file verification for CNAME or no-domain access.
      * `content` - Content of the verification file.
      * `path` - URL Path part that needs to be created.
    * `ns_verification` - Information required to switch DNS servers for NS access.
      * `name_servers` - DNS server addresses assigned to the user for NS access.
  * `status` - Domain status. Valid values: Online: effective; Process: deploying; Offline: disabled; Init: not effective, waiting to activate the site; Forbidden: domain is forbidden.
  * `update_time` - Modification time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `read` - (Defaults to `3m`) Used when reading the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

TEO inference domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_inference_domain.example zone-39quuimqg8r6#sid-xxxxxx#www.demo.com
```

