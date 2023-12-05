package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodRecordList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The domain to which the resolution record belongs.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The domain ID to which the resolution record belongs. If DomainId is provided, the system will ignore the Domain parameter. You can find all Domain and DomainId through the DescribeDomainList interface.",
			},

			"sub_domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Retrieve resolution records based on the host header of the resolution record. Fuzzy matching is used by default. You can set the IsExactSubdomain parameter to true for precise searching.",
			},

			"record_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrieve certain types of resolution records, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, SPF, etc.",
			},

			"record_line": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrieve resolution records for certain line IDs. You can view the allowed line information for the current domain through the DescribeRecordLineList interface.",
			},

			"group_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "When retrieving resolution records under certain groups, pass this group ID. You can obtain the GroupId field through the DescribeRecordGroupList interface.",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search for resolution records by keyword, currently supporting searching host headers and record values.",
			},

			"sort_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, supporting NAME, LINE, TYPE, VALUE, WEIGHT, MX, TTL, UPDATED_ON fields. NAME: The host header of the resolution record LINE: The resolution record line TYPE: The resolution record type VALUE: The resolution record value WEIGHT: The weight MX: MX priority TTL: The resolution record cache time UPDATED_ON: The resolution record update time.",
			},

			"sort_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting method, ascending: ASC, descending: DESC. The default value is ASC.",
			},

			"record_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get the resolution record based on the resolution record value.",
			},

			"record_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get the resolution record based on the resolution record status. The possible values are ENABLE and DISABLE. ENABLE: Normal DISABLE: Paused.",
			},

			"weight_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record weight query interval.",
			},

			"weight_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record weight query interval.",
			},

			"mx_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record MX priority query interval.",
			},

			"mx_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record MX priority query interval.",
			},

			"ttl_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record TTL query interval.",
			},

			"ttl_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record TTL query interval.",
			},

			"updated_at_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The starting point of the resolution record update time query interval.",
			},

			"updated_at_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The endpoint of the resolution record update time query interval.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get the resolution record based on the resolution record remark.",
			},

			"is_exact_sub_domain": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to perform an exact search based on the SubDomain parameter.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"record_count_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Statistics of the number of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subdomain_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of subdomains.",
						},
						"list_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of records returned in the list.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of records.",
						},
					},
				},
			},

			"record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record ID.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record value.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record status, enabled: ENABLE, paused: DISABLE.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host header.",
						},
						"line": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record line.",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record type.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record weight, used for load balancing records. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"monitor_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record monitoring status, normal: OK, alarm: WARN, downtime: DOWN, empty if monitoring is not set or paused.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record remark description.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record cache time.",
						},
						"mx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX value, only available for MX records Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"default_ns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is the default NS record.",
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

func dataSourceTencentCloudDnspodRecordListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnspod_record_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("sub_domain"); ok {
		paramMap["SubDomain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_type"); ok {
		recordTypeSet := v.(*schema.Set).List()
		paramMap["RecordType"] = helper.InterfacesStringsPoint(recordTypeSet)
	}

	if v, ok := d.GetOk("record_line"); ok {
		recordLineSet := v.(*schema.Set).List()
		paramMap["RecordLine"] = helper.InterfacesStringsPoint(recordLineSet)
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupIds := make([]*uint64, 0)
		for _, item := range v.(*schema.Set).List() {
			groupIds = append(groupIds, helper.IntUint64(item.(int)))
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

	if v, ok := d.GetOk("record_value"); ok {
		paramMap["RecordValue"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_status"); ok {
		recordStatusSet := v.(*schema.Set).List()
		paramMap["RecordStatus"] = helper.InterfacesStringsPoint(recordStatusSet)
	}

	if v, ok := d.GetOkExists("weight_begin"); ok {
		paramMap["WeightBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("weight_end"); ok {
		paramMap["WeightEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("mx_begin"); ok {
		paramMap["MXBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("mx_end"); ok {
		paramMap["MXEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ttl_begin"); ok {
		paramMap["TTLBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ttl_end"); ok {
		paramMap["TTLEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("updated_at_begin"); ok {
		paramMap["UpdatedAtBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_end"); ok {
		paramMap["UpdatedAtEnd"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		paramMap["Remark"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_exact_sub_domain"); ok {
		paramMap["IsExactSubDomain"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	var recordList []*dnspod.RecordListItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodRecordListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		recordList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(recordList))
	tmpList := make([]map[string]interface{}, 0, len(recordList))
	if recordList != nil {
		for _, recordListItem := range recordList {
			recordListItemMap := map[string]interface{}{}

			if recordListItem.RecordId != nil {
				recordListItemMap["record_id"] = recordListItem.RecordId
			}

			if recordListItem.Value != nil {
				recordListItemMap["value"] = recordListItem.Value
			}

			if recordListItem.Status != nil {
				recordListItemMap["status"] = recordListItem.Status
			}

			if recordListItem.UpdatedOn != nil {
				recordListItemMap["updated_on"] = recordListItem.UpdatedOn
			}

			if recordListItem.Name != nil {
				recordListItemMap["name"] = recordListItem.Name
			}

			if recordListItem.Line != nil {
				recordListItemMap["line"] = recordListItem.Line
			}

			if recordListItem.LineId != nil {
				recordListItemMap["line_id"] = recordListItem.LineId
			}

			if recordListItem.Type != nil {
				recordListItemMap["type"] = recordListItem.Type
			}

			if recordListItem.Weight != nil {
				recordListItemMap["weight"] = recordListItem.Weight
			}

			if recordListItem.MonitorStatus != nil {
				recordListItemMap["monitor_status"] = recordListItem.MonitorStatus
			}

			if recordListItem.Remark != nil {
				recordListItemMap["remark"] = recordListItem.Remark
			}

			if recordListItem.TTL != nil {
				recordListItemMap["ttl"] = recordListItem.TTL
			}

			if recordListItem.MX != nil {
				recordListItemMap["mx"] = recordListItem.MX
			}

			if recordListItem.DefaultNS != nil {
				recordListItemMap["default_ns"] = recordListItem.DefaultNS
			}

			ids = append(ids, helper.UInt64ToStr(*recordListItem.RecordId))
			tmpList = append(tmpList, recordListItemMap)
		}

		_ = d.Set("record_list", tmpList)
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
