package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTcmMeshResource_basic -v
func TestAccTencentCloudTcmMeshResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMeshDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmMesh,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMeshExists("tencentcloud_tcm_mesh.basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "display_name", "test_mesh"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "mesh_version", "1.12.5"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "type", "HOSTED"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.disable_http_retry"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.disable_policy_checks"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.enable_pilot_http"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.outbound_traffic_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.smart_dns.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.smart_dns.0.istio_meta_dns_auto_allocate"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.istio.0.smart_dns.0.istio_meta_dns_capture"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.0.custom_prom.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.0.custom_prom.0.auth_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.0.custom_prom.0.is_public_addr"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.0.custom_prom.0.url"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.prometheus.0.custom_prom.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.0.enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.0.sampling"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.0.apm.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.0.apm.0.enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.tracing.0.zipkin.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.inject.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.inject.0.exclude_ip_ranges.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "config.0.inject.0.hold_application_until_proxy_starts", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "config.0.inject.0.hold_proxy_until_application_ends", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.sidecar_resources.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.sidecar_resources.0.limits.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "config.0.sidecar_resources.0.requests.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "tag_list.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "tag_list.0.key"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "tag_list.0.passthrough"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.basic", "tag_list.0.value"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_mesh.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMeshDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_mesh" {
			continue
		}

		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if err != nil {
			if isExpectError(err, []string{"ResourceNotFound"}) {
				return nil
			}
		}
		if mesh != nil {
			return fmt.Errorf("tcm mesh %v still exists", *mesh.Mesh.State)
		}
	}
	return nil
}

func testAccCheckMeshExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if mesh.Mesh == nil {
			return fmt.Errorf("tcm mesh %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmMesh = `

resource "tencentcloud_tcm_mesh" "basic" {
	display_name = "test_mesh"
	mesh_version = "1.12.5"
	type         = "HOSTED"
	config {
	  istio {
		outbound_traffic_policy = "ALLOW_ANY"
		disable_policy_checks   = true
		enable_pilot_http       = true
		disable_http_retry      = true
		smart_dns {
		  istio_meta_dns_capture       = true
		  istio_meta_dns_auto_allocate = true
		}
		tracing {
		  enable = false
		}
	  }
	  tracing {
		enable   = true
		sampling = 1
		apm {
		  enable = true
		  region = "ap-guangzhou"
		}
	  }
	  prometheus {
		custom_prom {
		  url       = "https://10.0.0.1:1000"
		  auth_type = "none"
		  vpc_id    = "vpc-j9yhbzpn"
		}
	  }
	  inject {
		exclude_ip_ranges                   = ["172.16.0.0/16"]
		hold_application_until_proxy_starts = true
		hold_proxy_until_application_ends   = true
	  }
  
	  sidecar_resources {
		limits {
		  name     = "cpu"
		  quantity = "2"
		}
		limits {
		  name     = "memory"
		  quantity = "1Gi"
		}
		requests {
		  name     = "cpu"
		  quantity = "100m"
		}
		requests {
		  name     = "memory"
		  quantity = "128Mi"
		}
	  }
	}
	tag_list {
	  key         = "key"
	  value       = "value"
	  passthrough = false
	}
}
`
