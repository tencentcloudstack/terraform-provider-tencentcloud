package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMariadbCreateTmpInstanceResource_basic -v
func TestAccTencentCloudNeedFixMariadbCreateTmpInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCreateTmpInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_create_tmp_instance.create_tmp_instance", "id"),
				),
			},
		},
	})
}

const testAccMariadbCreateTmpInstance = `
resource "tencentcloud_mariadb_create_tmp_instance" "create_tmp_instance" {
  instance_id   = "tdsql-9vqvls95"
  rollback_time = "2023-06-05 01:00:00"
}
`
