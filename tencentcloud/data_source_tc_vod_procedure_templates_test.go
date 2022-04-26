package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVodProcedureTemplates(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodProcedureTemplates,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_procedure_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.name", "tf-procedure1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.adaptive_dynamic_streaming_task_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.snapshot_by_time_offset_task_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.image_sprite_task_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.0", "3.5s"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.adaptive_dynamic_streaming_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.snapshot_by_time_offset_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.media_process_task.0.image_sprite_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_procedure_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccDataSourceVodProcedureTemplate = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodSnapshotByTimeOffsetTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure1"
  comment = "test"
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
    }
    snapshot_by_time_offset_task_list {
      definition           = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tencentcloud_vod_image_sprite_template.foo.id
    }
  }
}
`

const testAccVodProcedureTemplates = testAccDataSourceVodProcedureTemplate + `
data "tencentcloud_vod_procedure_templates" "foo" {
  type = "Custom"
  name = tencentcloud_vod_procedure_template.foo.id
}
`
