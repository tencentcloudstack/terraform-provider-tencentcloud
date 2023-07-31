package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: fmt.Sprintf(testAccKubernetesClusterAuthenticationOptionsDataSource, defaultTkeClusterId),
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_authentication_options.cluster_authentication_options")),
			},
		},
	})
}

const testAccKubernetesClusterAuthenticationOptionsDataSource = `
variable "env_default_tke_cluster_id" {
  type = string
}

data "tencentcloud_kubernetes_cluster_authentication_options" "cluster_authentication_options" {
  cluster_id = var.env_default_tke_cluster_id != "" ? var.env_default_tke_cluster_id : "%s"
}

`
