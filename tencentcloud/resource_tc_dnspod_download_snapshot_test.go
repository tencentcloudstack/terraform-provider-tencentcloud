package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDnspodDownloadSnapshotResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDownloadSnapshot,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dnspod_download_snapshot.download_snapshot", "id")),
			},
			{
				ResourceName:      "tencentcloud_dnspod_download_snapshot.download_snapshot",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodDownloadSnapshot = `

resource "tencentcloud_dnspod_download_snapshot" "download_snapshot" {
  domain = "dnspod.cn"
  snapshot_id = "456"
  domain_id = 123
}

`
