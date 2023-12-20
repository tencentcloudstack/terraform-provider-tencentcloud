package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmInstancesModificationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstancesModificationDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_instances_modification.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_instances_modification.foo", "instance_type_config_status_list.#"),
				),
			},
		},
	})
}

const testAccCvmInstancesModificationDataSource = tcacctest.DefaultCvmModificationVariable + `
data "tencentcloud_cvm_instances_modification" "foo" {
	instance_ids = [var.cvm_id]
}
`
