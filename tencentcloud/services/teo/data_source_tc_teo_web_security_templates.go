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
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    100,
				Description: "Zone ID list. Up to 100 zone IDs per query.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"security_policy_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security policy template list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID that the policy template belongs to.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy template ID.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy template name.",
						},
						"bind_domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain binding information for the policy template.",
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
										Description: "Zone ID that the domain belongs to.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binding status. Valid values: process (binding in progress), online (binding succeeded), fail (binding failed).",
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

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.([]interface{})
		zoneIds := make([]*string, 0, len(zoneIdsSet))
		for i := range zoneIdsSet {
			zoneId := zoneIdsSet[i].(string)
			zoneIds = append(zoneIds, &zoneId)
		}
		paramMap["ZoneIds"] = zoneIds
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

	ids := make([]string, 0, len(respData))
	templatesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, template := range respData {
			templateMap := map[string]interface{}{}

			if template.ZoneId != nil {
				templateMap["zone_id"] = template.ZoneId
			}

			if template.TemplateId != nil {
				templateMap["template_id"] = template.TemplateId
				ids = append(ids, *template.TemplateId)
			}

			if template.TemplateName != nil {
				templateMap["template_name"] = template.TemplateName
			}

			bindDomainsList := make([]map[string]interface{}, 0, len(template.BindDomains))
			if template.BindDomains != nil {
				for _, bindDomain := range template.BindDomains {
					bindDomainMap := map[string]interface{}{}

					if bindDomain.Domain != nil {
						bindDomainMap["domain"] = bindDomain.Domain
					}

					if bindDomain.ZoneId != nil {
						bindDomainMap["zone_id"] = bindDomain.ZoneId
					}

					if bindDomain.Status != nil {
						bindDomainMap["status"] = bindDomain.Status
					}

					bindDomainsList = append(bindDomainsList, bindDomainMap)
				}

				templateMap["bind_domains"] = bindDomainsList
			}

			templatesList = append(templatesList, templateMap)
		}

		_ = d.Set("security_policy_templates", templatesList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
