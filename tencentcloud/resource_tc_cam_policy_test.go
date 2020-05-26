package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "name", "cam-policy-test4"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy.policy_basic", "document"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "description", "test"),
				),
			}, {
				Config: testAccCamPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamPolicyExists("tencentcloud_cam_policy.policy_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy.policy_basic", "name", "cam-policy-test4"),
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_policy" {
			continue
		}

		instance, err := camService.DescribePolicyById(ctx, rs.Primary.ID)
		if err == nil && (instance != nil && instance.Response != nil && instance.Response.PolicyName != nil) {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CAM policy][Desctroy] check: CAM policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CAM policy][Exists] check: CAM policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CAM policy][Exists] check: CAM policy id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribePolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil || instance.Response == nil || instance.Response.PolicyName == nil {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CAM policy][Exists] check: CAM policy %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamPolicy_basic = `
resource "tencentcloud_cam_policy" "policy_basic" {
  name        = "cam-policy-test4"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"cos:*\"],\"resource\":[\"*\"],\"effect\":\"allow\"},{\"effect\":\"allow\",\"action\":[\"monitor:*\",\"cam:ListUsersForGroup\",\"cam:ListGroups\",\"cam:GetGroup\"],\"resource\":[\"*\"]}]}"
  description = "test"
}
`

const testAccCamPolicy_update = `
resource "tencentcloud_cam_policy" "policy_basic" {
  name     = "cam-policy-test4"
  document = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"cos:*\"],\"resource\":[\"*\"],\"effect\":\"allow\"},{\"effect\":\"allow\",\"action\":[\"cam:ListUsersForGroup\",\"cam:ListGroups\",\"cam:GetGroup\"],\"resource\":[\"*\"]}]}"
  description = "test2"
}
`
