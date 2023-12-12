Provides a resource to create a cvm chc_config

Example Usage

```hcl
resource "tencentcloud_cvm_chc_config" "chc_config" {
  chc_id = "chc-xxxxxx"
  instance_name = "xxxxxx"
  bmc_user = "admin"
  password = "xxxxxx"
    bmc_virtual_private_cloud {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"

  }
  bmc_security_group_ids = ["sg-xxxxxx"]

  deploy_virtual_private_cloud {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"
  }
  deploy_security_group_ids = ["sg-xxxxxx"]
}
```

Import

cvm chc_config can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_config.chc_config chc_config_id
```