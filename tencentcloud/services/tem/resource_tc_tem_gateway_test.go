package tem_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctem "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTemGatewayResource_basic -v
func TestAccTencentCloudTemGatewayResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		// CheckDestroy: testAccCheckTemGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemGatewayExists("tencentcloud_tem_gateway.gateway"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_gateway.gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.ingress_name", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.environment_id", tcacctest.DefaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.address_ip_version", "IPV4"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rewrite_type", "NONE"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.mixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.0.host", "test.com"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.0.protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.0.http.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.0.http.0.paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.0.http.0.paths.0.path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.1.host", "hello.com"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.1.protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.1.http.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.1.http.0.paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_gateway.gateway", "ingress.0.rules.1.http.0.paths.0.path", "/"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_gateway.gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// func testAccCheckTemGatewayDestroy(s *terraform.State) error {
// 	logId := tccommon.GetLogId(tccommon.ContextNil)
// 	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
// 	service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "tencentcloud_tem_gateway" {
// 			continue
// 		}

// 		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
// 		if len(idSplit) != 2 {
// 			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
// 		}
// 		environmentId := idSplit[0]
// 		ingressName := idSplit[1]

// 		res, err := service.DescribeTemGateway(ctx, environmentId, ingressName)
// 		if err != nil {
// 			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
// 			if !ok {
// 				return err
// 			}
// 			if ee.Code == "InternalError.DefaultInternalError" {
// 				return nil
// 			}
// 			return err
// 		}

// 		if res != nil {
// 			return fmt.Errorf("tem gateway %s still exists", rs.Primary.ID)
// 		}
// 	}
// 	return nil
// }

func testAccCheckTemGatewayExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		ingressName := idSplit[1]

		service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTemGateway(ctx, environmentId, ingressName)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem gateway %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemGatewayVar = `
variable "environment_id" {
  default = "` + tcacctest.DefaultEnvironmentId + `"
}
`

const testAccTemGateway = testAccTemGatewayVar + `

resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
    ingress_name = "demo"
    environment_id = var.environment_id
    address_ip_version = "IPV4"
    rewrite_type = "NONE"
    mixed = false
    rules {
      host = "test.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "demo"
            service_port = 80
          }
        }
      }
    }
    rules {
      host = "hello.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "hello"
            service_port = 36000
          }
        }
      }
    }
  }
}
`
