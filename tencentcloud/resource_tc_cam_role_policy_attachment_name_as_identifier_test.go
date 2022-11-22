package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCamRolePolicyAttachmentNameAsIdentifier_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentNameAsIdentifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachmentNameAsIdentifier_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRolePolicyAttachmentNameAsIdentifierExists("tencentcloud_cam_role_policy_attachment_name_as_identifier.role_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment_name_as_identifier.role_policy_attachment_basic", "role_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment_name_as_identifier.role_policy_attachment_basic", "policy_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role_policy_attachment_name_as_identifier.role_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamRolePolicyAttachmentNameAsIdentifierDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role_policy_attachment_name_as_identifier" {
			continue
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 2 {
			return fmt.Errorf("RolePolicyAttachmentId is invalid!")
		}
		roleName, policyName := items[0], items[1]
		params := make(map[string]interface{})
		params["policy_name"] = policyName
		instance, err := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)

		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Desctroy] check: CAM role policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRolePolicyAttachmentNameAsIdentifierExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 2 {
			return fmt.Errorf("RolePolicyAttachmentId is invalid!")
		}
		roleName, policyName := items[0], items[1]
		params := make(map[string]interface{})
		params["policy_name"] = policyName
		instance, err := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)

		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

//need to add policy resource definition
func testAccCamRolePolicyAttachmentNameAsIdentifier_basic() string {
	return defaultCamVariables + `
resource "tencentcloud_cam_role_policy_attachment_name_as_identifier" "role_policy_attachment_basic" {
  role_name   = var.cam_role_basic
  policy_name = var.cam_policy_basic
}
`
}
