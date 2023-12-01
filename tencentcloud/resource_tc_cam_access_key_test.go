package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamAccessKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamAccessKey,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_access_key.access_key", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_access_key.access_key", "target_uin", "100033690181"),
				),
			},
			{
				Config: testAccCamAccessKeyUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_access_key.access_key", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_access_key.access_key", "target_uin", "100033690181"),
					resource.TestCheckResourceAttr("tencentcloud_cam_access_key.access_key", "status", "Inactive"),
				),
			},
			{
				Config: testAccCamAccessKeyEncryted,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_access_key.access_key", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_access_key.access_key", "target_uin", "100032767426"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_access_key.access_key", "pgp_key"),
				),
			},
			{
				ResourceName:            "tencentcloud_cam_access_key.access_key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_access_key", "encrypted_secret_access_key", "pgp_key", "key_fingerprint"},
			},
		},
	})
}

const testAccCamAccessKey = `

resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
}

`
const testAccCamAccessKeyUpdate = `

resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  status = "Inactive"
}

`
const testAccCamAccessKeyEncryted = `

resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100032767426
  pgp_key = "keybase:test_iac_tc"
}

`
