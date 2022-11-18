package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudRumProject_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumProjectExists("tencentcloud_rum_project.project"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "project_key", "ZEYrYfvaYQ30jRdmPx"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "unique_id", "100027012456"),
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project.project", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRumProjectDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_rum_project" {
			continue
		}

		project, err := service.DescribeRumProject(ctx, rs.Primary.ID)
		if project != nil {
			return fmt.Errorf("rum project %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRumProjectExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		project, err := service.DescribeRumProject(ctx, rs.Primary.ID)
		if project == nil {
			return fmt.Errorf("rum project %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccRumProject = `

resource "tencentcloud_rum_project" "project" {
  name = ""
  instance_id = ""
  rate = ""
  enable_url_group = ""
  type = ""
  repo = ""
  url = ""
  desc = ""
}

`
