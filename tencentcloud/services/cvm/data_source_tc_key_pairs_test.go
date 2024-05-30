package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKeyPairsDataSource -v
func TestAccTencentCloudKeyPairsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists("tencentcloud_key_pair.example"),
					resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.key_id", "key_pair_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.key_id", "key_pair_list.0.key_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.key_id", "key_pair_list.0.key_name", "tf_example"),
					resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.key_id", "key_pair_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.key_id", "key_pair_list.0.public_key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.key_id", "key_pair_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.key_name", "key_pair_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_key_pairs.key_name", "key_pair_list.0.key_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_key_pairs.key_name", "key_pair_list.0.key_name", "tf_example"),
				),
			},
		},
	})
}

const testAccKeyPairDataSource = `
resource "tencentcloud_key_pair" "example" {
  key_name   = "tf_example"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

data "tencentcloud_key_pairs" "key_id" {
  key_id = tencentcloud_key_pair.example.id
}

data "tencentcloud_key_pairs" "key_name" {
  key_name = "^${tencentcloud_key_pair.example.key_name}$"
}
`
