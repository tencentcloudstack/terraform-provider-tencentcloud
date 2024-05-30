package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKeyPairsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_key_pairs.data_key"), resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.data_key", "key_pair_list.0.public_key"), resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.data_key", "key_pair_list.0.create_time"), resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.data_key", "key_pair_list.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.data_key", "key_pair_list.0.key_id"), resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.data_key", "key_pair_list.0.key_name", "tf_test_key"), resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.data_key", "key_pair_list.0.project_id", "0")),
			},
		},
	})
}

const testAccKeyPairsDataSource_BasicCreate = `

data "tencentcloud_key_pairs" "data_key" {
    key_id = tencentcloud_key_pair.key.id
}
data "tencentcloud_key_pairs" "key_name" {
    key_name = "^${tencentcloud_key_pair.key.key_name}$"
}
resource "tencentcloud_key_pair" "key" {
    key_name = "tf_test_key"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

`
