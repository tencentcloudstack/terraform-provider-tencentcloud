package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIgtmInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "global_ttl"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "package_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_sub_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "resource_id"),
				),
			},
			{
				Config: testAccIgtmInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "global_ttl"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "package_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "access_sub_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_instance.example", "resource_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_igtm_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIgtmInstance = `
resource "tencentcloud_igtm_instance" "example" {
  domain            = "xtermpro.com"
  access_type       = "CUSTOM"
  global_ttl        = 60
  package_type      = "STANDARD"
  instance_name     = "tf-example"
  access_domain     = "xtermpro.com"
  access_sub_domain = "demo.com"
  remark            = "remark."
  resource_id       = "ins-lnpnnwvwqmr"
}
`

const testAccIgtmInstanceUpdate = `
resource "tencentcloud_igtm_instance" "example" {
  domain            = "xtermpro.com"
  access_type       = "CUSTOM"
  global_ttl        = 80
  package_type      = "STANDARD"
  instance_name     = "tf-example-update"
  access_domain     = "xtermpro.com"
  access_sub_domain = "demo.com"
  remark            = "remark update."
  resource_id       = "ins-lnpnnwvwqmr"
}
`
