package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCamPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamPolicyExists("tencentcloud_cam_policy.policy_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "name", "cam-policy-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy.policy_basic", "document"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "description", "test"),
				),
			}, {
				Config: testAccCamPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamPolicyExists("tencentcloud_cam_policy.policy_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "name", "cam-policy-test2"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy.policy_basic", "document"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_policy.policy_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamPolicyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_policy" {
			continue
		}

		_, err := camService.DescribePolicyById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("cam policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cam policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cam policy id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribePolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccCamPolicy_basic = `
resource "tencentcloud_cam_policy" "policy_basic" {
	name        = "cam-policy-test"
	document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
	description = "test"
}
`

const testAccCamPolicy_update = `
resource "tencentcloud_cam_policy" "policy_basic" {
	name     = "cam-policy-test2"
	document = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]},{\"action\":[\"name/cos:PutObject\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
}
`
