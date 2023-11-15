package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodSnapshotConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_snapshot_config.snapshot_config", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_snapshot_config.snapshot_config", "period", "hourly"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_snapshot_config.snapshot_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnspodSnapshotConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_snapshot_config.snapshot_config", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_snapshot_config.snapshot_config", "period", "daily"),
				),
			},
		},
	})
}

const testAccDnspodSnapshotConfig = `

resource "tencentcloud_dnspod_snapshot_config" "snapshot_config" {
  domain = "iac-tf.cloud"
  period = "hourly"
}

`

const testAccDnspodSnapshotConfigUp = `

resource "tencentcloud_dnspod_snapshot_config" "snapshot_config" {
  domain = "iac-tf.cloud"
  period = "daily"
}

`
