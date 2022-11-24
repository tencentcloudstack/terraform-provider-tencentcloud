package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRumProjectResource_basic -v
func TestAccTencentCloudRumProjectResource_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "desc", "desc"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "enable_url_group", "0"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "instance_id", defaultRumInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "name", "name-2"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "rate", "100"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "type", "web"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "url", "*.iac-tf.com"),
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

const testAccRumProjectVar = `
variable "instance_id" {
  default = "` + defaultRumInstanceId + `"
}
`
const testAccRumProject = testAccRumProjectVar + `

resource "tencentcloud_rum_project" "project" {
    desc             = "desc"
    enable_url_group = 0
    instance_id      = var.instance_id
    name             = "name-2"
    rate             = "100"
    type             = "web"
    url              = "*.iac-tf.com"
}

`
