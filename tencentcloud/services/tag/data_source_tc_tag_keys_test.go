package tag_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTagKeysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTagKeysDataSource,
			Check:  resource.ComposeTestCheckFunc(resource.AccCheckTencentCloudDataSourceID("data.tencentcloud_tag_keys.tag_keys")),
		}},
	})
}

const testAccTagKeysDataSource = `

data "tencentcloud_tag_keys" "tag_keys" {
}
`
