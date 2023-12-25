package vod_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_vod_procedure_template
	resource.AddTestSweepers("tencentcloud_vod_procedure_template", &resource.Sweeper{
		Name: "tencentcloud_vod_procedure_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			vodService := svcvod.NewVodService(client.GetAPIV3Conn())
			filter := make(map[string]interface{})
			templates, err := vodService.DescribeProcedureTemplatesByFilter(ctx, filter)
			if err != nil {
				return err
			}
			for _, template := range templates {
				ee := vodService.DeleteProcedureTemplate(ctx, *template.Name, uint64(0))
				if ee != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudVodProcedureTemplateResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVodProcedureTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodProcedureTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodProcedureTemplateExists("tencentcloud_vod_procedure_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "name", "tf-procedure0"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "comment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.adaptive_dynamic_streaming_task_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.image_sprite_task_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.0", "3.5s"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.adaptive_dynamic_streaming_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.image_sprite_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "update_time"),
				),
			},
			{
				Config: testAccVodProcedureTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "comment", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.adaptive_dynamic_streaming_task_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.0", "3.5s"),
					resource.TestCheckResourceAttr("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.ext_time_offset_list.1", "4.0s"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.adaptive_dynamic_streaming_task_list.0.definition"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.snapshot_by_time_offset_task_list.0.definition"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_procedure_template.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sub_app_id"},
			},
		},
	})
}

func testAccCheckVodProcedureTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_procedure_template" {
			continue
		}
		var (
			filter = map[string]interface{}{
				"name": []string{rs.Primary.ID},
			}
		)

		templates, err := vodService.DescribeProcedureTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 {
			return nil
		}
		return fmt.Errorf("vod procedure template still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodProcedureTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod procedure template %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod procedure template id is not set")
		}
		vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		var (
			filter = map[string]interface{}{
				"name": []string{rs.Primary.ID},
			}
		)
		templates, err := vodService.DescribeProcedureTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return fmt.Errorf("vod procedure template doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodProcedureTemplate = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodSnapshotByTimeOffsetTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure0"
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

const testAccVodProcedureTemplateUpdate = testAccVodAdaptiveDynamicStreamingTemplate + testAccVodSnapshotByTimeOffsetTemplate + testAccVodImageSpriteTemplate + `
resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure0"
  comment = "test-update"
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
    }
    snapshot_by_time_offset_task_list {
      definition           = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
      ext_time_offset_list = [
        "3.5s",
        "4.0s"
      ]
    }
  }
}
`
