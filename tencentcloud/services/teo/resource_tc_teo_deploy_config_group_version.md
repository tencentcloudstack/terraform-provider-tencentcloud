Provides a resource to create a teo deploy config group version

Example Usage

```hcl
resource "tencentcloud_teo_deploy_config_group_version" "teo_deploy_config_group_version" {
  zone_id = "zone-2xkazzl8yf6k"
  env_id = "env-3lchxiq1h855"
  description = "Deploy config group version for production"
  # l7_acceleration
  config_group_version_infos {
    version_id = "ver-3lchxizh2mqn"
  }
  # edge_functions
  config_group_version_infos {
    version_id = "ver-3lchxjdciuzx"
  }
}
```