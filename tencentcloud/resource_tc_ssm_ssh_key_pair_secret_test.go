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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret", "description", "for tf test"),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret", "status", "Disabled"),
				),
			},
			{
				Config: testAccSsmSshKeyPairSecretUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret", "description", "for test"),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret", "status", "Enabled"),
				),
			},
			{
				ResourceName:            "tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clean_ssh_key"},
			},
		},
	})
}

const testAccSsmSshKeyPairSecret = `
data "tencentcloud_kms_keys" "kms" {
  key_state = 1
}
resource "tencentcloud_ssm_ssh_key_pair_secret" "ssh_key_pair_secret" {
  secret_name  = "tf-ssh-key-secret"
  project_id   = 0
  description  = "for tf test"
  kms_key_id   = data.tencentcloud_kms_keys.kms.key_list.0.key_id
  ssh_key_name = "tf_ssh_test"
  status       = "Disabled"
  clean_ssh_key = true
}
`

const testAccSsmSshKeyPairSecretUpdate = `
data "tencentcloud_kms_keys" "kms" {
  key_state = 1
}
resource "tencentcloud_ssm_ssh_key_pair_secret" "ssh_key_pair_secret" {
  secret_name  = "tf-ssh-key-secret"
  project_id   = 0
  description  = "for test"
  kms_key_id   = data.tencentcloud_kms_keys.kms.key_list.0.key_id
  ssh_key_name = "tf_ssh_test"
  status       = "Enabled"
  clean_ssh_key = true
}
`
