Use this data source to query detailed information of vpc tenant_ccn

Example Usage

```hcl
data "tencentcloud_ccn_tenant_instances" "tenant_ccn" {
  ccn_ids = ["ccn-39lqkygf"]
  is_security_lock = ["true"]
}

```