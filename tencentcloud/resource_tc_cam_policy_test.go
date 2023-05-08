package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cam_policy
	resource.AddTestSweepers("tencentcloud_cam_policy", &resource.Sweeper{
		Name: "tencentcloud_cam_policy",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := CamService{client: client}

			policies, err := service.DescribePoliciesByFilter(ctx, map[string]interface{}{"name": "cam-policy-test"})
			if err != nil {
				return nil
			}
			var (
				request             = cam.NewDeletePolicyRequest()
				presetPolicy uint64 = 2
			)

			for _, v := range policies {
				name := *v.PolicyName
				if *v.Type == presetPolicy || !strings.Contains(name, "cam-policy-test") {
					continue
				}
				request.PolicyId = append(request.PolicyId, v.PolicyId)
			}

			_, err = client.UseCamClient().DeletePolicy(request)

			if err != nil {
				log.Printf("[%s] error, request: %s \nreason: %s", request.GetAction(), request.ToJsonString(), err.Error())
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudCamPolicy_basic(t *testing.T) {
	t.Parallel()
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
			return fmt.Errorf("[CHECK][CAM policy][Desctroy] check: CAM policy still exists: %s", rs.Primary.ID)
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
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribePolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil || instance.Response == nil || instance.Response.PolicyName == nil {
			return fmt.Errorf("[CHECK][CAM policy][Exists] check: CAM policy %s is not exist", rs.Primary.ID)
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
