/*
Use this data source to query detailed information of kubernetes cluster_authentication_options

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_authentication_options" "cluster_authentication_options" {
  cluster_id = ""
      }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kubernetes "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesClusterAuthenticationOptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterAuthenticationOptionsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"service_accounts": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "ServiceAccount authentication configurationNote: this field may return `null`, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_t_k_e_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Use TKE default issuer and jwksuriNote: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service-account-issuerNote: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"j_w_k_s_u_r_i": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service-account-jwks-uriNote: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"auto_create_discovery_anonymous_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If it is set to `true`, a RABC rule is automatically created to allow anonymous users to access `/.well-known/openid-configuration` and `/openid/v1/jwks`.Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"latest_operation_state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Result of the last modification. Values: `Updating`, `Success`, `Failed` or `TimeOut`.Note: this field may return `null`, indicating that no valid values can be obtained.",
			},

			"o_i_d_c_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "OIDC authentication configurationsNote: This field may return `null`, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_o_i_d_c_config": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Creating an identity providerNote: This field may return `null`, indicating that no valid value can be obtained.",
						},
						"auto_create_client_id": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Creating ClientId of the identity providerNote: This field may return `null`, indicating that no valid value can be obtained.",
						},
						"auto_install_pod_identity_webhook_addon": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Creating the PodIdentityWebhook componentNote: This field may return `null`, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKubernetesClusterAuthenticationOptionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kubernetes_cluster_authentication_options.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	service := KubernetesService{client: meta.(*TencentCloudClient).apiV3Conn}

	var serviceAccounts []*kubernetes.ServiceAccountAuthenticationOptions

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterAuthenticationOptionsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		serviceAccounts = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(serviceAccounts))
	if serviceAccounts != nil {
		serviceAccountAuthenticationOptionsMap := map[string]interface{}{}

		if serviceAccounts.UseTKEDefault != nil {
			serviceAccountAuthenticationOptionsMap["use_t_k_e_default"] = serviceAccounts.UseTKEDefault
		}

		if serviceAccounts.Issuer != nil {
			serviceAccountAuthenticationOptionsMap["issuer"] = serviceAccounts.Issuer
		}

		if serviceAccounts.JWKSURI != nil {
			serviceAccountAuthenticationOptionsMap["j_w_k_s_u_r_i"] = serviceAccounts.JWKSURI
		}

		if serviceAccounts.AutoCreateDiscoveryAnonymousAuth != nil {
			serviceAccountAuthenticationOptionsMap["auto_create_discovery_anonymous_auth"] = serviceAccounts.AutoCreateDiscoveryAnonymousAuth
		}

		ids = append(ids, *serviceAccounts.ClusterId)
		_ = d.Set("service_accounts", serviceAccountAuthenticationOptionsMap)
	}

	if latestOperationState != nil {
		_ = d.Set("latest_operation_state", latestOperationState)
	}

	if oIDCConfig != nil {
		oIDCConfigAuthenticationOptionsMap := map[string]interface{}{}

		if oIDCConfig.AutoCreateOIDCConfig != nil {
			oIDCConfigAuthenticationOptionsMap["auto_create_o_i_d_c_config"] = oIDCConfig.AutoCreateOIDCConfig
		}

		if oIDCConfig.AutoCreateClientId != nil {
			oIDCConfigAuthenticationOptionsMap["auto_create_client_id"] = oIDCConfig.AutoCreateClientId
		}

		if oIDCConfig.AutoInstallPodIdentityWebhookAddon != nil {
			oIDCConfigAuthenticationOptionsMap["auto_install_pod_identity_webhook_addon"] = oIDCConfig.AutoInstallPodIdentityWebhookAddon
		}

		ids = append(ids, *oIDCConfig.ClusterId)
		_ = d.Set("o_i_d_c_config", oIDCConfigAuthenticationOptionsMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), serviceAccountAuthenticationOptionsMap); e != nil {
			return e
		}
	}
	return nil
}
