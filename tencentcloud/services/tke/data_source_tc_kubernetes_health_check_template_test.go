package tke

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesHealthCheckTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesHealthCheckTemplateDataSource,
			Check:  resource.ComposeTestCheckFunc(resource.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_health_check_template.kubernetes_health_check_template")),
		}},
	})
}

const testAccKubernetesHealthCheckTemplateDataSource = `

data "tencentcloud_kubernetes_health_check_template" "kubernetes_health_check_template" {
}
`
