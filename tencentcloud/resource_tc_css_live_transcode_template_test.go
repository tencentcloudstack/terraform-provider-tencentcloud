package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_live_transcode_template", &resource.Sweeper{
		Name: "tencentcloud_css_live_transcode_template",
		F:    testSweepCSSLiveTranscodeTemplate,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_live_transcode_template
func testSweepCSSLiveTranscodeTemplate(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	cssService := CssService{client: cli.(*TencentCloudClient).apiV3Conn}

	temps, err := cssService.DescribeCssLiveTranscodeTemplates(ctx)
	if err != nil {
		return err
	}
	if temps == nil {
		return fmt.Errorf("live transcode template not exists.")
	}

	for _, v := range temps {
		delName := v.TemplateName
		delId := v.TemplateId

		if strings.HasPrefix(*delName, defaultCSSPrefix) {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssLiveTranscodeTemplateById(ctx, delId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] sweeper live transcode template %s:%v failed! reason:[%s]", *delName, *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCSSLiveTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssLiveTranscodeTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveTranscodeTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssLiveTranscodeTemplateExists("tencentcloud_css_live_transcode_template.live_transcode_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_template.live_transcode_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "template_name", "tftest900p"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "acodec", "aac"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "audio_bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "video_bitrate", "100"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "vcodec", "origin"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "description", "This_is_a_tf_test_temp."),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "need_video", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "width", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "need_audio", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "height", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "fps", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "gop", "2"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "rotate", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "profile", "baseline"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "bitrate_to_orig", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "height_to_orig", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "fps_to_orig", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "ai_trans_code", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "adapt_bitrate_percent", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "short_edge_as_height", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "drm_type", "fairplay"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "drm_tracks", "SD"),
				),
			},
			{
				Config: testAccCssLiveTranscodeTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssLiveTranscodeTemplateExists("tencentcloud_css_live_transcode_template.live_transcode_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_template.live_transcode_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "template_name", "tftest900p"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "acodec", "aac"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_template.live_transcode_template", "audio_bitrate"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "video_bitrate", "200"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "vcodec", "h264"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "description", "This_is_a_tf_test_temp_changed."),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "need_video", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "width", "10"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "need_audio", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "height", "10"),
					resource.TestCheckResourceAttr("tencentcloud_css_live_transcode_template.live_transcode_template", "fps", "36"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_template.live_transcode_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCssLiveTranscodeTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_live_transcode_template" {
			continue
		}

		watermark, err := cssService.DescribeCssLiveTranscodeTemplate(ctx, helper.StrToInt64Point(rs.Primary.ID))
		if err != nil {
			return nil
		}

		if watermark != nil {
			return fmt.Errorf("css live transcode template still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssLiveTranscodeTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css live transcode template %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css live transcode template id is not set")
		}

		cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		watermark, err := cssService.DescribeCssLiveTranscodeTemplate(ctx, helper.StrToInt64Point(rs.Primary.ID))
		if err != nil {
			return err
		}

		if watermark == nil {
			return fmt.Errorf("css live transcode template not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCssLiveTranscodeTemplate = `
resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = "tftest900p"
  acodec = "aac"
  audio_bitrate = 128
  video_bitrate = 100
  vcodec = "origin"
  description = "This_is_a_tf_test_temp."
  need_video = 1
  width = 0
  need_audio = 1
  height = 0
  fps = 0
  gop = 2
  rotate = 0
  profile = "baseline"
  bitrate_to_orig = 0
  height_to_orig = 0
  fps_to_orig = 0
  ai_trans_code = 0
  adapt_bitrate_percent = 0
  short_edge_as_height = 0
  drm_type = "fairplay"
  drm_tracks = "SD"
}

`

const testAccCssLiveTranscodeTemplate_update = `
resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = "tftest900p"
  acodec = "aac"
  audio_bitrate = 128
  video_bitrate = 200
  vcodec = "h264"
  description = "This_is_a_tf_test_temp_changed."
  need_video = 0
  width = 10
  need_audio = 0
  height = 10
  fps = 36
  profile = "baseline"
}

`
