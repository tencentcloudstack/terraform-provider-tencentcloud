package ssm_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSsmSshKeyPairSecretResource_basic -v
func TestAccTencentCloudSsmSshKeyPairSecretResource_basic(t *testing.T) {
	t.Parallel()
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))
	rSshName := fmt.Sprintf("%d", time.Now().Unix())
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecret, rName, rSshName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "description", "desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example", "status", "Disabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecretUpdate, rName, rSshName),
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
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecretNoId, rSshName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example1", "description", "desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example1", "status", "Disabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSsmSshKeyPairSecretNoIdUpdate, rSshName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example1", "description", "update desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example1", "status", "Enabled"),
				),
			},
			{
				ResourceName:            "tencentcloud_ssm_ssh_key_pair_secret.example1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clean_ssh_key"},
			},
			{
				Config: testAccSsmSshKeyPairSecretNoName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example2", "description", "desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example2", "status", "Disabled"),
				),
			},
			{
				Config: testAccSsmSshKeyPairSecretNoNameUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example2", "description", "update desc."),
					resource.TestCheckResourceAttr("tencentcloud_ssm_ssh_key_pair_secret.example2", "status", "Enabled"),
				),
			},
			{
				ResourceName:            "tencentcloud_ssm_ssh_key_pair_secret.example2",
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
  ssh_key_name  = "tf_ssh_name_%s"
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
  ssh_key_name  = "tf_ssh_name_%s"
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraformUpdate"
  }
}
`

const testAccSsmSshKeyPairSecretNoId = `
resource "tencentcloud_ssm_ssh_key_pair_secret" "example1" {
  secret_name   = "tf-example-ssh-test-no-id"
  project_id    = 0
  description   = "desc."
  ssh_key_name  = "tf_noid_name_%s"
  status        = "Disabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
`

const testAccSsmSshKeyPairSecretNoIdUpdate = `
resource "tencentcloud_ssm_ssh_key_pair_secret" "example1" {
  secret_name   = "tf-example-ssh-test-no-id"
  project_id    = 0
  description   = "update desc."
  ssh_key_name  = "tf_noid_name_%s"
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
`

const testAccSsmSshKeyPairSecretNoName = `
resource "tencentcloud_ssm_ssh_key_pair_secret" "example2" {
  secret_name   = "tf-example-ssh-test-no-id"
  project_id    = 0
  description   = "desc."
  status        = "Disabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
`

const testAccSsmSshKeyPairSecretNoNameUpdate = `
resource "tencentcloud_ssm_ssh_key_pair_secret" "example2" {
  secret_name   = "tf-example-ssh-test-no-id"
  project_id    = 0
  description   = "update desc."
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
`
