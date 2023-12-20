package apigateway_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcapigateway "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apigateway"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAPIGatewayAPIKeyAttachmentResourceName = "tencentcloud_api_gateway_api_key_attachment"
var testAPIGatewayAPIKeyAttachmentResourceKey = testAPIGatewayAPIKeyAttachmentResourceName + ".attach"

func TestAccTencentCloudAPIGateWayAPIKeyAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIKeyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayAPIKeyAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayAPIKeyAttachmentExists(testAPIGatewayAPIKeyAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyAttachmentResourceKey, "api_key_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIKeyAttachmentResourceKey, "usage_plan_id"),
				),
			},
			{
				ResourceName:      testAPIGatewayAPIKeyAttachmentResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIGatewayAPIKeyAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayAPIKeyAttachmentResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		apiKeyId := idSplit[0]
		usagePlanId := idSplit[1]
		if apiKeyId == "" || usagePlanId == "" {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		info, has, err := service.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			info, has, err = service.DescribeUsagePlan(ctx, usagePlanId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		for _, v := range info.BindSecretIds {
			if *v == apiKeyId {
				return fmt.Errorf("unattach API key %s fail, still on server", rs.Primary.ID)
			}
		}

		return nil

	}
	return nil
}

func testAccCheckAPIGatewayAPIKeyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		apiKeyId := idSplit[0]
		usagePlanId := idSplit[1]
		if apiKeyId == "" || usagePlanId == "" {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		info, has, err := service.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			info, has, err = service.DescribeUsagePlan(ctx, usagePlanId)
		}
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("usage plan %s is not exist", usagePlanId)
		}

		for _, v := range info.BindSecretIds {
			if *v == apiKeyId {
				return nil
			}
		}
		return fmt.Errorf("attach API key %s fail, still on server", rs.Primary.ID)
	}
}

const testAccAPIGatewayAPIKeyAttachment = `
resource "tencentcloud_api_gateway_api_key" "key" {
  secret_name = "my_api_key"
  status      = "on"
}

resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_api_key_attachment" "attach" {
  api_key_id    = tencentcloud_api_gateway_api_key.key.id
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
`
