package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixAuditKeyAliassDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditKeyAliasDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_audit_key_alias.all"),
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
