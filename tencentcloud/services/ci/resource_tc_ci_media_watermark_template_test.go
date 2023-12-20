package ci_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localci "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_watermark_template
	resource.AddTestSweepers("tencentcloud_ci_media_watermark_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_watermark_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := localci.NewCiService(client)

			response, _, err := client.UseCiClient(tcacctest.DefaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "watermark_template",
			})
			if err != nil {
				return err
			}

			for _, v := range response.TemplateList {
				err := service.DeleteCiMediaTemplateById(ctx, tcacctest.DefaultCiBucket, v.TemplateId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudCiMediaWatermarkTemplateResource_basic -v
func TestAccTencentCloudCiMediaWatermarkTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiMediaWatermarkTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaWatermarkTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaWatermarkTemplateExists("tencentcloud_ci_media_watermark_template.media_watermark_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_watermark_template.media_watermark_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "bucket", tcacctest.DefaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "name", "watermark_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.type", "Text"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.pos", "TopRight"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.loc_mode", "Absolute"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.dx", "128"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.dy", "128"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.start_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.end_time", "100.5"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.0.font_size", "30"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.0.font_type", "simfang.ttf"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.0.font_color", "0xF0F8F0"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.0.transparency", "30"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_watermark_template.media_watermark_template", "watermark.0.text.0.text", "watermark-content"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_watermark_template.media_watermark_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaWatermarkTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_voice_separate_template" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
			return fmt.Errorf("ci media video separate template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaWatermarkTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media video separate template %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
			return fmt.Errorf("ci media video separate template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaWatermarkTemplateVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
  }
`

const testAccCiMediaWatermarkTemplate = testAccCiMediaWatermarkTemplateVar + `

resource "tencentcloud_ci_media_watermark_template" "media_watermark_template" {
	bucket = var.bucket
	name = "watermark_template"
	watermark {
		type = "Text"
		pos = "TopRight"
		loc_mode = "Absolute"
		dx = "128"
		dy = "128"
		start_time = "0"
		end_time = "100.5"
		text {
			font_size = "30"
			font_type = "simfang.ttf"
			font_color = "0xF0F8F0"
			transparency = "30"
			text = "watermark-content"
		}
	}
}

`
