package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSsmSshKeyPairSecretResource_basic -v
func TestAccTencentCloudSsmSshKeyPairSecretResource_basic(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecret, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "description", "desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "status", "Disabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecretUpdate, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "description", "update desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "status", "Enabled"),
				),
			},
			{
				ResourceName:            "tencentcloud_ssm_ssh_key_pair_secret.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clean_ssh_key"},
			},
		},
	})
}

const testAccSsmSshKeyPairSecret = `
resource "tencentcloud_kms_key" "example" {
  alias                = "%s"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_ssh_key_pair_secret" "example" {
  secret_name   = "tf-example-ssh-test"
  project_id    = 0
  description   = "desc."
  kms_key_id    = tencentcloud_kms_key.example.id
  ssh_key_name  = "tf_example_ssh"
  status        = "Disabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
`

const testAccSsmSshKeyPairSecretUpdate = `
resource "tencentcloud_kms_key" "example" {
  alias                = "%s"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_ssh_key_pair_secret" "example" {
  secret_name   = "tf-example-ssh-test"
  project_id    = 0
  description   = "update desc."
  kms_key_id    = tencentcloud_kms_key.example.id
  ssh_key_name  = "tf_example_ssh"
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraformUpdate"
  }
}
`
