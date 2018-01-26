package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudKeyPair_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_key_pair.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_key_pair.foo"),
					resource.TestCheckResourceAttr("tencentcloud_key_pair.foo", "key_name", "from_terraform"),
				),
			},
		},
	})
}

func TestAccTencentCloudKeyPair_pubcliKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_key_pair.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairPublicKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_key_pair.foo"),
					resource.TestCheckResourceAttr("tencentcloud_key_pair.foo", "key_name", "from_terraform_public_key"),
				),
			},
		},
	})
}

func testAccCheckKeyPairDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*TencentCloudClient).commonConn
	var keyId string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_key_pair" {
			continue
		}
		keyId = rs.Primary.ID
	}
	_, _, err := findKeyPairById(client, keyId)
	if err == nil {
		return fmt.Errorf("key pair can still be found after deleted")
	}
	if err.Error() == `tencentcloud_key_pair not found` {
		return nil
	}
	return err
}

const testAccKeyPairBasic = `
resource "tencentcloud_key_pair" "foo" {
  key_name = "from_terraform"
}
`

const testAccKeyPairPublicKey = `
resource "tencentcloud_key_pair" "foo" {
  key_name = "from_terraform_public_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA7TIqdj1BfWk2GoWSwPKT2wkpDTsAGNIgu6SZI8w5/zz02kWMA78WGEv1/I1JCg8BcsagmZeZohSYwgHJXL6QuWDAL6sjGMVAhMQvXRT87apyPX5rzM1MVHlL87btvmc7gLv7OhyDwwpq3Lp62mLxLDVOOu+InTcspQ4eIkHxxuTRo6aXXnimTLezu6CrjWvYxiSkPPZYLG7hwkfcewnsBxp8Tb9i68CLPlenZZph6DUzolNMCGhWmSiqi4BZCaZ8sFrRcU19mWi9gFZYSYSVq3uYZWp4zfwhSJ0dUj592dcAr9/Fuqy7YVUT8KnfR43oyfaeWJLJ4FS8FIpriZAFEw== foo@bar"
}
`
