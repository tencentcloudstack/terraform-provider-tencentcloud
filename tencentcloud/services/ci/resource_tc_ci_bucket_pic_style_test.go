package ci_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	localci "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_bucket_pic_style
	resource.AddTestSweepers("tencentcloud_ci_bucket_pic_style", &resource.Sweeper{
		Name: "tencentcloud_ci_bucket_pic_style",
		F:    testSweepCiBucketPicStyle,
	})
}

func testSweepCiBucketPicStyle(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := localci.NewCiService(client)

	bucket := tcacctest.DefaultCiBucket
	styleName := tcacctest.DefaultStyleName

	for {
		bucketPicStyle, err := service.DescribeCiBucketPicStyleById(ctx, bucket, styleName)
		if err != nil {
			return nil
		}

		if bucketPicStyle == nil {
			return nil
		}

		err = service.DeleteCiBucketPicStyleById(ctx, bucket, styleName)
		if err != nil {
			return err
		}
	}
}

// go test -i; go test -test.run TestAccTencentCloudCiBucketPicStyleResource_basic -v
func TestAccTencentCloudCiBucketPicStyleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_bucket_pic_style" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci bucket pic style %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
	default = "` + tcacctest.DefaultCiBucket + `"
}
variable "style_name" {
	default = "` + tcacctest.DefaultStyleName + `"
}
`

const testAccCiBucketPicStyle = testAccCiBucketPicStyleVar + `

resource "tencentcloud_ci_bucket_pic_style" "bucket_pic_style" {
	bucket     = var.bucket
	style_name = var.style_name
	style_body = "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"
  }

`
