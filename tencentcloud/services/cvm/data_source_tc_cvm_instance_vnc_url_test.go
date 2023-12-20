package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmInstanceVncUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceVncUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_instance_vnc_url.instance_vnc_url")),
			},
		},
	})
}

const testAccCvmInstanceVncUrlDataSource = tcacctest.DefaultCvmModificationVariable + `

data "tencentcloud_cvm_instance_vnc_url" "instance_vnc_url" {
  instance_id = var.cvm_id
}
`
