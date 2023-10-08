---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_ssh_key_pair_value"
sidebar_current: "docs-tencentcloud-datasource-ssm_ssh_key_pair_value"
description: |-
  Use this data source to query detailed information of ssm ssh_key_pair_value
---

# tencentcloud_ssm_ssh_key_pair_value

Use this data source to query detailed information of ssm ssh_key_pair_value

~> **NOTE:** Must set at least one of `secret_name` or `ssh_key_id`.

## Example Usage

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  secret_name = "keep_terraform"
  ssh_key_id  = "skey-2ae2snwd"
}
```

### Or

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  secret_name = "keep_terraform"
}
```

### Or

```hcl
data "tencentcloud_ssm_ssh_key_pair_value" "example" {
  ssh_key_id = "skey-2ae2snwd"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `secret_name` - (Optional, String) Secret name.
* `ssh_key_id` - (Optional, String) The key pair ID is the unique identifier of the key pair in the cloud server.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `private_key` - Private key plain text, encoded using base64.
* `project_id` - The project ID to which this key pair belongs.
* `public_key` - Public key plain text, encoded using base64.
* `ssh_key_description` - Description of the SSH key pair. Users can modify the description information of the key pair in the CVM console.
* `ssh_key_name` - SSH key name.


