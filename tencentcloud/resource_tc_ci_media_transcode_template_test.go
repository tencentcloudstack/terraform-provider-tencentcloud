package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaTranscodeTemplateResource_basic -v
func TestAccTencentCloudNeedFixCiMediaTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaTranscodeTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTranscodeTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaTranscodeTemplateExists("tencentcloud_ci_media_transcode_template.media_transcode_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_transcode_template.media_transcode_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "name", "transcode_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "container.0.format", "mp4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.codec", "H.264"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.width", "1280"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.fps", "30"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.remove", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.profile", "high"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.bitrate", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.preset", "medium"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "video.0.long_short_mode", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "time_interval.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "time_interval.0.start", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "time_interval.0.duration", "60"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.codec", "aac"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.samplerate", "44100"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.channels", "4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.remove", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.keep_two_tracks", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio.0.switch_track", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.adj_dar_method", "scale"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.is_check_reso", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.reso_adj_method", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.is_check_video_bitrate", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.video_bitrate_adj_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.is_check_audio_bitrate", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.audio_bitrate_adj_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.delete_metadata", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.is_hdr2_sdr", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.hls_encrypt.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "trans_config.0.hls_encrypt.0.is_hls_encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.audio_source", "https://"+defaultCiBucket+".cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.mix_mode", "Once"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.replace", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.enable_start_fadein", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.start_fadein_time", "3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.enable_end_fadeout", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.end_fadeout_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.enable_bgm_fade", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_transcode_template.media_transcode_template", "audio_mix.0.effect_config.0.bgm_fade_time", "1.7"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_ci_media_transcode_template.media_transcode_template",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckCiMediaTranscodeTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_transcode_template" {
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
			return fmt.Errorf("ci media transcode template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaTranscodeTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media transcode template %s is not found", re)
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
			return fmt.Errorf("ci media transcode template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaTranscodeTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }
`

const testAccCiMediaTranscodeTemplate = testAccCiMediaTranscodeTemplateVar + `

resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
	bucket = var.bucket
	name = "transcode_template"
	container {
		format = "mp4"
	}
	video {
		codec = "H.264"
		width = "1280"
		fps = "30"
		remove = "false"
		profile = "high"
		bitrate = "1000"
		preset = "medium"
		long_short_mode = "false"
	}
	time_interval {
		start = "0"
		duration = "60"
	}
	audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"
		keep_two_tracks = "false"
		switch_track = "false"
		sample_format = ""
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
		hls_encrypt {
			is_hls_encrypt = "true"
			# uri_key = ""
		}
	}
	audio_mix {
		audio_source = "https://` + defaultCiBucket + `.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
		mix_mode = "Once"
		replace = "true"
		effect_config {
			enable_start_fadein = "true"
			start_fadein_time = "3"
			enable_end_fadeout = "false"
			end_fadeout_time = "0"
			enable_bgm_fade = "true"
			bgm_fade_time = "1.7"
		}
	}
  }

`
