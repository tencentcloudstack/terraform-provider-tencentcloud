package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaTtsTemplateResource_basic -v
func TestAccTencentCloudCiMediaTtsTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaTtsTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTtsTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaTtsTemplateExists("tencentcloud_ci_media_tts_template.media_tts_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_tts_template.media_tts_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "name", "tts_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "mode", "Asyc"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "codec", "pcm"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "voice_type", "ruxue"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "volume", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_tts_template.media_tts_template", "speed", "100"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_tts_template.media_tts_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaTtsTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_tts_template" {
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
			return fmt.Errorf("ci media tts template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaTtsTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media tts template %s is not found", re)
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
			return fmt.Errorf("ci media tts template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaTtsTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }
`

const testAccCiMediaTtsTemplate = testAccCiMediaTtsTemplateVar + `

resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
	bucket = var.bucket
	name = "tts_template"
	mode = "Asyc"
	codec = "pcm"
	voice_type = "ruxue"
	volume = "0"
	speed = "100"
  }

`
