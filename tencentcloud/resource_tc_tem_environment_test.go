package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTemEnvironmentResource_basic -v
func TestAccTencentCloudTemEnvironmentResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemEnvironment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemEnvironmentExists("tencentcloud_tem_environment.environment"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_environment.environment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_environment.environment", "environment_name", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_tem_environment.environment", "description", "demo for test"),
					resource.TestCheckResourceAttr("tencentcloud_tem_environment.environment", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_environment.environment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemEnvironmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_environment" {
			continue
		}

		res, err := service.DescribeTemEnvironment(ctx, rs.Primary.ID)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.VersionNamespaceNotFound" {
					return nil
				}
			}
			return err
		}

		if res.Result != nil {
			return fmt.Errorf("tem environment %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemEnvironmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTemEnvironment(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res.Result == nil {
			return fmt.Errorf("tem environment %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemEnvironmentVar = `
variable "vpc_id" {
	default = "` + defaultTemVpcId + `"
}

variable "subnet_id" {
	default = "` + defaultTemSubnetId + `"
}
`

const testAccTemEnvironment = testAccTemEnvironmentVar + `

resource "tencentcloud_tem_environment" "environment" {
	environment_name = "demo"
	description      = "demo for test"
	vpc              = var.vpc_id
	subnet_ids       = [var.subnet_id]
	tags = {
		"createdBy" = "terraform"
	}
  }

`
