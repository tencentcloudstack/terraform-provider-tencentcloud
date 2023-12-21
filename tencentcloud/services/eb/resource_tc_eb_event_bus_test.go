package eb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svceb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/eb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventBusResource_basic -v
func TestAccTencentCloudEbEventBusResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudEbEventBusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventBus,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudEbEventBusExists("tencentcloud_eb_event_bus.event_bus"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_bus.event_bus", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "event_bus_name", "tf-event_bus"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "description", "event bus desc"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "enable_store", "false"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "save_days", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_eb_event_bus.event_bus",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEbEventBusUp,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudEbEventBusExists("tencentcloud_eb_event_bus.event_bus"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_bus.event_bus", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "event_bus_name", "tf-event_bus-test"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "description", "event bus"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "enable_store", "true"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "save_days", "2"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_bus.event_bus", "tags.createdBy", "terraform-test"),
				),
			},
		},
	})
}

func testAccTencentCloudEbEventBusExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		securityGroup, err := service.DescribeEbEventBusById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if securityGroup == nil {
			return fmt.Errorf("eb eventBus %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccTencentCloudEbEventBusDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_bus" {
			continue
		}

		securityGroup, err := service.DescribeEbEventBusById(ctx, rs.Primary.ID)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound.EventBus" {
				return nil
			}
			return err
		}
		if securityGroup != nil {
			return fmt.Errorf("eb eventBus %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccEbEventBus = `

resource "tencentcloud_eb_event_bus" "event_bus" {
	event_bus_name = "tf-event_bus"
	description    = "event bus desc"
	enable_store   = false
	save_days      = 1
	tags = {
	  "createdBy" = "terraform"
	}
}

`

const testAccEbEventBusUp = `

resource "tencentcloud_eb_event_bus" "event_bus" {
	event_bus_name = "tf-event_bus-test"
	description    = "event bus"
	enable_store   = true
	save_days      = 2
	tags = {
	  "createdBy" = "terraform-test"
	}
}

`
