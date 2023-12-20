package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbInstanceByCertIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceByCertIdDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_clb_instance_by_cert_id.instance_by_cert_id")),
			},
		},
	})
}

const testAccClbInstanceByCertIdDataSource = `

data "tencentcloud_clb_instance_by_cert_id" "instance_by_cert_id" {
  cert_ids = ["3a6B5y8v"]
}

`
