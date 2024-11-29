package scf

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfCustomDomainCreate,
		Read:   resourceTencentCloudScfCustomDomainRead,
		Update: resourceTencentCloudScfCustomDomainUpdate,
		Delete: resourceTencentCloudScfCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain names, pan-domain names are not supported.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol, value range: HTTP, HTTPS, HTTP&HTTPS.",
			},

			"endpoints_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Routing configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Function namespace.",
						},
						"function_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Function name.",
						},
						"qualifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Function alias or version.",
						},
						"path_match": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path, value specification: /,/*,/xxx,/xxx/a,/xxx/*.",
						},
						"path_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Path rewriting policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Path that needs to be rerouted, value specification: /,/*,/xxx,/xxx/a,/xxx/*.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching rules, value range: WildcardRules wildcard matching, ExactRules exact matching.",
									},
									"rewrite": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Replacement values: such as/, /$.",
									},
								},
							},
						},
					},
				},
			},

			"cert_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Certificate configuration information, required for HTTPS protocol.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL Certificates ID.",
						},
					},
				},
			},

			"waf_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Web Application Firewall Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waf_open": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether the Web Application Firewall is turned on, value range:OPEN, CLOSE.",
						},
						"waf_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Web Application Firewall Instance ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudScfCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_custom_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		domain string
	)
	var (
		request  = scf.NewCreateCustomDomainRequest()
		response = scf.NewCreateCustomDomainResponse()
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoints_config"); ok {
		for _, item := range v.([]interface{}) {
			endpointsConfigMap := item.(map[string]interface{})
			endpointsConf := scf.EndpointsConf{}
			if v, ok := endpointsConfigMap["namespace"]; ok {
				endpointsConf.Namespace = helper.String(v.(string))
			}
			if v, ok := endpointsConfigMap["function_name"]; ok {
				endpointsConf.FunctionName = helper.String(v.(string))
			}
			if v, ok := endpointsConfigMap["qualifier"]; ok {
				endpointsConf.Qualifier = helper.String(v.(string))
			}
			if v, ok := endpointsConfigMap["path_match"]; ok {
				endpointsConf.PathMatch = helper.String(v.(string))
			}
			if v, ok := endpointsConfigMap["path_rewrite"]; ok {
				for _, item := range v.([]interface{}) {
					pathRewriteMap := item.(map[string]interface{})
					pathRewriteRule := scf.PathRewriteRule{}
					if v, ok := pathRewriteMap["path"]; ok {
						pathRewriteRule.Path = helper.String(v.(string))
					}
					if v, ok := pathRewriteMap["type"]; ok {
						pathRewriteRule.Type = helper.String(v.(string))
					}
					if v, ok := pathRewriteMap["rewrite"]; ok {
						pathRewriteRule.Rewrite = helper.String(v.(string))
					}
					endpointsConf.PathRewrite = append(endpointsConf.PathRewrite, &pathRewriteRule)
				}
			}
			request.EndpointsConfig = append(request.EndpointsConfig, &endpointsConf)
		}
	}

	if certConfigMap, ok := helper.InterfacesHeadMap(d, "cert_config"); ok {
		certConf := scf.CertConf{}
		if v, ok := certConfigMap["certificate_id"]; ok {
			certConf.CertificateId = helper.String(v.(string))
		}
		request.CertConfig = &certConf
	}

	if wafConfigMap, ok := helper.InterfacesHeadMap(d, "waf_config"); ok {
		wafConf := scf.WafConf{}
		if v, ok := wafConfigMap["waf_open"]; ok {
			wafConf.WafOpen = helper.String(v.(string))
		}
		if v, ok := wafConfigMap["waf_instance_id"]; ok {
			wafConf.WafInstanceId = helper.String(v.(string))
		}
		request.WafConfig = &wafConf
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().CreateCustomDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf custom domain failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(domain)

	return resourceTencentCloudScfCustomDomainRead(d, meta)
}

func resourceTencentCloudScfCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_custom_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	domain := d.Id()

	_ = d.Set("domain", domain)

	respData, err := service.DescribeScfCustomDomainById(ctx, domain)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `scf_custom_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	endpointsConfigList := make([]map[string]interface{}, 0, len(respData.EndpointsConfig))
	if respData.EndpointsConfig != nil {
		for _, endpointsConfig := range respData.EndpointsConfig {
			endpointsConfigMap := map[string]interface{}{}

			if endpointsConfig.Namespace != nil {
				endpointsConfigMap["namespace"] = endpointsConfig.Namespace
			}

			if endpointsConfig.FunctionName != nil {
				endpointsConfigMap["function_name"] = endpointsConfig.FunctionName
			}

			if endpointsConfig.Qualifier != nil {
				endpointsConfigMap["qualifier"] = endpointsConfig.Qualifier
			}

			if endpointsConfig.PathMatch != nil {
				endpointsConfigMap["path_match"] = endpointsConfig.PathMatch
			}

			pathRewriteList := make([]map[string]interface{}, 0, len(endpointsConfig.PathRewrite))
			if endpointsConfig.PathRewrite != nil {
				for _, pathRewrite := range endpointsConfig.PathRewrite {
					pathRewriteMap := map[string]interface{}{}

					if pathRewrite.Path != nil {
						pathRewriteMap["path"] = pathRewrite.Path
					}

					if pathRewrite.Type != nil {
						pathRewriteMap["type"] = pathRewrite.Type
					}

					if pathRewrite.Rewrite != nil {
						pathRewriteMap["rewrite"] = pathRewrite.Rewrite
					}

					pathRewriteList = append(pathRewriteList, pathRewriteMap)
				}

				endpointsConfigMap["path_rewrite"] = pathRewriteList
			}
			endpointsConfigList = append(endpointsConfigList, endpointsConfigMap)
		}

		_ = d.Set("endpoints_config", endpointsConfigList)
	}

	certConfigMap := map[string]interface{}{}

	if respData.CertConfig != nil {
		if respData.CertConfig.CertificateId != nil {
			certConfigMap["certificate_id"] = respData.CertConfig.CertificateId
		}

		_ = d.Set("cert_config", []interface{}{certConfigMap})
	}

	wafConfigMap := map[string]interface{}{}

	if respData.WafConfig != nil {
		if respData.WafConfig.WafOpen != nil {
			wafConfigMap["waf_open"] = respData.WafConfig.WafOpen
		}

		if respData.WafConfig.WafInstanceId != nil {
			wafConfigMap["waf_instance_id"] = respData.WafConfig.WafInstanceId
		}

		_ = d.Set("waf_config", []interface{}{wafConfigMap})
	}

	return nil
}

func resourceTencentCloudScfCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_custom_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"domain"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	domain := d.Id()

	needChange := false
	mutableArgs := []string{"domain", "protocol", "cert_config", "waf_config", "endpoints_config"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := scf.NewUpdateCustomDomainRequest()

		if v, ok := d.GetOk("domain"); ok {
			request.Domain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("protocol"); ok {
			request.Protocol = helper.String(v.(string))
		}

		if certConfigMap, ok := helper.InterfacesHeadMap(d, "cert_config"); ok {
			certConf := scf.CertConf{}
			if v, ok := certConfigMap["certificate_id"]; ok {
				certConf.CertificateId = helper.String(v.(string))
			}
			request.CertConfig = &certConf
		}

		if wafConfigMap, ok := helper.InterfacesHeadMap(d, "waf_config"); ok {
			wafConf := scf.WafConf{}
			if v, ok := wafConfigMap["waf_open"]; ok {
				wafConf.WafOpen = helper.String(v.(string))
			}
			if v, ok := wafConfigMap["waf_instance_id"]; ok {
				wafConf.WafInstanceId = helper.String(v.(string))
			}
			request.WafConfig = &wafConf
		}

		if v, ok := d.GetOk("endpoints_config"); ok {
			for _, item := range v.([]interface{}) {
				endpointsConfigMap := item.(map[string]interface{})
				endpointsConf := scf.EndpointsConf{}
				if v, ok := endpointsConfigMap["namespace"]; ok {
					endpointsConf.Namespace = helper.String(v.(string))
				}
				if v, ok := endpointsConfigMap["function_name"]; ok {
					endpointsConf.FunctionName = helper.String(v.(string))
				}
				if v, ok := endpointsConfigMap["qualifier"]; ok {
					endpointsConf.Qualifier = helper.String(v.(string))
				}
				if v, ok := endpointsConfigMap["path_match"]; ok {
					endpointsConf.PathMatch = helper.String(v.(string))
				}
				if v, ok := endpointsConfigMap["path_rewrite"]; ok {
					for _, item := range v.([]interface{}) {
						pathRewriteMap := item.(map[string]interface{})
						pathRewriteRule := scf.PathRewriteRule{}
						if v, ok := pathRewriteMap["path"]; ok {
							pathRewriteRule.Path = helper.String(v.(string))
						}
						if v, ok := pathRewriteMap["type"]; ok {
							pathRewriteRule.Type = helper.String(v.(string))
						}
						if v, ok := pathRewriteMap["rewrite"]; ok {
							pathRewriteRule.Rewrite = helper.String(v.(string))
						}
						endpointsConf.PathRewrite = append(endpointsConf.PathRewrite, &pathRewriteRule)
					}
				}
				request.EndpointsConfig = append(request.EndpointsConfig, &endpointsConf)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().UpdateCustomDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update scf custom domain failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = domain
	return resourceTencentCloudScfCustomDomainRead(d, meta)
}

func resourceTencentCloudScfCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_custom_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	domain := d.Id()

	var (
		request  = scf.NewDeleteCustomDomainRequest()
		response = scf.NewDeleteCustomDomainResponse()
	)

	request.Domain = helper.String(domain)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().DeleteCustomDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete scf custom domain failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
