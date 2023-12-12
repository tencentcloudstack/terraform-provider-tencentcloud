Provides a resource to create a tsf cluster

Example Usage

```hcl
resource "tencentcloud_tsf_cluster" "cluster" {
	cluster_name = "terraform-test"
	cluster_type = "C"
	vpc_id = "vpc-xxxxxx"
	cluster_cidr = "9.165.120.0/24"
	cluster_desc = "test"
	tsf_region_id = "ap-guangzhou"
	cluster_version = "1.18.4"
	max_node_pod_num = 32
	max_cluster_service_num = 128
	tags = {
	  "createdBy" = "terraform"
	}
}
```