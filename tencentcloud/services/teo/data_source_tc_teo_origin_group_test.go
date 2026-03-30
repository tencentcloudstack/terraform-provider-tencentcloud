package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoOriginGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoOriginGroupDataSourceBasic,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_origin_group.foo"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_origin_group.foo", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "origin_group_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "name"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "type"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "records.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "create_time"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "update_time"),
			),
		}},
	})
}

const testAccTeoOriginGroupDataSourceBasic = `

data "tencentcloud_teo_origin_group" "foo" {
  zone_id = "zone-2xkazzl8yf6k"
}
`

func TestAccTencentCloudTeoOriginGroupDataSource_withRecords(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoOriginGroupDataSourceWithRecords,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_origin_group.foo"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_origin_group.foo", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "records.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "records.0.record"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "records.0.type"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "records.0.record_id"),
			),
		}},
	})
}

const testAccTeoOriginGroupDataSourceWithRecords = `

data "tencentcloud_teo_origin_group" "foo" {
  zone_id = "zone-2xkazzl8yf6k"
}
`

func TestAccTencentCloudTeoOriginGroupDataSource_withReferences(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoOriginGroupDataSourceWithReferences,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_origin_group.foo"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_origin_group.foo", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_group.foo", "references.#"),
			),
		}},
	})
}

const testAccTeoOriginGroupDataSourceWithReferences = `

data "tencentcloud_teo_origin_group" "foo" {
  zone_id = "zone-2xkazzl8yf6k"
}
`
