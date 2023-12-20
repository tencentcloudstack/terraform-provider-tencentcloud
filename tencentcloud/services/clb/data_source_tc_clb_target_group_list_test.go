package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbTargetGroupListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupListDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_clb_target_group_list.target_group_list")),
			},
		},
	})
}

const testAccClbTargetGroupListDataSource = `

data "tencentcloud_clb_target_group_list" "target_group_list" {
  filters {
    name = "TargetGroupName"
    values = ["keep-tgg"]
  }
}

`
