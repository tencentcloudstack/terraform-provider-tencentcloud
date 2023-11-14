package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDeployCertificateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDeployCertificateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_deploy_certificate_instance.deploy_certificate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssl_deploy_certificate_instance.deploy_certificate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslDeployCertificateInstance = `

resource "tencentcloud_ssl_deploy_certificate_instance" "deploy_certificate_instance" {
  certificate_id = ""
  instance_id_list = 
  resource_type = ""
  status = 
}

`
