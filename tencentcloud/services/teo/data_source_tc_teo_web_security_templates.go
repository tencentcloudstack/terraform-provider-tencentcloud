package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoWebSecurityTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoWebSecurityTemplatesRead,
		Schema: map[string]*schema.Schema{
			"zone_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of zone IDs. A maximum of 100 zones can be queried in a single request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"security_policy_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of policy templates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone ID to which the policy template belongs.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy template ID.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy template.",
						},
						"bind_domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about domains bound to the policy template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name.",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone ID to which the domain belongs.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binding status. valid values:. \n<li>`process`: binding in progress</li>\n<li>`online`: binding succeeded.</li>\n<Li>`fail`: binding failed.</li>.",
									},
								},
							},
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

func dataSourceTencentCloudTeoWebSecurityTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_web_security_templates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsList := []*string{}
		zoneIdsSet := v.(*schema.Set).List()
		for i := range zoneIdsSet {
			zoneIds := zoneIdsSet[i].(string)
			zoneIdsList = append(zoneIdsList, helper.String(zoneIds))
		}

		paramMap["ZoneIds"] = zoneIdsList
	}

	var respData []*teov20220901.SecurityPolicyTemplateInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoWebSecurityTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	securityPolicyTemplatesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, securityPolicyTemplates := range respData {
			securityPolicyTemplatesMap := map[string]interface{}{}
			if securityPolicyTemplates.ZoneId != nil {
				securityPolicyTemplatesMap["zone_id"] = securityPolicyTemplates.ZoneId
			}

			if securityPolicyTemplates.TemplateId != nil {
				securityPolicyTemplatesMap["template_id"] = securityPolicyTemplates.TemplateId
			}

			if securityPolicyTemplates.TemplateName != nil {
				securityPolicyTemplatesMap["template_name"] = securityPolicyTemplates.TemplateName
			}

			bindDomainsList := make([]map[string]interface{}, 0, len(securityPolicyTemplates.BindDomains))
			if securityPolicyTemplates.BindDomains != nil {
				for _, bindDomains := range securityPolicyTemplates.BindDomains {
					bindDomainsMap := map[string]interface{}{}

					if bindDomains.Domain != nil {
						bindDomainsMap["domain"] = bindDomains.Domain
					}

					if bindDomains.ZoneId != nil {
						bindDomainsMap["zone_id"] = bindDomains.ZoneId
					}

					if bindDomains.Status != nil {
						bindDomainsMap["status"] = bindDomains.Status
					}

					bindDomainsList = append(bindDomainsList, bindDomainsMap)
				}

				securityPolicyTemplatesMap["bind_domains"] = bindDomainsList
			}

			securityPolicyTemplatesList = append(securityPolicyTemplatesList, securityPolicyTemplatesMap)
		}

		_ = d.Set("security_policy_templates", securityPolicyTemplatesList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), securityPolicyTemplatesList); e != nil {
			return e
		}
	}

	return nil
}
