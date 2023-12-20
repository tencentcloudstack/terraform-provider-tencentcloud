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
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

var testAPIGatewayUsagePlanAttachmentResourceName = "tencentcloud_api_gateway_usage_plan_attachment"
var testAPIGatewayUsagePlanAttachmentResourceKey = testAPIGatewayUsagePlanAttachmentResourceName + ".example"

// go test -i; go test -test.run TestAccTencentCloudAPIGateWayUsagePlanAttachmentResource -v
func TestAccTencentCloudAPIGateWayUsagePlanAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAPIGatewayUsagePlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayUsagePlanAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAPIGatewayUsagePlanAttachmentExists(testAPIGatewayUsagePlanAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanAttachmentResourceKey, "service_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayUsagePlanAttachmentResourceKey, "usage_plan_id"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanAttachmentResourceKey, "environment", "release"),
					resource.TestCheckResourceAttr(testAPIGatewayUsagePlanAttachmentResourceKey, "bind_type", "API"),
				),
			},
			{
				ResourceName:      testAPIGatewayUsagePlanAttachmentResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func base(ctx context.Context, rs *terraform.ResourceState) (plans []*apigateway.ApiUsagePlan, err error) {
	ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
	if len(ids) != 6 {
		return nil, fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
	}

	var (
		usagePlanId = ids[0]
		serviceId   = ids[1]
		environment = ids[2]
		bindType    = ids[3]
		apiId       = ids[4]
		//accessKeysStr = ids[5]
		outErr error
		has    bool
	)

	if usagePlanId == "" || serviceId == "" || environment == "" || bindType == "" {
		return nil, fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
	}
	if bindType == svcapigateway.API_GATEWAY_TYPE_API && apiId == "" {
		return nil, fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
	}

	service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId)
	if outErr != nil {
		_, has, outErr = service.DescribeUsagePlan(ctx, usagePlanId)
	}
	if outErr != nil {
		return nil, outErr
	}
	if !has {
		return nil, nil
	}

	_, has, outErr = service.DescribeService(ctx, serviceId)
	if outErr != nil {
		_, has, outErr = service.DescribeService(ctx, serviceId)
	}
	if outErr != nil {
		return nil, outErr
	}
	if !has {
		return nil, nil
	}

	if bindType == svcapigateway.API_GATEWAY_TYPE_API {
		plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
		if outErr != nil {
			plans, outErr = service.DescribeApiUsagePlan(ctx, serviceId)
		}
		if outErr != nil {
			return nil, outErr
		}
	} else {
		plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
		if outErr != nil {
			plans, outErr = service.DescribeServiceUsagePlan(ctx, serviceId)
		}
		if outErr != nil {
			return nil, outErr
		}
	}

	return plans, nil
}

func testAccCheckAPIGatewayUsagePlanAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAPIGatewayUsagePlanAttachmentResourceName {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		plans, err := base(ctx, rs)
		if err != nil {
			return err
		}

		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		var (
			usagePlanId = ids[0]
			environment = ids[2]
			bindType    = ids[3]
			apiId       = ids[4]
			//accessKeysStr = ids[5]
		)

		for _, plan := range plans {
			if *plan.UsagePlanId == usagePlanId && *plan.Environment == environment {
				if bindType == svcapigateway.API_GATEWAY_TYPE_API {
					if plan.ApiId != nil && *plan.ApiId == apiId {
						return fmt.Errorf("attachment  %s still exist on server", rs.Primary.ID)
					}
				} else {
					return fmt.Errorf("attachment  %s still exist on server", rs.Primary.ID)
				}
			}
		}

		return nil
	}
	return nil
}

func testAccCheckAPIGatewayUsagePlanAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		plans, err := base(ctx, rs)
		if err != nil {
			return err
		}

		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		var (
			usagePlanId = ids[0]
			environment = ids[2]
			bindType    = ids[3]
			apiId       = ids[4]
		)

		for _, plan := range plans {
			if *plan.UsagePlanId == usagePlanId && *plan.Environment == environment {
				if bindType == svcapigateway.API_GATEWAY_TYPE_API {
					if plan.ApiId != nil && *plan.ApiId == apiId {
						return nil
					}
				} else {
					return nil
				}
			}
		}
		return fmt.Errorf("attachment  %s not exist on server", rs.Primary.ID)
	}
}

const testAccAPIGatewayUsagePlanAttachment = `
resource "tencentcloud_api_gateway_usage_plan" "example" {
  usage_plan_name         = "tf_example"
  usage_plan_desc         = "desc."
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example"
  protocol     = "http&https"
  service_desc = "desc."
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "hello_update"
  api_desc              = "my hello api update"
  auth_type             = "SECRET"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "POST"
  request_parameters {
    name          = "email"
    position      = "QUERY"
    type          = "string"
    desc          = "desc."
    default_value = "test@qq.com"
    required      = true
  }
  service_config_type      = "HTTP"
  service_config_timeout   = 10
  service_config_url       = "http://www.tencent.com"
  service_config_path      = "/user"
  service_config_method    = "POST"
  response_type            = "XML"
  response_success_example = "<note>success</note>"
  response_fail_example    = "<note>fail</note>"
  response_error_codes {
    code           = 500
    msg            = "system error"
    desc           = "system error code"
    converted_code = 5000
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_api_key" "example" {
  secret_name = "tf_example"
  status      = "on"
}

resource "tencentcloud_api_gateway_usage_plan_attachment" "example" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.example.id
  service_id    = tencentcloud_api_gateway_service.example.id
  environment   = "release"
  bind_type     = "API"
  api_id        = tencentcloud_api_gateway_api.example.id

  access_key_ids = [
    tencentcloud_api_gateway_api_key.example.id,
  ]
}
`
