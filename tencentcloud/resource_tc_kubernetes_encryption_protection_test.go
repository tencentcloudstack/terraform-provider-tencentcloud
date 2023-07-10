package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTkeEncryptionProtectionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeEncryptionProtection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_encryption_protection.encryption_protection", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_encryption_protection.encryption_protection", "id", "cls-cpsqobnp"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_encryption_protection.encryption_protection", "status", "Opened"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_encryption_protection.encryption_protection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTkeEncryptionProtection = `

resource "tencentcloud_kubernetes_encryption_protection" "encryption_protection" {
  cluster_id = "cls-cpsqobnp"
  k_m_s_configuration {
		key_id = "my_key_id"
		kms_region = "ap-guangzhou"

  }
}

`
