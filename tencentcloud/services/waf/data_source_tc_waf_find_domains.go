package waf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafFindDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafFindDomainsRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition.",
			},
			"is_waf_domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether access to waf or not.",
			},
			"by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting parameter, eg: FindTime.",
			},
			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     ORDER_DESC,
				Description: "Sorting type, eg: desc, asc.",
			},
			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"appid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User appid.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"ips": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Domain ip.",
						},
						"find_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Find time.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance unique id.",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain unique id.",
						},
						"edition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.",
						},
						"is_waf_domain": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether access to waf or not.",
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

func dataSourceTencentCloudWafFindDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_find_domains.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		list    []*waf.FindAllDomainDetail
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key"); ok {
		paramMap["Key"] = helper.String(v.(string))
	} else {
		paramMap["Key"] = helper.String("")
	}

	if v, ok := d.GetOk("is_waf_domain"); ok {
		paramMap["IsWafDomain"] = helper.String(v.(string))
	} else {
		paramMap["IsWafDomain"] = helper.String("")
	}

	if v, ok := d.GetOk("by"); ok {
		paramMap["By"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafFindDomainsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		list = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(list))
	tmpList := make([]map[string]interface{}, 0, len(list))

	if list != nil {
		for _, findAllDomainDetail := range list {
			findAllDomainDetailMap := map[string]interface{}{}

			if findAllDomainDetail.Appid != nil {
				findAllDomainDetailMap["appid"] = findAllDomainDetail.Appid
			}

			if findAllDomainDetail.Domain != nil {
				findAllDomainDetailMap["domain"] = findAllDomainDetail.Domain
			}

			if findAllDomainDetail.Ips != nil {
				findAllDomainDetailMap["ips"] = findAllDomainDetail.Ips
			}

			if findAllDomainDetail.FindTime != nil {
				findAllDomainDetailMap["find_time"] = findAllDomainDetail.FindTime
			}

			if findAllDomainDetail.InstanceId != nil {
				findAllDomainDetailMap["instance_id"] = findAllDomainDetail.InstanceId
			}

			if findAllDomainDetail.DomainId != nil {
				findAllDomainDetailMap["domain_id"] = findAllDomainDetail.DomainId
			}

			if findAllDomainDetail.Edition != nil {
				findAllDomainDetailMap["edition"] = findAllDomainDetail.Edition
			}

			if findAllDomainDetail.IsWafDomain != nil {
				findAllDomainDetailMap["is_waf_domain"] = findAllDomainDetail.IsWafDomain
			}

			ids = append(ids, *findAllDomainDetail.DomainId)
			tmpList = append(tmpList, findAllDomainDetailMap)
		}

		_ = d.Set("list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
