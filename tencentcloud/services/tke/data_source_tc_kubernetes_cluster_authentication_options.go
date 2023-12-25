package tke

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kubernetes "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

func DataSourceTencentCloudKubernetesClusterAuthenticationOptions() *schema.Resource {
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
				Description: "ServiceAccount authentication configuration. Note: this field may return `null`, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_tke_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Use TKE default issuer and jwksuri. Note: This field may return `null`, indicating that no valid values can be obtained.",
						},
						"issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "service-account-issuer. Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"jwks_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "service-account-jwks-uri. Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"auto_create_discovery_anonymous_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If it is set to `true`, a RABC rule is automatically created to allow anonymous users to access `/.well-known/openid-configuration` and `/openid/v1/jwks`. Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"latest_operation_state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Result of the last modification. Values: `Updating`, `Success`, `Failed` or `TimeOut`. Note: this field may return `null`, indicating that no valid values can be obtained.",
			},

			"oidc_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "OIDC authentication configurations. Note: This field may return `null`, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_oidc_config": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Creating an identity provider. Note: This field may return `null`, indicating that no valid value can be obtained.",
						},
						"auto_create_client_id": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Creating ClientId of the identity provider. Note: This field may return `null`, indicating that no valid value can be obtained.",
						},
						"auto_install_pod_identity_webhook_addon": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Creating the PodIdentityWebhook component. Note: This field may return `null`, indicating that no valid value can be obtained.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_cluster_authentication_options.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clusterId := d.Get("cluster_id").(string)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		options    *kubernetes.ServiceAccountAuthenticationOptions
		oidcConfig *kubernetes.OIDCConfigAuthenticationOptions
		state      string
		e          error
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		options, state, oidcConfig, e = service.DescribeClusterAuthenticationOptions(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)

	if options != nil {
		serviceAccountAuthenticationOptionsMap := map[string]interface{}{}

		if options.UseTKEDefault != nil {
			serviceAccountAuthenticationOptionsMap["use_tke_default"] = options.UseTKEDefault
		}

		if options.Issuer != nil {
			serviceAccountAuthenticationOptionsMap["issuer"] = options.Issuer
		}

		if options.JWKSURI != nil {
			serviceAccountAuthenticationOptionsMap["jwks_uri"] = options.JWKSURI
		}

		if options.AutoCreateDiscoveryAnonymousAuth != nil {
			serviceAccountAuthenticationOptionsMap["auto_create_discovery_anonymous_auth"] = options.AutoCreateDiscoveryAnonymousAuth
		}
		tmpList = append(tmpList, serviceAccountAuthenticationOptionsMap)
		_ = d.Set("service_accounts", []interface{}{serviceAccountAuthenticationOptionsMap})
	}

	if state != "" {
		_ = d.Set("latest_operation_state", state)
	}

	if oidcConfig != nil {
		oIDCConfigAuthenticationOptionsMap := map[string]interface{}{}

		if oidcConfig.AutoCreateOIDCConfig != nil {
			oIDCConfigAuthenticationOptionsMap["auto_create_oidc_config"] = oidcConfig.AutoCreateOIDCConfig
		}

		if oidcConfig.AutoCreateClientId != nil {
			oIDCConfigAuthenticationOptionsMap["auto_create_client_id"] = oidcConfig.AutoCreateClientId
		}

		if oidcConfig.AutoInstallPodIdentityWebhookAddon != nil {
			oIDCConfigAuthenticationOptionsMap["auto_install_pod_identity_webhook_addon"] = oidcConfig.AutoInstallPodIdentityWebhookAddon
		}
		tmpList = append(tmpList, oIDCConfigAuthenticationOptionsMap)
		_ = d.Set("oidc_config", []interface{}{oIDCConfigAuthenticationOptionsMap})
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
