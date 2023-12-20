package cdn_test

import (
	"log"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCdnUrlPush(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnUrlPushBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_push.foo", "task_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_push.foo", "push_history.#"),
				),
			},
			{
				PreConfig: func() {
					log.Printf("waiting 10 for next push")
					time.Sleep(time.Second * 10)
				},
				Config: testAccCdnUrlPushBasicUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_push.foo", "task_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_url_push.foo", "push_history.#"),
				),
			},
		},
	})
}

const testAccCdnUrlPushBasic = testAccDomainCosForCDN + `
resource "tencentcloud_cdn_url_push" "foo" {
  urls = [
    "http://keep.${local.domain}/alive",
    "https://keep.${local.domain}/alive"
  ]
  area = "overseas"
}
`

const testAccCdnUrlPushBasicUpdate = testAccDomainCosForCDN + `
resource "tencentcloud_cdn_url_push" "foo" {
  urls = [
    "http://keep.${local.domain}/alive",
    "https://keep.${local.domain}/alive"
  ]
  area = "overseas"
  redo = 111222
}
`
