package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseKeyPairResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseKeyPair,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_key_pair.key_pair", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_key_pair.key_pair",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudLighthouseKeyPairResource_createByImportPublicKey(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseKeyPairCreateByImportPublicKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_key_pair.key_pair_import", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_key_pair.key_pair_import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseKeyPair = `

resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
}

`
const testAccLighthouseKeyPairCreateByImportPublicKey = `

resource "tencentcloud_lighthouse_key_pair" "key_pair_import" {
  key_name = "key_name_test_import"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0Z/Z2cOr7+pO7hJbdGBlpDqmHT8uXzJ1/mKY1AIAtphKiGPiF/iG+WVREaEJi9eWtuwAyWGZdBnz1/CzeH0WP9k8uihHHba/fLc5OB1BIBLJ7kJmWAWxdk8oKar9MhD+QbsfLniz6RWmAJPpjtEUUcNpwpr6Bn5yC+Rv3xVgY8mc/osPOxZGyok0aqLQXRlT0boYTdUIHE/lf4aTgviUloQzXtSHp6Yn7j0tvK48XSiNDp8NY8gg91fx9oMPv8yzw3TSWiMDueVGhPliiSHE6Zm/mwRgg9giEvSBhBszVCzUeXUeC0PKq+Prl5QHe0xM1xoMDt1cfDePy5555emoP skey-lvs2jhi5"
}

`
