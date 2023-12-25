package scf_test

import (
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudScfNamespaces_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudScfNamespaces,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_scf_namespaces.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespace", tcacctest.DefaultScfNamespace),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_namespaces.foo", "namespaces.0.namespace", tcacctest.DefaultScfNamespace),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_namespaces.foo", "namespaces.0.type"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfNamespaces_desc(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudScfNamespacesDesc,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_scf_namespaces.foo"),
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
data "tencentcloud_scf_namespaces" "foo" {
  namespace = "` + tcacctest.DefaultScfNamespace + `"
}
`

const TestAccDataSourceTencentCloudScfNamespacesDesc = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-desc-namespace"
  description = "test"
}

data "tencentcloud_scf_namespaces" "foo" {
  description = tencentcloud_scf_namespace.foo.description
}
`
