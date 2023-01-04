package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaAnimationTemplateResource_basic -v
func TestAccTencentCloudCiMediaAnimationTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaAnimationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaAnimationTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaAnimationTemplateExists("tencentcloud_ci_media_animation_template.media_animation_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_animation_template.media_animation_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "name", "animation_template_test"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "container.0.format", "gif"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.0.codec", "gif"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.0.width", "1280"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.0.height", ""),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.0.fps", "20"),
					// resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "video.0.animate_only_keep_key_frame", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "time_interval.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "time_interval.0.start", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_animation_template.media_animation_template", "time_interval.0.duration", "60"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_ci_media_animation_template.media_animation_template",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckCiMediaAnimationTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_animation_template" {
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
			return fmt.Errorf("ci media animation template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaAnimationTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media animation template %s is not found", re)
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
			return fmt.Errorf("ci media animation template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaAnimationTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }

`

const testAccCiMediaAnimationTemplate = testAccCiMediaAnimationTemplateVar + `

resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
	bucket = var.bucket
	name = "animation_template_test"
	container {
		  format = "gif"
	}
	video {
		  codec = "gif"
		  width = "1280"
		  height = ""
		  fps = "20"
		  animate_only_keep_key_frame = ""
		  animate_time_interval_of_frame = ""
		  animate_frames_per_second = ""
		  quality = ""
  
	}
	time_interval {
		  start = "0"
		  duration = "60"
  
	}
}

`
