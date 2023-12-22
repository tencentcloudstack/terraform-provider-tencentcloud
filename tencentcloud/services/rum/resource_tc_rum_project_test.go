package rum_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcrum "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/rum"
)

// go test -i; go test -test.run TestAccTencentCloudRumProjectResource_basic -v
func TestAccTencentCloudRumProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRumProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumProjectExists("tencentcloud_rum_project.project"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "desc", "desc"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "enable_url_group", "0"),
					resource.TestCheckResourceAttr("tencentcloud_rum_project.project", "instance_id", tcacctest.DefaultRumInstanceId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcrum.NewRumService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
  default = "` + tcacctest.DefaultRumInstanceId + `"
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
