---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eks_cluster_credential"
sidebar_current: "docs-tencentcloud-datasource-eks_cluster_credential"
description: |-
  Provide a datasource to query EKS cluster credential info.
---

# tencentcloud_eks_cluster_credential

Provide a datasource to query EKS cluster credential info.

## Example Usage

```hcl
data "tencentcloud_eks_cluster_credential" "foo" {
  cluster_id = "cls-xxxxxxxx"
}

# example outputs
output "addresses" {
  value = data.tencentcloud_eks_cluster_credential.cred.addresses
}

output "ca_cert" {
  value = data.tencentcloud_eks_cluster_credential.cred.credential.ca_cert
}

output "token" {
  value = data.tencentcloud_eks_cluster_credential.cred.credential.token
}

output "public_lb_param" {
  value = data.tencentcloud_eks_cluster_credential.cred.public_lb.0.extra_param
}

output "internal_lb_subnet" {
  value = data.tencentcloud_eks_cluster_credential.cred.internal_lb.0.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) EKS Cluster ID.
* `result_output_file` - (Optional, String) Used for save result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `addresses` - List of IP Address information.
  * `ip` - IP Address.
  * `port` - Port.
  * `type` - Type of IP, can be `advertise`, `public`, etc.
* `credential` - Credential info.
  * `ca_cert` - CA root certification.
  * `token` - Certification token.
* `internal_lb` - Cluster internal access LoadBalancer info.
  * `enabled` - Indicates weather the internal access LB enabled.
  * `subnet_id` - ID of subnet which related to Internal LB.
* `kube_config` - EKS cluster kubeconfig.
* `proxy_lb` - Indicates whether the new internal/public network function.
* `public_lb` - Cluster public access LoadBalancer info.
  * `allow_from_cidrs` - List of CIDRs which allowed to access.
  * `enabled` - Indicates weather the public access LB enabled.
  * `extra_param` - Extra param text json.
  * `security_group` - Security group.
  * `security_policies` - List of security allow IP or CIDRs, default deny all.


