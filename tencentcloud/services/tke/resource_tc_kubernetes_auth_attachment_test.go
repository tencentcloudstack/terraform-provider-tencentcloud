package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesAuthAttachResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAuthAttachDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_discovery_anonymous_auth", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "use_tke_default", "true"),
				),
			},
			{
				Config: testAccTkeAuthAttachNonDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_discovery_anonymous_auth", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "jwks_uri", "https://ap-guangzhou-oidc.tke.tencentcs.com/id/7cbe7ca92eba3abc76a17de1/openid/v1/jwks"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "issuer", "https://ap-guangzhou-oidc.tke.tencentcs.com/id/7cbe7ca92eba3abc76a17de1"),
				),
			},
			{
				Config: testAccTkeAuthAttachOidcUpdateOidc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_discovery_anonymous_auth", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "use_tke_default", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_oidc_config", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_create_client_id.#"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_auth_attachment.test_auth_attach", "auto_install_pod_identity_webhook_addon", "true"),
				),
			},
		},
	})
}

const testAccTkeAuthAttachDefault = tcacctest.TkeDataSource + `

resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = local.cluster_id
  auto_create_discovery_anonymous_auth = true
  use_tke_default                      = true
}`

const testAccTkeAuthAttachNonDefault = tcacctest.TkeDataSource + `
resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = local.cluster_id
  auto_create_discovery_anonymous_auth = true
  jwks_uri = "https://ap-guangzhou-oidc.tke.tencentcs.com/id/7cbe7ca92eba3abc76a17de1/openid/v1/jwks"
  issuer = "https://ap-guangzhou-oidc.tke.tencentcs.com/id/7cbe7ca92eba3abc76a17de1"
}`

const testAccTkeAuthAttachOidcUpdateOidc = tcacctest.TkeDataSource + `
resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = local.cluster_id
  auto_create_discovery_anonymous_auth = true
  use_tke_default                      = true
  auto_create_oidc_config = true
  auto_create_client_id = ["xxx"]
  auto_install_pod_identity_webhook_addon=true
}`
