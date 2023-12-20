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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var testAPIGatewayAPIDocResourceName = "tencentcloud_api_gateway_api_doc"
var testAPIGatewayAPIDocResourceKey = testAPIGatewayAPIDocResourceName + ".test"

// go test -i; go test -test.run TestAccTencentCloudAPIGateWayAPIDocResource_basic -v
func TestAccTencentCloudAPIGateWayAPIDocResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDocDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIDoc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIDocExists(testAPIGatewayAPIDocResourceKey),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIDocResourceKey, "api_doc_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIDocResourceKey, "service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIDocResourceKey, "environment"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIDocResourceKey, "api_ids.#"),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIDocResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAPIGatewayAPIDocUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIDocExists(testAPIGatewayAPIDocResourceKey),
					resource.TestCheckResourceAttr(testAPIGatewayAPIDocResourceKey, "api_doc_name", "update_doc_name_test"),
				),
			},
		},
	})
}

func testAccCheckAPIGatewayAPIDocDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIDocResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		apiDoc, err := service.DescribeApiDoc(ctx, rs.Primary.ID)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.InvalidApiDoc" || sdkerr.Code == "InvalidParameterValue.InvalidCommandId" {
					return nil
				}
			}
			return err
		}

		if apiDoc != nil {
			return fmt.Errorf("api_gateway api_doc %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckAPIGatewayAPIDocExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		command, err := service.DescribeApiDoc(ctx, rs.Primary.ID)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.InvalidApiDoc" || sdkerr.Code == "InvalidParameterValue.InvalidCommandId" {
					return nil
				}
			}
			return err
		}

		if command == nil {
			return fmt.Errorf("api_gateway api_doc %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccAPIGatewayAPIDoc = `
resource "tencentcloud_api_gateway_api_doc" "test" {
  api_doc_name = "doc_test1"
  service_id   = "service-nxz6yync"
  environment  = "release"
  api_ids      = ["api-jvqlzolk"]
}
`

const testAccAPIGatewayAPIDocUpdate = `
resource "tencentcloud_api_gateway_api_doc" "test" {
  api_doc_name = "update_doc_name_test"
  service_id   = "service-nxz6yync"
  environment  = "release"
  api_ids      = ["api-jvqlzolk"]
}
`
