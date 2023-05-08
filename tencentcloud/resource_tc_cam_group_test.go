package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cam_group
	resource.AddTestSweepers("tencentcloud_cam_group", &resource.Sweeper{
		Name: "tencentcloud_cam_group",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := CamService{client: client}

			groups, err := service.DescribeGroupsByFilter(ctx, nil)
			if err != nil {
				return err
			}
			for _, v := range groups {
				name := *v.GroupName

				if persistResource.MatchString(name) {
					continue
				}

				request := cam.NewDeleteGroupRequest()
				request.GroupId = v.GroupId
				if _, err := client.UseCamClient().DeleteGroup(request); err != nil {
					log.Printf("[%s] error, request: %s \nreason: %s ", request.GetAction(), request.ToJsonString(), err.Error())
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCamGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckCamGroupExists("tencentcloud_cam_group.group_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group.group_basic", "name", "cam-group-test1"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group.group_basic", "remark", "test"),
				),
			}, {
				Config: testAccCamGroup_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupExists("tencentcloud_cam_group.group_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group.group_basic", "name", "cam-group-test2"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group.group_basic", "remark", "test-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_group.group_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_group" {
			continue
		}

		instance, err := camService.DescribeGroupById(ctx, rs.Primary.ID)
		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM group][Destroy] check: CAM group still exists: %s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckCamGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM group][Exists] check: CAM group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM group][Exists] check: CAM group id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribeGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM group][Exists] check: CAM group is not exist")
		}
		return nil
	}
}

const testAccCamGroup_basic = `
resource "tencentcloud_cam_group" "group_basic" {
  name   = "cam-group-test1"
  remark = "test"
}
`

const testAccCamGroup_update = `
resource "tencentcloud_cam_group" "group_basic" {
  name   = "cam-group-test2"
  remark = "test-update"
}
`
