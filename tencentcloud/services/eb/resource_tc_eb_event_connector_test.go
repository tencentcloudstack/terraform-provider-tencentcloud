package eb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svceb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/eb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventConnectorResource_basic -v
func TestAccTencentCloudEbEventConnectorResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "type", "tdmq"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_connector.event_connector", "connection_description.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_connector.event_connector", "connection_description.0.resource_description"),
				),
			},
			{
				ResourceName: "tencentcloud_eb_event_connector.event_connector",
				ImportState:  true,
			},
		},
	})
}

func testAccCheckEbEventConnectorDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_connector" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		connectionId := idSplit[1]

		connector, err := service.DescribeEbEventConnectorById(ctx, connectionId, eventBusId)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound" {
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		connectionId := idSplit[1]

		service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

data "tencentcloud_user_info" "foo" {}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

locals {
  uin = data.tencentcloud_user_info.foo.owner_uin
}

resource "tencentcloud_eb_event_connector" "event_connector" {
  event_bus_id    = tencentcloud_eb_event_bus.foo.id
  connection_name = "tf-event-connector"
  description     = "event connector desc"
  enable          = true
  type            = "tdmq"

  connection_description {
    resource_description = "qcs::tdmq:ap-guangzhou:uin/${local.uin}:subscriptionName/pulsar-pb5xp2wwxepd/keep-tdmq-ns/keep-tdmq-topic1/test"
  }
}

`
