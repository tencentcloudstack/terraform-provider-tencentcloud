package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDnspodSnapshotResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodSnapshot,
				Check:  resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_snapshot.snapshot", "domain", "iac-tf.cloud")
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodSnapshot = `

resource "tencentcloud_dnspod_snapshot" "snapshot" {
  domain = "iac-tf.cloud"
}

`
