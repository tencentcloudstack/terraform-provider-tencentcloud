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

// go test -i; go test -test.run TestAccTencentCloudTsfNamespaceResource_basic -v
func TestAccTencentCloudTsfNamespaceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfNamespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfNamespaceExists("tencentcloud_tsf_namespace.namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_namespace.namespace", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_namespace.namespace", "namespace_name", "terraform-namespace-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_namespace.namespace", "namespace_desc", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_namespace.namespace", "namespace_type", "DEF"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_namespace.namespace", "is_ha_enable", "0"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_namespace.namespace",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfNamespaceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_namespace" {
			continue
		}

		res, err := service.DescribeTsfNamespaceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf namespace %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfNamespaceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfNamespaceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf namespace %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfNamespace = `

resource "tencentcloud_tsf_namespace" "namespace" {
	namespace_name = "terraform-namespace-name"
	namespace_desc = "terraform-test"
	namespace_type = "DEF"
	is_ha_enable = "0"
}

`
