Provides a resource to create a tsf instances_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_instances_attachment" "instances_attachment" {
  cluster_id = "cluster-123456"
  instance_id_list = [""]
  os_name = "Ubuntu 20.04"
  image_id = "img-123456"
  password = "MyP@ssw0rd"
  key_id = "key-123456"
  sg_id = "sg-123456"
  instance_import_mode = "R"
  os_customize_type = "my_customize"
  feature_id_list =
  instance_advanced_settings {
	mount_target = "/mnt/data"
	docker_graph_path = "/var/lib/docker"
  }
  security_group_ids = [""]
}
```