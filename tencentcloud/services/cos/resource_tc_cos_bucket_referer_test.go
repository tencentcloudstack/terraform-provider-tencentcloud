package cos_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCosBucketRefererResource_basic -v
func TestAccTencentCloudCosBucketRefererResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCosBucketRefererDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketReferer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosBucketRefererExists("tencentcloud_cos_bucket_referer.bucket_referer"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_referer.bucket_referer", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "status", "Enabled"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "referer_type", "Black-List"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "empty_refer_configuration", "Allow"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "domain_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "domain_list.0", "127.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "domain_list.1", "terraform.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_cos_bucket_referer.bucket_referer",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCosBucketRefererUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_referer.bucket_referer", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "status", "Disabled"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_referer.bucket_referer", "referer_type", "Black-List"),
				),
			},
		},
	})
}

func testAccCheckCosBucketRefererDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cos_bucket_referer" {
			continue
		}

		res, err := service.DescribeCosBucketRefererById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil && res.Status == "Enabled" {
			return fmt.Errorf("cos referer still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCosBucketRefererExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("cos referer %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		result, err := service.DescribeCosBucketRefererById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("cos referer not found, Id: %v", rs.Primary.ID)
		}

		if result != nil && result.Status == "Disabled" {
			return fmt.Errorf("cos referer is not enabled, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCosBucketRefererVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
}

`

const testAccCosBucketReferer = testAccCosBucketRefererVar + `

resource "tencentcloud_cos_bucket_referer" "bucket_referer" {
	bucket = var.bucket
	status = "Enabled"
	referer_type = "Black-List"
	domain_list = ["127.0.0.1", "terraform.com"]
	empty_refer_configuration = "Allow"
}
`

const testAccCosBucketRefererUp = testAccCosBucketRefererVar + `

resource "tencentcloud_cos_bucket_referer" "bucket_referer" {
	bucket = var.bucket
	status = "Disabled"
	referer_type = "Black-List"
	domain_list = ["127.0.0.1", "terraform.com"]
	empty_refer_configuration = "Allow"
}
`
