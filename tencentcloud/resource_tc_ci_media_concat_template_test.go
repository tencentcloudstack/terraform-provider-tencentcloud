package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_concat_template
	resource.AddTestSweepers("tencentcloud_ci_media_concat_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_concat_template",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := CiService{client: client}

			response, _, err := service.client.UseCiClient(defaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "concat_templates",
			})
			if err != nil {
				return err
			}

			for _, v := range response.TemplateList {
				err := service.DeleteCiMediaTemplateById(ctx, defaultCiBucket, v.TemplateId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudCiMediaConcatTemplateResource_basic -v
func TestAccTencentCloudCiMediaConcatTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaConcatTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaConcatTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaConcatTemplateTemplateExists("tencentcloud_ci_media_concat_template.media_concat_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_concat_template.media_concat_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "name", "concat_templates"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.concat_fragment.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.concat_fragment.0.url", "https://"+defaultCiBucket+".cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.concat_fragment.0.mode", "Start"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.concat_fragment.1.url", "https://"+defaultCiBucket+".cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.concat_fragment.1.mode", "End"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio.0.codec", "mp3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.video.0.codec", "H.264"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.video.0.width", "1280"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.video.0.bitrate", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.video.0.fps", "25"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.container.0.format", "mp4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.audio_source", "https://"+defaultCiBucket+".cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.mix_mode", "Once"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.replace", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.0.enable_start_fadein", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.0.start_fadein_time", "3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.0.enable_end_fadeout", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.0.enable_bgm_fade", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_concat_template.media_concat_template", "concat_template.0.audio_mix.0.effect_config.0.bgm_fade_time", "1.7"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_concat_template.media_concat_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaConcatTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_concat_template" {
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
			return fmt.Errorf("ci media concat template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaConcatTemplateTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media concat template %s is not found", re)
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
			return fmt.Errorf("ci media concat template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaConcatTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }

`

const testAccCiMediaConcatTemplate = testAccCiMediaConcatTemplateVar + `

resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
	bucket = var.bucket
	name = "concat_templates"
	concat_template {
		concat_fragment {
			url = "https://` + defaultCiBucket + `.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "Start"
		}
		concat_fragment {
			url = "https://` + defaultCiBucket + `.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "End"
		}
		audio {
			codec = "mp3"
			samplerate = ""
			bitrate = ""
			channels = ""
		}
		video {
			codec = "H.264"
			width = "1280"
			height = ""
			bitrate = "1000"
			fps = "25"
			crf = ""
			remove = ""
			rotate = ""
		}
		container {
			format = "mp4"
		}
		audio_mix {
			audio_source = "https://` + defaultCiBucket + `.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
			mix_mode = "Once"
			replace = "true"
			effect_config {
				enable_start_fadein = "true"
				start_fadein_time = "3"
				enable_end_fadeout = "false"
				# end_fadeout_time = "1"
				enable_bgm_fade = "true"
				bgm_fade_time = "1.7"
			}
		}
	}
  }

`
