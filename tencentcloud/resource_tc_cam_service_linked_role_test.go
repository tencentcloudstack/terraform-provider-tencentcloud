package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCamServiceLinkedRoleResource_basic -v
func TestAccTencentCloudCamServiceLinkedRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamServiceLinkedRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamServiceLinkedRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamServiceLinkedRoleExists("tencentcloud_cam_service_linked_role.service_linked_role"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_service_linked_role.service_linked_role", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_service_linked_role.service_linked_role", "qcs_service_name.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cam_service_linked_role.service_linked_role", "custom_suffix", "x-1"),
					resource.TestCheckResourceAttr("tencentcloud_cam_service_linked_role.service_linked_role", "description", "desc cam"),
					resource.TestCheckResourceAttr("tencentcloud_cam_service_linked_role.service_linked_role", "tags.createdBy", "terraform"),
				),
			},
		},
	})
}

func testAccCheckCamServiceLinkedRoleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_service_linked_role" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK]CAM ServiceLinkedRole id is not set")
		}

		instance, err := camService.DescribeCamServiceLinkedRole(ctx, rs.Primary.ID)
		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK]CAM ServiceLinkedRole still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamServiceLinkedRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK]CAM ServiceLinkedRole %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK]CAM ServiceLinkedRole id is not set")
		}

		roleId := rs.Primary.ID

		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribeCamServiceLinkedRole(ctx, roleId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK] CAM ServiceLinkedRole %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamServiceLinkedRole = `

resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
	qcs_service_name = ["cvm.qcloud.com","ekslog.tke.cloud.tencent.com"]
	custom_suffix = "x-1"
	description = "desc cam"
	tags = {
	  "createdBy" = "terraform"
	}
  }

`
