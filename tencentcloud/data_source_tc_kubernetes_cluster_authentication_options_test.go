package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKubernetesClusterAuthenticationOptionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterAuthenticationOptionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_authentication_options.cluster_authentication_options")),
			},
		},
	})
}

const testAccKubernetesClusterAuthenticationOptionsDataSource = `

data "tencentcloud_kubernetes_cluster_authentication_options" "cluster_authentication_options" {
  cluster_id = ""
      }

`
