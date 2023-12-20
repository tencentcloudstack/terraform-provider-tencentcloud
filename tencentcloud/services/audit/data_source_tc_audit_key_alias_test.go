package audit_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixAuditKeyAliassDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditKeyAliasDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_audit_key_alias.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_key_alias.all", "audit_key_alias_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_key_alias.all", "audit_key_alias_list.0.key_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_key_alias.all", "audit_key_alias_list.0.key_alias"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditKeyAliasDataSource = `
data "tencentcloud_audit_key_alias" "all" {
	region = "ap-hongkong"
}
`
