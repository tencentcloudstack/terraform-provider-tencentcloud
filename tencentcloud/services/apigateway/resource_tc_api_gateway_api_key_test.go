package apigateway_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcapigateway "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apigateway"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAPIGatewayAPIKeyResourceName = "tencentcloud_api_gateway_api_key"
var testAPIGatewayAPIKeyResourceKeyAuto = testAPIGatewayAPIKeyResourceName + ".example_auto"
var testAPIGatewayAPIKeyResourceKeyManual = testAPIGatewayAPIKeyResourceName + ".example_manual"

// go test -i; go test -test.run TestAccTencentCloudAPIGateWayAPIKeyResource -v
func TestAccTencentCloudAPIGateWayAPIKeyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIKeyAuto,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKeyAuto),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyAuto, "secret_name", "tf_example_auto"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyAuto, "status", "on"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "create_time"),
				),
			},
			{
				Config: testAccAPIGatewayAPIKeyAutoUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKeyAuto),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyAuto, "secret_name", "tf_example_auto"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyAuto, "status", "off"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyAuto, "create_time"),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIKeyResourceKeyAuto,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayAPIKeyManual,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKeyManual),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyManual, "secret_name", "tf_example_manual"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyManual, "status", "on"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "create_time"),
				),
			},
			{
				Config: testAccAPIGatewayAPIKeyManualUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyExists(testAPIGatewayAPIKeyResourceKeyManual),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyManual, "secret_name", "tf_example_manual"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIKeyResourceKeyManual, "status", "off"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "access_key_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "modify_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyResourceKeyManual, "create_time"),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIKeyResourceKeyManual,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIGatewayAPIKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIKeyResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeApiKey(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApiKey(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete API key for API gateway %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, has, err := service.DescribeApiKey(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApiKey(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("API key for API gateway %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccAPIGatewayAPIKeyAuto = `
resource "tencentcloud_api_gateway_api_key" "example_auto" {
  secret_name = "tf_example_auto"
  status      = "on"
}
`
const testAccAPIGatewayAPIKeyAutoUpdate = `
resource "tencentcloud_api_gateway_api_key" "example_auto" {
  secret_name = "tf_example_auto"
  status      = "off"
}
`

const testAccAPIGatewayAPIKeyManual = `
resource "tencentcloud_api_gateway_api_key" "example_manual" {
  secret_name       = "tf_example_manual"
  status            = "on"
  access_key_type   = "manual"
  access_key_id     = "28e287e340507fa147b2c8284dab542f"
  access_key_secret = "0198a4b8c3105080f4acd9e507599eff"
}
`
const testAccAPIGatewayAPIKeyManualUpdate = `
resource "tencentcloud_api_gateway_api_key" "example_manual" {
  secret_name       = "tf_example_manual"
  status            = "off"
  access_key_type   = "manual"
  access_key_id     = "28e287e340507fa147b2c8284dab542f"
  access_key_secret = "7f276e1cc4fa13cc82edae93c54b6e9b"
}
`
