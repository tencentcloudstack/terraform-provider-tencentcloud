package tcm_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTcmPrometheusAttachmentResource_basic -v
func TestAccTencentCloudTcmPrometheusAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPrometheusAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmPrometheusAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrometheusAttachmentExists("tencentcloud_tcm_prometheus_attachment.prometheus_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "mesh_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.auth_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.is_public_addr", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.url", "http://10.0.0.1:9090"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.username", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.vpc_id", tcacctest.DefaultMeshVpcId),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_prometheus_attachment.prometheus_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPrometheusAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctcm.NewTcmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_prometheus_attachment" {
			continue
		}

		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if err != nil {
			if tccommon.IsExpectError(err, []string{"ResourceNotFound"}) {
				return nil
			}
		}
		if mesh.Mesh.Config.Prometheus != nil {
			return fmt.Errorf("tcm prometheusAttachment %v still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPrometheusAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctcm.NewTcmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if mesh.Mesh.Config.Prometheus == nil {
			return fmt.Errorf("tcm prometheusAttachment %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmPrometheusAttachmentVar = `
variable "vpc_id" {
  default = "` + tcacctest.DefaultMeshVpcId + `"
}
`

const testAccTcmPrometheusAttachment = testAccTcmPrometheusAttachmentVar + `

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
	}
	tag_list {
	  key = "key"
	  value = "value"
	  passthrough = false
	}
  }

resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
	mesh_id = tencentcloud_tcm_mesh.basic.id
	prometheus {
	  	# vpc_id = "vpc-pewdpxxx"
	  	# subnet_id = "subnet-driddxxx"
	  	# region = "ap-guangzhou"
	  	# instance_id = ""
		custom_prom {
		  is_public_addr = false
		  vpc_id = var.vpc_id
		  url = "http://10.0.0.1:9090"
		  auth_type = "basic"
		  username = "test"
		  password = "test"
	    }
	}
}

`
