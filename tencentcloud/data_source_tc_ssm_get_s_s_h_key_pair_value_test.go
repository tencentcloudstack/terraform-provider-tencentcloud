package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmGetSSHKeyPairValueDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmGetSSHKeyPairValueDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_get_s_s_h_key_pair_value.get_s_s_h_key_pair_value")),
			},
		},
	})
}

const testAccSsmGetSSHKeyPairValueDataSource = `

data "tencentcloud_ssm_get_s_s_h_key_pair_value" "get_s_s_h_key_pair_value" {
  secret_name = ""
  s_s_h_key_id = ""
            }

`
