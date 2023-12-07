Provides a resource to create a tcm prometheus_attachment

~> **NOTE:** Instructions for use: 1. Use Tencent Cloud Prometheus to monitor TMP, please enter `vpc_id`, `subnet_id`, `region` or `instance_id`, it is recommended to use an existing tmp instance; 2. To use the third-party Prometheus service, please enter `custom_prom`; 3. `tencentcloud_tcm_prometheus_attachment` does not support modification; 4. If you use Tencent Cloud Prometheus to monitor TMP, enter `vpc_id`, `subnet_id`, `region` to create a new Prometheus monitoring instance, destroy will not destroy the Prometheus monitoring instance
~> **NOTE:** If you use the config attribute prometheus in tencentcloud_tcm_mesh, do not use tencentcloud_tcm_prometheus_attachment

Example Usage

```hcl
resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
	mesh_id = "mesh-rofjmxxx"
	prometheus {
	  vpc_id = "vpc-pewdpxxx"
	  subnet_id = "subnet-driddxxx"
	  region = "ap-guangzhou"
	  instance_id = ""
	  # custom_prom {
		#   is_public_addr = false
		#   vpc_id = "vpc-pewdpxxx"
		#   url = "http://10.0.0.1:9090"
		#   auth_type = "basic"
		#   username = "test"
		#   password = "test"
	  # }
	}
}

```
Import

tcm prometheus_attachment can be imported using the mesh_id, e.g.
```
$ terraform import tencentcloud_tcm_prometheus_attachment.prometheus_attachment mesh-rofjmxxx
```