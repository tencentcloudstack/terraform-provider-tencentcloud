---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_authentication_options"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_authentication_options"
description: |-
  Use this data source to query detailed information of kubernetes cluster_authentication_options
---

# tencentcloud_kubernetes_cluster_authentication_options

Use this data source to query detailed information of kubernetes cluster_authentication_options

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_authentication_options" "cluster_authentication_options" {
  cluster_id = "cls-kzilgv5m"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `latest_operation_state` - Result of the last modification. Values: `Updating`, `Success`, `Failed` or `TimeOut`. Note: this field may return `null`, indicating that no valid values can be obtained.
* `oidc_config` - OIDC authentication configurations. Note: This field may return `null`, indicating that no valid value can be obtained.
  * `auto_create_client_id` - Creating ClientId of the identity provider. Note: This field may return `null`, indicating that no valid value can be obtained.
  * `auto_create_oidc_config` - Creating an identity provider. Note: This field may return `null`, indicating that no valid value can be obtained.
  * `auto_install_pod_identity_webhook_addon` - Creating the PodIdentityWebhook component. Note: This field may return `null`, indicating that no valid value can be obtained.
* `service_accounts` - ServiceAccount authentication configuration. Note: this field may return `null`, indicating that no valid values can be obtained.
  * `auto_create_discovery_anonymous_auth` - If it is set to `true`, a RABC rule is automatically created to allow anonymous users to access `/.well-known/openid-configuration` and `/openid/v1/jwks`. Note: this field may return `null`, indicating that no valid values can be obtained.
  * `issuer` - service-account-issuer. Note: this field may return `null`, indicating that no valid values can be obtained.
  * `jwks_uri` - service-account-jwks-uri. Note: this field may return `null`, indicating that no valid values can be obtained.
  * `use_tke_default` - Use TKE default issuer and jwksuri. Note: This field may return `null`, indicating that no valid values can be obtained.


