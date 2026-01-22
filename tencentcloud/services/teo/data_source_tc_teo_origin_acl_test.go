package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoOriginAclDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoOriginAclDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_origin_acl.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_acl.example", "zone_id"),
			),
		}},
	})
}

const testAccTeoOriginAclDataSource = `
data "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
`
