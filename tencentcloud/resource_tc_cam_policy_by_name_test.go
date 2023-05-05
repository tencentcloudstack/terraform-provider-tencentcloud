package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCamPolicyByNameResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamPolicyByNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamPolicyByName_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamPolicyByNameExists("tencentcloud_cam_policy_by_name.policy_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy_by_name.policy_basic", "name", "cam_policy_name_as_identifier_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy_by_name.policy_basic", "document"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy_by_name.policy_basic", "description", "test"),
				),
			}, {
				Config: testAccCamPolicyByName_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamPolicyByNameExists("tencentcloud_cam_policy_by_name.policy_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_policy_by_name.policy_basic", "name", "cam_policy_name_as_identifier_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_policy_by_name.policy_basic", "document"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_policy_by_name.policy_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamPolicyByNameDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_policy_by_name" {
			continue
		}

		params := make(map[string]interface{})
		params["name"] = rs.Primary.ID
		instances, err := camService.DescribePoliciesByFilter(ctx, params)
		if err == nil && len(instances) != 0 {
			return fmt.Errorf("[CHECK][CAM policy][Desctroy] check: CAM policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamPolicyByNameExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		params := make(map[string]interface{})
		params["name"] = rs.Primary.ID
		instances, err := camService.DescribePoliciesByFilter(ctx, params)
		if err != nil {
			return err
		}
		if len(instances) == 0 {
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamPolicyByName_basic = `
resource "tencentcloud_cam_policy_by_name" "policy_basic" {
  name        = "cam_policy_name_as_identifier_test"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"cos:*\"],\"resource\":[\"*\"],\"effect\":\"allow\"},{\"effect\":\"allow\",\"action\":[\"monitor:*\",\"cam:ListUsersForGroup\",\"cam:ListGroups\",\"cam:GetGroup\"],\"resource\":[\"*\"]}]}"
  description = "test"
}
`

const testAccCamPolicyByName_update = `
resource "tencentcloud_cam_policy_by_name" "policy_basic" {
  name     = "cam_policy_name_as_identifier_test"
  document = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"cos:*\"],\"resource\":[\"*\"],\"effect\":\"allow\"},{\"effect\":\"allow\",\"action\":[\"cam:ListUsersForGroup\",\"cam:ListGroups\",\"cam:GetGroup\"],\"resource\":[\"*\"]}]}"
  description = "test2"
}
`
