package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTemApplicationResource_basic -v
func TestAccTencentCloudNeedFixTemApplicationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
}

`
