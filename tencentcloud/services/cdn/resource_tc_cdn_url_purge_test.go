package cdn_test

import (
	"log"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCdnUrlPurge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnUrlPurgeBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_purge.foo", "task_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_purge.foo", "purge_history.#"),
				),
			},
			{
				PreConfig: func() {
					log.Printf("waiting 10 for next purge")
					time.Sleep(time.Second * 10)
				},
				Config: testAccCdnUrlPurgeBasicUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_purge.foo", "task_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_purge.foo", "purge_history.#"),
				),
			},
		},
	})
}

const testAccCdnUrlPurgeBasic = testAccDomainCosForCDN + `
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "http://keep.${local.domain}/ping",
    "https://keep.${local.domain}/ping"
  ]
  area = "overseas"
}
`

const testAccCdnUrlPurgeBasicUpdate = testAccDomainCosForCDN + `
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "http://keep.${local.domain}/ping",
    "https://keep.${local.domain}/ping"
  ]
  area = "overseas"
  redo = 123456
}
`
