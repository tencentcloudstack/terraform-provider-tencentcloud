package tem_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctem "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTemApplicationResource_basic -v
func TestAccTencentCloudTemApplicationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTemApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemApplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemApplicationExists("tencentcloud_tem_application.application"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_application.application", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "application_name", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "description", "demo for test"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "coding_language", "JAVA"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "use_default_image_service", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "repo_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "repo_name", "qcloud/nginx"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "repo_server", "ccr.ccs.tencentyun.com"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application.application", "tags.createdBy", "terraform"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tem_application.application",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTemApplicationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_application" {
			continue
		}

		res, err := service.DescribeTemApplication(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(res.Result.Records) > 0 {
			return fmt.Errorf("tem application %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemApplicationExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTemApplication(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(res.Result.Records) < 1 {
			return fmt.Errorf("tem application %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemApplication = `

resource "tencentcloud_tem_application" "application" {
  application_name = "demo"
  description = "demo for test"
  coding_language = "JAVA"
  use_default_image_service = 0
  repo_type = 2
  repo_name = "qcloud/nginx"
  repo_server = "ccr.ccs.tencentyun.com"
  tags = {
	"createdBy" = "terraform"
  }
}

`
