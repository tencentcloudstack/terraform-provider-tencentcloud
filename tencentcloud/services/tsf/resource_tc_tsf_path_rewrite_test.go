package tsf_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfPathRewriteResource_basic -v
func TestAccTencentCloudTsfPathRewriteResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfPathRewriteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPathRewrite,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfPathRewritekExists("tencentcloud_tsf_path_rewrite.path_rewrite"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_path_rewrite.path_rewrite", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_path_rewrite.path_rewrite", "gateway_group_id", tcacctest.DefaultTsfGroupId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_path_rewrite.path_rewrite", "regex", "/test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_path_rewrite.path_rewrite", "replacement", "/tt"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_path_rewrite.path_rewrite", "blocked", "N"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_path_rewrite.path_rewrite", "order", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_path_rewrite.path_rewrite",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfPathRewriteDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_path_rewrite" {
			continue
		}

		res, err := service.DescribeTsfPathRewriteById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf PathRewrite %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfPathRewritekExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfPathRewriteById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf PathRewrite %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfPathRewriteVar = `
variable "group_id" {
	default = "` + tcacctest.DefaultTsfGroupId + `"
}
`

const testAccTsfPathRewrite = testAccTsfPathRewriteVar + `

resource "tencentcloud_tsf_path_rewrite" "path_rewrite" {
	gateway_group_id = var.group_id
	regex = "/test"
	replacement = "/tt"
	blocked = "N"
	order = 2
}

`
