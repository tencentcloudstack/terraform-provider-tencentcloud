package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVodSnapshotByTimeOffsetTemplates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSnapshotByTimeOffsetTemplates,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.name", "tf-snapshot"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.width", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.height", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.resolution_adaptive", "close"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.format", "png"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.fill_type", "white"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccVodSnapshotByTimeOffsetTemplates = testAccVodSnapshotByTimeOffsetTemplate + `
data "tencentcloud_vod_snapshot_by_time_offset_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
}
`
