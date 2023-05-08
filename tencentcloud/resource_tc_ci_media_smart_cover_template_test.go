package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_smart_cover_template
	resource.AddTestSweepers("tencentcloud_ci_media_smart_cover_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_smart_cover_template",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := CiService{client: client}

			response, _, err := service.client.UseCiClient(defaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "smart_cover_template",
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

// go test -i; go test -test.run TestAccTencentCloudCiMediaSmartCoverTemplateResource_basic -v
func TestAccTencentCloudCiMediaSmartCoverTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaSmartCoverTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSmartCoverTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaSmartCoverTemplateTemplateExists("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "name", "smart_cover_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.0.format", "jpg"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.0.width", "1280"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.0.height", "960"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.0.count", "10"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "smart_cover.0.delete_duplicates", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_smart_cover_template.media_smart_cover_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaSmartCoverTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_smart_cover_template" {
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
			return fmt.Errorf("ci media smart cover template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaSmartCoverTemplateTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media smart cover template %s is not found", re)
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
			return fmt.Errorf("ci media smart cover template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaSmartCoverTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }

`

const testAccCiMediaSmartCoverTemplate = testAccCiMediaSmartCoverTemplateVar + `

resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
	bucket = var.bucket
	name = "smart_cover_template"
	smart_cover {
		format = "jpg"
		width = "1280"
		height = "960"
		count = "10"
		delete_duplicates = "true"
	}
}

`
