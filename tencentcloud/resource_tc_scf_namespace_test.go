package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

func TestAccTencentCloudScfNamespace_basic(t *testing.T) {
	var nsId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfNamespaceDestroy(&nsId),
		Steps: []resource.TestStep{
			{
				Config: testAccScfNamespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfNamespaceExists("tencentcloud_scf_namespace.foo", &nsId),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "namespace", "ci-test-scf"),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "description", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "type"),
				),
			},
			{
				ResourceName:      "tencentcloud_scf_namespace.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudScfNamespace_desc(t *testing.T) {
	var nsId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfNamespaceDestroy(&nsId),
		Steps: []resource.TestStep{
			{
				Config: testAccScfNamespaceDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfNamespaceExists("tencentcloud_scf_namespace.foo", &nsId),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "namespace", "ci-test-scf"),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "description", "test1"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "type"),
				),
			},
			{
				Config: testAccScfNamespaceDescUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfNamespaceExists("tencentcloud_scf_namespace.foo", &nsId),
					resource.TestCheckResourceAttr("tencentcloud_scf_namespace.foo", "description", "test2"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_namespace.foo", "modify_time"),
				),
			},
		},
	})
}

func testAccCheckScfNamespaceExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errors.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no scf namespace id is set")
		}

		service := ScfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		fn, err := service.DescribeNamespace(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if fn == nil {
			return errors.Errorf("scf namespace not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckScfNamespaceDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := ScfService{client: client}

		namespace, err := service.DescribeNamespace(context.TODO(), *id)
		if err != nil {
			return err
		}

		if namespace != nil {
			return fmt.Errorf("scf namespace still exists")
		}

		return nil
	}
}

const testAccScfNamespace = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}
`

const testAccScfNamespaceDesc = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test1"
}
`

const testAccScfNamespaceDescUpdate = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test2"
}
`
