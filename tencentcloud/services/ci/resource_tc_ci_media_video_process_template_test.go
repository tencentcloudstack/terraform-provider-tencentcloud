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
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_video_process_template
	resource.AddTestSweepers("tencentcloud_ci_media_video_process_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_video_process_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := localci.NewCiService(client)

			response, _, err := client.UseCiClient(tcacctest.DefaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "video_process_template",
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

// go test -i; go test -test.run TestAccTencentCloudCiMediaVideoProcessTemplateResource_basic -v
func TestAccTencentCloudCiMediaVideoProcessTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiMediaVideoProcessTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaVideoProcessTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaVideoProcessTemplateExists("tencentcloud_ci_media_video_process_template.media_video_process_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_video_process_template.media_video_process_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "bucket", tcacctest.DefaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "name", "video_process_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "color_enhance.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "color_enhance.0.enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "ms_sharpen.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_process_template.media_video_process_template", "ms_sharpen.0.enable", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_video_process_template.media_video_process_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaVideoProcessTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_video_process_template" {
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
			return fmt.Errorf("ci media video process template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaVideoProcessTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media video process template %s is not found", re)
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
			return fmt.Errorf("ci media video process template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaVideoProcessTemplateVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
  }
`

const testAccCiMediaVideoProcessTemplate = testAccCiMediaVideoProcessTemplateVar + `

resource "tencentcloud_ci_media_video_process_template" "media_video_process_template" {
	bucket = var.bucket
	name = "video_process_template"
	color_enhance {
		enable = "true"
	}
	ms_sharpen {
		enable = "false"
	}
}

`
