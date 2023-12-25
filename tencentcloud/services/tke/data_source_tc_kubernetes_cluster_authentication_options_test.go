package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesClusterAuthenticationOptionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterAuthenticationOptionsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_authentication_options.cluster_authentication_options")),
			},
		},
	})
}

const testAccKubernetesClusterAuthenticationOptionsDataSource = `

data "tencentcloud_kubernetes_cluster_authentication_options" "cluster_authentication_options" {
  cluster_id = "cls-kzilgv5m"
}

`
