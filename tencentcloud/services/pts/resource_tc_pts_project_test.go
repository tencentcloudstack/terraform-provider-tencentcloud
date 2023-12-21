package pts_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pts"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudPtsProjectResource_basic -v
func TestAccTencentCloudPtsProjectResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPtsProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPtsProjectExists("tencentcloud_pts_project.project"),
					resource.TestCheckResourceAttr("tencentcloud_pts_project.project", "name", "iac-pts-projectName"),
					resource.TestCheckResourceAttr("tencentcloud_pts_project.project", "description", "desc"),
					resource.TestCheckResourceAttr("tencentcloud_pts_project.project", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_project.project", "tags.0.tag_key", "createdBy"),
					resource.TestCheckResourceAttr("tencentcloud_pts_project.project", "tags.0.tag_value", "terraform"),
				),
			},
			{
				ResourceName:            "tencentcloud_pts_project.project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated_at"},
			},
		},
	})
}

func testAccCheckPtsProjectDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_pts_project" {
			continue
		}

		project, err := service.DescribePtsProject(ctx, rs.Primary.ID)
		if project != nil {
			return fmt.Errorf("pts project %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckPtsProjectExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		project, err := service.DescribePtsProject(ctx, rs.Primary.ID)
		if project == nil {
			return fmt.Errorf("pts project %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccPtsProject = `

resource "tencentcloud_pts_project" "project" {
	name = "iac-pts-projectName"
	description = "desc"
	tags {
	  tag_key = "createdBy"
	  tag_value = "terraform"
	}
}

`
