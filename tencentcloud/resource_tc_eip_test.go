package tencentcloud

import (
	"fmt"
	"testing"

	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEip_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBasicWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "gateway_eip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithoutName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_eip.bar"),
					resource.TestCheckNoResourceAttr("tencentcloud_eip.bar", "name"),
					resource.TestCheckResourceAttr("tencentcloud_eip.bar", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.bar", "public_ip"),
				),
			},
		},
	})
}

func testAccCheckEipDestroy(s *terraform.State) error {
	cvmConn := testAccProvider.Meta().(*TencentCloudClient).cvmConn
	var eipId string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip" {
			continue
		}
		eipId = rs.Primary.ID
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, _, err := findEipById(cvmConn, eipId)
		if err == nil {
			err = fmt.Errorf("eip can still be found after deleted")
			return resource.RetryableError(err)
		}
		if err == errEIPNotFound {
			return nil
		}
		return resource.RetryableError(errEIPStillDeleting)
	})
	return err
}

const testAccEipBasicWithName = `
resource "tencentcloud_eip" "foo" {
	name = "gateway_eip"
}
`
const testAccEipBasicWithNewName = `
resource "tencentcloud_eip" "foo" {
	name = "new_name"
}
`

const testAccEipBasicWithoutName = `
resource "tencentcloud_eip" "bar" {
}
`
