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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfGroupResource_basic -v
func TestAccTencentCloudTsfGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfGroupExists("tencentcloud_tsf_group.group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_group.group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "application_id", tcacctest.DefaultTsfApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "namespace_id", tcacctest.DefaultNamespaceId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "group_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "cluster_id", tcacctest.DefaultTsfClustId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "group_desc", "terraform desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "alias", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_group.group", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_group" {
			continue
		}

		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "ResourceNotFound.GroupNotExist" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf group %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf group %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfGroupVar = `
variable "application_id" {
	default = "` + tcacctest.DefaultTsfApplicationId + `"
}
variable "namespace_id" {
	default = "` + tcacctest.DefaultNamespaceId + `"
}
variable "cluster_id" {
	default = "` + tcacctest.DefaultTsfClustId + `"
}
`

const testAccTsfGroup = testAccTsfGroupVar + `

resource "tencentcloud_tsf_group" "group" {
	application_id = var.application_id
	namespace_id = var.namespace_id
	group_name = "terraform-test"
	cluster_id = var.cluster_id
	group_desc = "terraform desc"
	alias = "terraform test"
	tags = {
	  "createdBy" = "terraform"
	}
  }

`
