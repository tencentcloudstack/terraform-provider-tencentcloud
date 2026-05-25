---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_temp_credential"
sidebar_current: "docs-tencentcloud-ephemeral_resource-temp_credential"
description: |-
  Provides an ephemeral resource that returns a synthetic short-lived
credential bundle. The returned value lives only inside a single
Terraform graph walk and is never persisted to the state file. This
reference is suitable for handing transient secrets to downstream
provider configuration blocks.
---

# tencentcloud_temp_credential

Provides an ephemeral resource that returns a synthetic short-lived
credential bundle. The returned value lives only inside a single
Terraform graph walk and is never persisted to the state file. This
reference is suitable for handing transient secrets to downstream
provider configuration blocks.

## Example Usage

```hcl
ephemeral "tencentcloud_temp_credential" "this" {
  ttl_seconds = 900
}

provider "vault" {
  address = "https://vault.example.com"
  token   = ephemeral.tencentcloud_temp_credential.this.session_token
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Region the placeholder credential is bound to. Falls back to the provider's configured region when omitted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `expires_at` - RFC3339 timestamp 5 minutes after Open. Mirrors the shape of a real short-term credential.
* `secret_id` - Locally-constructed secret id, prefixed with "STS-fake-". Not a real credential.
* `secret_key` - Locally-constructed random hex string. Not a real credential.
* `token` - Locally-constructed random hex string. Not a real credential.


