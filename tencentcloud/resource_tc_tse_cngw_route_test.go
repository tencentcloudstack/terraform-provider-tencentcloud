package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwRouteResource_basic -v
func TestAccTencentCloudTseCngwRouteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseCngwRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwRoute,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwRouteExists("tencentcloud_tse_cngw_route.cngw_route"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_route.cngw_route", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "gateway_id", defaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "hosts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "https_redirect_status_code", "426"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.0.key", "req"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.0.value", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "preserve_host", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "protocols.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "route_name", "terraform-route"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "service_id", "b6017eaf-2363-481e-9e93-8d65aaf498cd"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "strip_path", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_route.cngw_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTseCngwRouteUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwRouteExists("tencentcloud_tse_cngw_route.cngw_route"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_route.cngw_route", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "gateway_id", defaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "hosts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "https_redirect_status_code", "301"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.0.key", "req"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "headers.0.value", "terraform1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "preserve_host", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "protocols.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "route_name", "terraform-route"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "service_id", "b6017eaf-2363-481e-9e93-8d65aaf498cd"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_route.cngw_route", "strip_path", "false"),
				),
			},
		},
	})
}

func testAccCheckTseCngwRouteDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_route" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		serviceID := idSplit[1]
		routeName := idSplit[2]

		res, err := service.DescribeTseCngwRouteById(ctx, gatewayId, serviceID, routeName)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwRoute %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwRouteExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		serviceID := idSplit[1]
		routeName := idSplit[2]

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTseCngwRouteById(ctx, gatewayId, serviceID, routeName)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwRoute %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwRoute = DefaultTseVar + `

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  destination_ports = []
  gateway_id        = var.gateway_id
  hosts = [
    "192.168.0.1:9090",
  ]
  https_redirect_status_code = 426
  paths = [
    "/user",
  ]
  headers {
  	key = "req"
  	value = "terraform"
  }
  preserve_host = false
  protocols = [
    "http",
    "https",
  ]
  route_name = "terraform-route"
  service_id = "b6017eaf-2363-481e-9e93-8d65aaf498cd"
  strip_path = true
}
`

const testAccTseCngwRouteUp = DefaultTseVar + `

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  destination_ports = []
  gateway_id        = var.gateway_id
  hosts = [
    "192.168.0.1:9091",
  ]
  https_redirect_status_code = 301
  paths = [
    "/user1",
  ]
  headers {
  	key = "req"
  	value = "terraform1"
  }
  preserve_host = true
  protocols = [
    "http",
  ]
  route_name = "terraform-route"
  service_id = "b6017eaf-2363-481e-9e93-8d65aaf498cd"
  strip_path = false
}
`
