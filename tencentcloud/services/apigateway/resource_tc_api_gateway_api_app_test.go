package apigateway_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcapigateway "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apigateway"
)

var testAPIGatewayAPIAppResourceName = "tencentcloud_api_gateway_api_app"
var testAPIGatewayAPIAppResourceKey = testAPIGatewayAPIAppResourceName + ".example"

// go test -i; go test -test.run TestAccTencentCloudAPIGateWayAPIAppResource_basic -v
func TestAccTencentCloudAPIGateWayAPIAppResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIApp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIAppExists(testAPIGatewayAPIAppResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIAppResourceKey, "api_app_name", "tf_example"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIAppResourceKey, "api_app_desc", "app desc."),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIAppResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayAPIAppUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIAppExists(testAPIGatewayAPIAppResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIAppResourceKey, "api_app_name", "tf_example_update"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIAppResourceKey, "api_app_desc", "update app desc."),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayAPIAppDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIAppResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		apiApp, err := service.DescribeApiApp(ctx, rs.Primary.ID)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "InvalidParameterValue.InvalidCommandId" {
					return nil
				}
			}
			return err
		}

		if apiApp != nil {
			if *apiApp.TotalCount == 0 {
				return nil
			}
			return fmt.Errorf("api_gateway api_app %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIAppExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		command, err := service.DescribeApiApp(ctx, rs.Primary.ID)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "InvalidParameterValue.InvalidCommandId" {
					return nil
				}
			}
			return err
		}

		if command == nil {
			return fmt.Errorf("api_gateway api_app %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccAPIGatewayAPIApp = `
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example"
  api_app_desc = "app desc."

  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccAPIGatewayAPIAppUpdate = `
resource "tencentcloud_api_gateway_api_app" "example" {
  api_app_name = "tf_example_update"
  api_app_desc = "update app desc."

  tags = {
    "createdBy" = "terraform"
  }
}
`
