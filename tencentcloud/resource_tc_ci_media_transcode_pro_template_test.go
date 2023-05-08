package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaTranscodeProTemplateResource_basic -v
func TestAccTencentCloudNeedFixCiMediaTranscodeProTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaTranscodeProTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTranscodeProTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaTranscodeProTemplateExists("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "name", "transcode_pro_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "container.0.format", "mxf"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.codec", "xavc"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.profile", "XAVC-HD_422_10bit"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.width", "1920"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.height", "1080"),
					// resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.interlaced", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.fps", "30000/1001"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "video.0.bitrate", "50000"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "audio.#", "1"),
					// resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "audio.0.codec", "pcm_s24le"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "audio.0.remove", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.adj_dar_method", "scale"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.is_check_reso", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.reso_adj_method", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.is_check_video_bitrate", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.video_bitrate_adj_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.is_check_audio_bitrate", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.audio_bitrate_adj_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.delete_metadata", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "trans_config.0.is_hdr2_sdr", "false"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckCiMediaTranscodeProTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_transcode_pro_template" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		templateId := idSplit[1]

		res, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("ci media transcode pro template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaTranscodeProTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media transcode pro template %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		templateId := idSplit[1]

		result, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("ci media transcode pro template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaTranscodeProTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }
`

const testAccCiMediaTranscodeProTemplate = testAccCiMediaTranscodeProTemplateVar + `

resource "tencentcloud_ci_media_transcode_pro_template" "media_transcode_pro_template" {
	bucket = var.bucket
	name = "transcode_pro_template"
	container {
		format = "mxf"
	}
	video {
		codec = "xavc"
		profile = "XAVC-HD_422_10bit"
		width = "1920"
		height = "1080"
	  	interlaced = "true"
		fps = "30000/1001"
		bitrate = "50000"
	}
	time_interval {
		start = ""
		duration = ""
	}
	audio {
		codec = "pcm_s24le"
		remove = "true"
	}
	trans_config {
		adj_dar_method = "scale"
		is_check_reso = "false"
		reso_adj_method = "1"
		is_check_video_bitrate = "false"
		video_bitrate_adj_method = "0"
		is_check_audio_bitrate = "false"
		audio_bitrate_adj_method = "0"
		delete_metadata = "false"
		is_hdr2_sdr = "false"
	}
}

`
