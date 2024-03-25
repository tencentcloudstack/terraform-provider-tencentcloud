package vod_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodProcedureTemplate,
				Check: resource.ComposeTestCheckFunc(
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
					resource.TestCheckResourceAttrSet("tencentcloud_vod_procedure_template.foo", "media_process_task.0.transcode_task_list.0.definition"),
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
				ResourceName:      "tencentcloud_vod_procedure_template.foo",
				ImportState:       true,
				ImportStateVerify: true,
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
		id := rs.Primary.ID
		filter := map[string]interface{}{}

		idSplit := strings.Split(id, tccommon.FILED_SP)
		if len(idSplit) == 2 {
			filter["name"] = []string{idSplit[0]}
			filter["sub_appid"] = helper.StrToInt(idSplit[1])
		} else {
			return fmt.Errorf("can not get sub_appid")
		}

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
		id := rs.Primary.ID
		filter := map[string]interface{}{}

		idSplit := strings.Split(id, tccommon.FILED_SP)
		if len(idSplit) == 2 {
			filter["name"] = []string{idSplit[0]}
			filter["sub_appid"] = helper.StrToInt(idSplit[1])
		} else {
			return fmt.Errorf("can not get sub_appid")
		}
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

const testAccVodProcedureTemplate = `
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "procedure-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
    tehd_config {
      type = "TEHD-100"
    }
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p update"
  remove_video = 0
  remove_audio = 0
  video_template {
		codec = "libx264"
		fps = 26
		bitrate = 1000
		resolution_adaptive = "open"
		width = 0
		height = 720
		fill_type = "stretch"
		vcrf = 1
		gop = 250
		preserve_hdr_switch = "OFF"
		codec_tag = "hvc1"

  }
  audio_template {
		codec = "libfdk_aac"
		bitrate = 128
		sample_rate = 44100
		audio_channel = 2

  }
  segment_type = "ts"
}

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure0"
  comment = "test"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tonumber(split("#", tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id)[1])
    }
    snapshot_by_time_offset_task_list {
      definition = tonumber(split("#", tencentcloud_vod_snapshot_by_time_offset_template.foo.id)[1])
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tonumber(split("#", tencentcloud_vod_image_sprite_template.foo.id)[1])
    }
    transcode_task_list {
      definition = tonumber(split("#", tencentcloud_vod_transcode_template.transcode_template.id)[1])
    }
  }
}
`

const testAccVodProcedureTemplateUpdate = `
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "procedure-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
    tehd_config {
      type = "TEHD-100"
    }
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p update"
  remove_video = 0
  remove_audio = 0
  video_template {
		codec = "libx264"
		fps = 26
		bitrate = 1000
		resolution_adaptive = "open"
		width = 0
		height = 720
		fill_type = "stretch"
		vcrf = 1
		gop = 250
		preserve_hdr_switch = "OFF"
		codec_tag = "hvc1"

  }
  audio_template {
		codec = "libfdk_aac"
		bitrate = 128
		sample_rate = 44100
		audio_channel = 2

  }
  segment_type = "ts"
}

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure0"
  comment = "test-update"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  media_process_task {
    adaptive_dynamic_streaming_task_list {
		definition = tonumber(split("#", tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id)[1])
    }
    snapshot_by_time_offset_task_list {
		definition = tonumber(split("#", tencentcloud_vod_snapshot_by_time_offset_template.foo.id)[1])
      ext_time_offset_list = [
        "3.5s",
        "4.0s"
      ]
    }
    transcode_task_list {
      definition = tonumber(split("#", tencentcloud_vod_transcode_template.transcode_template.id)[1])
    }
  }
}
`
