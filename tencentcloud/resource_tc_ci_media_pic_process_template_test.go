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
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_pic_process_template
	resource.AddTestSweepers("tencentcloud_ci_media_pic_process_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_pic_process_template",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := CiService{client: client}

			response, _, err := service.client.UseCiClient(defaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "pic_process_template",
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

// go test -i; go test -test.run TestAccTencentCloudCiMediaPicProcessTemplateResource_basic -v
func TestAccTencentCloudCiMediaPicProcessTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaPicProcessTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaPicProcessTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaPicProcessTemplateTemplateExists("tencentcloud_ci_media_pic_process_template.media_pic_process_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "name", "pic_process_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "pic_process.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "pic_process.0.is_pic_info", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "pic_process.0.process_rule", "imageMogr2/rotate/90"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_pic_process_template.media_pic_process_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaPicProcessTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_pic_process_template" {
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
			return fmt.Errorf("ci media pic process template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaPicProcessTemplateTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media pic process template %s is not found", re)
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
			return fmt.Errorf("ci media pic process template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaPicProcessTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }

`

const testAccCiMediaPicProcessTemplate = testAccCiMediaPicProcessTemplateVar + `

resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
	bucket = var.bucket
	name = "pic_process_template"
	pic_process {
		is_pic_info = "true"
		process_rule = "imageMogr2/rotate/90"
  
	}
}

`
