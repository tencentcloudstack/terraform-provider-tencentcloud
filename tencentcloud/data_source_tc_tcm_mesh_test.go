package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTcmMeshDataSource_basic -v
func TestAccTencentCloudTcmMeshDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmMeshDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcm_mesh.mesh"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.display_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.type"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.config.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.config.0.istio.#", "1"),
					// resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.config.0.tracing.#", "1"),
					// resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.config.0.prometheus.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcm_mesh.mesh", "mesh_list.0.tag_list.#", "1"),
				),
			},
		},
	})
}

const testAccTcmMeshDataSource = `

resource "tencentcloud_tcm_mesh" "basic" {
	display_name = "test_mesh"
	mesh_version = "1.12.5"
	type = "HOSTED"
	config {
	  istio {
		outbound_traffic_policy = "ALLOW_ANY"
		disable_policy_checks = true
		enable_pilot_http = true
		disable_http_retry = true
		smart_dns {
		  istio_meta_dns_capture = true
		  istio_meta_dns_auto_allocate = true
		}
	  }
	  tracing {
		  enable = true
		  sampling = 1
		  apm {
			  enable = false
		  }
		  zipkin {
			  address = "10.0.0.1:1000"
		  }
	  }
	  prometheus {
		  custom_prom {
			  url = "https://10.0.0.1:1000"
			  auth_type = "none"
			  vpc_id = "vpc-j9yhbzpn"
		  }
	  }
	}
	tag_list {
	  key = "key"
	  value = "value"
	  passthrough = false
	}
  }

data "tencentcloud_tcm_mesh" "mesh" {
  mesh_id = [tencentcloud_tcm_mesh.basic.id]
  mesh_name = ["test_mesh"]
  tags = ["key"]
}

`
