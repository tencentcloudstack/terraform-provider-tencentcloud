package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVodSnapshotByTimeOffsetTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVodSnapshotByTimeOffsetTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSnapshotByTimeOffsetTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodSnapshotByTimeOffsetTemplateExists("tencentcloud_vod_snapshot_by_time_offset_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "name", "tf-snapshot"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "width", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "height", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "resolution_adaptive", "false"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "format", "png"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "comment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "fill_type", "white"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_snapshot_by_time_offset_template.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_snapshot_by_time_offset_template.foo", "update_time"),
				),
			},
			{
				Config: testAccVodSnapshotByTimeOffsetTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "name", "tf-snapshot-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "width", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "height", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "resolution_adaptive", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "format", "jpg"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "comment", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "fill_type", "gauss"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_snapshot_by_time_offset_template.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sub_app_id"},
			},
		},
	})
}

func testAccCheckVodSnapshotByTimeOffsetTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	vodService := VodService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_snapshot_by_time_offset_template" {
			continue
		}

		_, has, err := vodService.DescribeSnapshotByTimeOffsetTemplatesById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("vod snapshot by time offset template still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodSnapshotByTimeOffsetTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod snapshot by time offset template %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod snapshot by time offset template id is not set")
		}
		vodService := VodService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, has, err := vodService.DescribeSnapshotByTimeOffsetTemplatesById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("vod snapshot by time offset template doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodSnapshotByTimeOffsetTemplate = `
resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}
`

const testAccVodSnapshotByTimeOffsetTemplateUpdate = `
resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot-update"
  width               = 129
  height              = 129
  resolution_adaptive = true
  format              = "jpg"
  comment             = "test-update"
  fill_type           = "gauss"
}
`
