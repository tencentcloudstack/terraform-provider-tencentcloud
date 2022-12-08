package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTcmPrometheusAttachmentResource_basic -v
func TestAccTencentCloudTcmPrometheusAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrometheusAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmPrometheusAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrometheusAttachmentExists("tencentcloud_tcm_prometheus_attachment.prometheus_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "mesh_id", defaultMeshId),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.auth_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.is_public_addr", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.url", "http://10.0.0.1:9090"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.username", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_prometheus_attachment.prometheus_attachment", "prometheus.0.custom_prom.0.vpc_id", defaultMeshVpcId),
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_prometheus_attachment" {
			continue
		}

		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if err != nil {
			if isExpectError(err, []string{"ResourceNotFound"}) {
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
variable "mesh_id" {
  default = "` + defaultMeshId + `"
}
variable "vpc_id" {
  default = "` + defaultMeshVpcId + `"
}
`

const testAccTcmPrometheusAttachment = testAccTcmPrometheusAttachmentVar + `

resource "tencentcloud_tcm_prometheus_attachment" "prometheus_attachment" {
	mesh_id = var.mesh_id
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
