package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudPtsFileResource_basic -v
func TestAccTencentCloudPtsFileResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPtsFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsFile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPtsFileExists("tencentcloud_pts_file.file"),
					resource.TestCheckResourceAttr("tencentcloud_pts_file.file", "kind", "3"),
					resource.TestCheckResourceAttr("tencentcloud_pts_file.file", "name", "iac.txt"),
					resource.TestCheckResourceAttr("tencentcloud_pts_file.file", "size", "10799"),
					resource.TestCheckResourceAttr("tencentcloud_pts_file.file", "type", "text/plain"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_file.file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPtsFileDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := PtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_pts_file" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		fileId := idSplit[1]

		file, err := service.DescribePtsFile(ctx, projectId, fileId)
		if file != nil {
			return fmt.Errorf("pts file %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckPtsFileExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		fileId := idSplit[1]

		service := PtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		file, err := service.DescribePtsFile(ctx, projectId, fileId)
		if file == nil {
			return fmt.Errorf("pts file %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccPtsFile = testAccPtsProject + `

  
resource "tencentcloud_pts_file" "file" {
	file_id = "file-${tencentcloud_pts_project.project.id}"
	project_id = tencentcloud_pts_project.project.id
	kind = 3
	name = "iac.txt"
	size = 10799
	type = "text/plain"
	# line_count = ""
	# head_lines = ""
	# tail_lines = ""
	# header_in_file = ""
	# header_columns = ""
	# file_infos {
	  # name = ""
	  # size = ""
	  # type = ""
	  # updated_at = ""
	# }
}

`
