package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSsmSshKeyPairSecretResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmSshKeyPairSecret,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSsmSshKeyPairSecret = `

resource "tencentcloud_ssm_ssh_key_pair_secret" "ssh_key_pair_secret" {
  secret_name = ""
  project_id = 
  description = ""
  kms_key_id = ""
  s_s_h_key_name = ""
}

`
