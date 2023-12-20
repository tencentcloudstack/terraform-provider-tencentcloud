package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmChcHostsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcHostsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_hosts.chc_hosts")),
			},
		},
	})
}

const testAccCvmChcHostsDataSource = `

data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
	chc_ids = ["chc-0brmw3wl"]
  filters {
    name = "zone"
    values = ["ap-guangzhou-7"]
  }
}
`
