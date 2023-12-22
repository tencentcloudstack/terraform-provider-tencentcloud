package vod_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_vod_adaptive_dynamic_streaming_template
	resource.AddTestSweepers("tencentcloud_vod_adaptive_dynamic_streaming_template", &resource.Sweeper{
		Name: "tencentcloud_vod_adaptive_dynamic_streaming_template",
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
			templates, e := vodService.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
			if e != nil {
				return nil
			}
			for _, template := range templates {
				ee := vodService.DeleteAdaptiveDynamicStreamingTemplate(ctx, strconv.FormatUint(*template.Definition, 10), uint64(0))
				if ee != nil {
					continue
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudVodAdaptiveDynamicStreamingTemplateResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVodAdaptiveDynamicStreamingTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodAdaptiveDynamicStreamingTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodAdaptiveDynamicStreamingTemplateExists("tencentcloud_vod_adaptive_dynamic_streaming_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "format", "HLS"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "name", "tf-adaptive"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "drm_type", "SimpleAES"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "disable_higher_video_bitrate", "false"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "disable_higher_video_resolution", "false"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "comment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.codec", "libx264"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.fps", "3"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.sample_rate", "32000"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.remove_audio", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.video.0.codec", "libx264"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.video.0.fps", "4"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.video.0.bitrate", "256"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.audio.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.audio.0.bitrate", "256"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.audio.0.sample_rate", "44100"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.1.remove_audio", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "update_time"),
				),
			},
			{
				Config: testAccVodAdaptiveDynamicStreamingTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "name", "tf-adaptive-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "disable_higher_video_bitrate", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "disable_higher_video_resolution", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "comment", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.codec", "libx265"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.fps", "4"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.bitrate", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.width", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.height", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.video.0.fill_type", "stretch"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.bitrate", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.sample_rate", "44100"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.audio.0.audio_channel", "dual"),
					resource.TestCheckResourceAttr("tencentcloud_vod_adaptive_dynamic_streaming_template.foo", "stream_info.0.remove_audio", "false"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_adaptive_dynamic_streaming_template.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sub_app_id"},
			},
		},
	})
}

func testAccCheckVodAdaptiveDynamicStreamingTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_adaptive_dynamic_streaming_template" {
			continue
		}
		var (
			filter = map[string]interface{}{
				"definitions": []string{rs.Primary.ID},
			}
		)

		templates, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return nil
		}
		return fmt.Errorf("vod adaptive dynamic streaming template still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodAdaptiveDynamicStreamingTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod adaptive dynamic streaming template %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod adaptive dynamic streaming template id is not set")
		}
		vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		var (
			filter = map[string]interface{}{
				"definitions": []string{rs.Primary.ID},
			}
		)

		templates, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return fmt.Errorf("vod adaptive dynamic streaming template doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodAdaptiveDynamicStreamingTemplate = `
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
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
  }
}
`

const testAccVodAdaptiveDynamicStreamingTemplateUpdate = `
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive-update"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = true
  disable_higher_video_resolution = true
  comment                         = "test-update"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libfdk_aac"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
  }
}
`
