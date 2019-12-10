package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudScfNamespaces_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudScfNamespaces,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_namespaces.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespace", "ci-test-scf"),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.0.namespace", regexp.MustCompile(`ci-test-scf`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.type"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfNamespaces_desc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudScfNamespacesDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_namespaces.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "description"),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.namespace"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.type"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudScfNamespaces = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}

data "tencentcloud_scf_namespaces" "foo" {
  namespace = "${tencentcloud_scf_namespace.foo.id}"
}
`

const TestAccDataSourceTencentCloudScfNamespacesDesc = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test"
}

data "tencentcloud_scf_namespaces" "foo" {
  description = "${tencentcloud_scf_namespace.foo.description}"
}
`
