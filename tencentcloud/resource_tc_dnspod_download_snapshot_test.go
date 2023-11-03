package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDownloadSnapshotResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDownloadSnapshot,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_download_snapshot.download_snapshot", "domain", "iac-tf.cloud"),
				),
			},
		},
	})
}

const testAccDnspodDownloadSnapshot = `

resource "tencentcloud_dnspod_download_snapshot" "download_snapshot" {
  domain = "iac-tf.cloud"
  snapshot_id = "456"
}

`
