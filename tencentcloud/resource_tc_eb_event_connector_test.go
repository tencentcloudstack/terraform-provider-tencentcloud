package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudEbEventConnectorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEbEventConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventConnector,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEbEventConnectorExists("tencentcloud_eb_event_connector.event_connector"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_connector.event_connector", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_connector.event_connector", "event_bus_id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "connection_name", "tf-event-connector"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "description", "event connector desc"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "type", "apigw"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "connection_description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "connection_description.0.method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "connection_description.0.protocol", "HTTP"),
				),
			},
			{
				ResourceName:      "tencentcloud_eb_event_connector.event_connector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEbEventConnectorDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_connector" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		connectionId := idSplit[1]

		connector, err := service.DescribeEbEventConnectorById(ctx, connectionId, eventBusId)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound.Rule" {
				return nil
			}
			return err
		}

		if connector != nil {
			return fmt.Errorf("eb eventConnector %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEbEventConnectorExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		connectionId := idSplit[1]

		service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		connector, err := service.DescribeEbEventConnectorById(ctx, connectionId, eventBusId)
		if err != nil {
			return err
		}

		if connector == nil {
			return fmt.Errorf("eb eventConnector %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccEbEventConnector = `

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_api_gateway_service" "service" {
  service_name = "tf-eb-service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

locals {
  service_id = tencentcloud_api_gateway_service.service.id
}

resource "tencentcloud_eb_event_connector" "event_connector" {
  event_bus_id    = tencentcloud_eb_event_bus.foo.id
  connection_name = "tf-event-connector"
  description     = "event connector desc1"
  enable          = false
  type            = "apigw"
  connection_description {
    resource_description = "qcs::apigw:ap-guangzhou:uin/100022975249:serviceid/${local.service_id}"
    api_gw_params {
      protocol = "HTTP"
      method   = "GET"
    }
  }
}

`
