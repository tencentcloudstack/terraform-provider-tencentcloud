/*
Use this data source to query detailed information of dnspod domain_list

Example Usage

```hcl

data "tencentcloud_dnspod_domain_list" "domain_list" {
	type = "ALL"
	group_id = [1]
	keyword = ""
	sort_field = "UPDATED_ON"
	sort_type = "DESC"
	status = ["PAUSE"]
	package = [""]
	remark = ""
	updated_at_begin = "2021-05-01 03:00:00"
	updated_at_end = "2024-05-10 20:00:00"
	record_count_begin = 0
	record_count_end = 100
	project_id = -1
}

```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodDomainList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainListRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Get domain names based on domain group type. Available values are ALL, MINE, SHARE, RECENT. ALL: All MINE: My domain names SHARE: Domain names shared with me RECENT: Recently operated domain names.",
			},

			"group_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Get domain names based on domain group id, which can be obtained through the GroupId field in DescribeDomain or DescribeDomainList interface.",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get domain names based on keywords.",
			},

			"sort_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field. Available values are NAME, STATUS, RECORDS, GRADE, UPDATED_ON. NAME: Domain name STATUS: Domain status RECORDS: Number of records GRADE: Package level UPDATED_ON: Update time.",
			},

			"sort_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting type, ascending: ASC, descending: DESC.",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get domain names based on domain status. Available values are ENABLE, LOCK, PAUSE, SPAM. ENABLE: Normal LOCK: Locked PAUSE: Paused SPAM: Banned.",
			},

			"package": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get domain names based on the package, which can be obtained through the Grade field in DescribeDomain or DescribeDomainList interface.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get domain names based on remark information.",
			},

			"updated_at_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The start time of the domain name&amp;#39;s update time to be obtained, such as &amp;#39;2021-05-01 03:00:00&amp;#39;.",
			},

			"updated_at_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The end time of the domain name&amp;#39;s update time to be obtained, such as &amp;#39;2021-05-10 20:00:00&amp;#39;.",
			},

			"record_count_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The start point of the domain name&amp;#39;s record count query range.",
			},

			"record_count_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The end point of the domain name&amp;#39;s record count query range.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"domain_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier assigned to the domain by the system.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original format of the domain.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain status, normal: ENABLE, paused: PAUSE, banned: SPAM.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Default TTL value for domain resolution records.",
						},
						"cname_speedup": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable CNAME acceleration, enabled: ENABLE, disabled: DISABLE.",
						},
						"dns_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS settings status, error: DNSERROR, normal: empty string.",
						},
						"grade": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain package level code.",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Group Id the domain belongs to.",
						},
						"search_engine_push": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable search engine push optimization, YES: YES, NO: NO.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain remark description.",
						},
						"punycode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Punycode encoded domain format.",
						},
						"effective_dns": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Valid DNS assigned to the domain by the system.",
						},
						"grade_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sequence number corresponding to the domain package level.",
						},
						"grade_title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package name.",
						},
						"is_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether it is a paid package.",
						},
						"vip_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Paid package activation time.",
						},
						"vip_end_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Paid package expiration time.",
						},
						"vip_auto_renew": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the domain has VIP auto-renewal enabled, YES: YES, NO: NO, DEFAULT: DEFAULT.",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of records under the domain.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain addition time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain update time.",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain owner account.",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain-related tag list Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Value. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudDnspodDomainListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnspod_domain_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupIds := make([]*int64, 0)
		for _, item := range v.(*schema.Set).List() {
			groupIds = append(groupIds, helper.IntInt64(item.(int)))
		}
		paramMap["GroupId"] = groupIds
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("package"); ok {
		packageSet := v.(*schema.Set).List()
		paramMap["Package"] = helper.InterfacesStringsPoint(packageSet)
	}

	if v, ok := d.GetOk("remark"); ok {
		paramMap["Remark"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_begin"); ok {
		paramMap["UpdatedAtBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_end"); ok {
		paramMap["UpdatedAtEnd"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("record_count_begin"); ok {
		paramMap["RecordCountBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("record_count_end"); ok {
		paramMap["RecordCountEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	var domainList []*dnspod.DomainListItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodDomainListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		domainList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domainList))
	tmpList := make([]map[string]interface{}, 0, len(domainList))
	if domainList != nil {
		for _, domainListItem := range domainList {
			domainListItemMap := map[string]interface{}{}

			if domainListItem.DomainId != nil {
				domainListItemMap["domain_id"] = domainListItem.DomainId
			}

			if domainListItem.Name != nil {
				domainListItemMap["name"] = domainListItem.Name
			}

			if domainListItem.Status != nil {
				domainListItemMap["status"] = domainListItem.Status
			}

			if domainListItem.TTL != nil {
				domainListItemMap["ttl"] = domainListItem.TTL
			}

			if domainListItem.CNAMESpeedup != nil {
				domainListItemMap["cname_speedup"] = domainListItem.CNAMESpeedup
			}

			if domainListItem.DNSStatus != nil {
				domainListItemMap["dns_status"] = domainListItem.DNSStatus
			}

			if domainListItem.Grade != nil {
				domainListItemMap["grade"] = domainListItem.Grade
			}

			if domainListItem.GroupId != nil {
				domainListItemMap["group_id"] = domainListItem.GroupId
			}

			if domainListItem.SearchEnginePush != nil {
				domainListItemMap["search_engine_push"] = domainListItem.SearchEnginePush
			}

			if domainListItem.Remark != nil {
				domainListItemMap["remark"] = domainListItem.Remark
			}

			if domainListItem.Punycode != nil {
				domainListItemMap["punycode"] = domainListItem.Punycode
			}

			if domainListItem.EffectiveDNS != nil {
				domainListItemMap["effective_dns"] = domainListItem.EffectiveDNS
			}

			if domainListItem.GradeLevel != nil {
				domainListItemMap["grade_level"] = domainListItem.GradeLevel
			}

			if domainListItem.GradeTitle != nil {
				domainListItemMap["grade_title"] = domainListItem.GradeTitle
			}

			if domainListItem.IsVip != nil {
				domainListItemMap["is_vip"] = domainListItem.IsVip
			}

			if domainListItem.VipStartAt != nil {
				domainListItemMap["vip_start_at"] = domainListItem.VipStartAt
			}

			if domainListItem.VipEndAt != nil {
				domainListItemMap["vip_end_at"] = domainListItem.VipEndAt
			}

			if domainListItem.VipAutoRenew != nil {
				domainListItemMap["vip_auto_renew"] = domainListItem.VipAutoRenew
			}

			if domainListItem.RecordCount != nil {
				domainListItemMap["record_count"] = domainListItem.RecordCount
			}

			if domainListItem.CreatedOn != nil {
				domainListItemMap["created_on"] = domainListItem.CreatedOn
			}

			if domainListItem.UpdatedOn != nil {
				domainListItemMap["updated_on"] = domainListItem.UpdatedOn
			}

			if domainListItem.Owner != nil {
				domainListItemMap["owner"] = domainListItem.Owner
			}

			ids = append(ids, strconv.FormatUint(*domainListItem.DomainId, 10))
			tmpList = append(tmpList, domainListItemMap)
		}

		_ = d.Set("domain_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
