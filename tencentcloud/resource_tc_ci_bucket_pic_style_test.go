package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiBucketPicStyleResource_basic -v
func TestAccTencentCloudCiBucketPicStyleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiBucketPicStyleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiBucketPicStyle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiBucketPicStyleExists("tencentcloud_ci_bucket_pic_style.bucket_pic_style"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_bucket_pic_style.bucket_pic_style", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_bucket_pic_style.bucket_pic_style", "style_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_ci_bucket_pic_style.bucket_pic_style", "style_body", "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_bucket_pic_style.bucket_pic_style",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiBucketPicStyleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_bucket_pic_style" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		styleName := idSplit[1]

		res, err := service.DescribeCiBucketPicStyleById(ctx, bucket, styleName)
		if err != nil {
			log.Printf("[ERROR] DescribeCiBucketPicStyleById err: %v\n", err)
			return nil
		}

		if res != nil {
			return fmt.Errorf("ci bucket pic style still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiBucketPicStyleExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci bucket pic style %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		styleName := idSplit[1]

		result, err := service.DescribeCiBucketPicStyleById(ctx, bucket, styleName)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("ci bucket pic style not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiBucketPicStyleVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }

`

const testAccCiBucketPicStyle = testAccCiBucketPicStyleVar + `

resource "tencentcloud_ci_bucket_pic_style" "bucket_pic_style" {
	bucket     = var.bucket
	style_name = "terraform_test"
	style_body = "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"
  }

`
