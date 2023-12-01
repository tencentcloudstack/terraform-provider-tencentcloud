/*
Use this data source to query dnspod record list.

Example Usage

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain = "example.com"
  subdomain = "www"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```

Use verbose filter

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain = "example.com"
  subdomain = "www"
  limit = 100
  record_type = "TXT"
  sort_field = "updated_on"
  sort_type = "DESC"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```

*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Description:  "The domain for which DNS records are to be obtained.",
				Optional:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"domain", "domain_id"},
			},
			"domain_id": {
				Description:  "The ID of the domain for which DNS records are to be obtained. If DomainId is passed in, the system will omit the parameter domain.",
				Optional:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"domain", "domain_id"},
			},
			"group_id": {
				Description: "The group ID.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"keyword": {
				Description: "The keyword for searching for DNS records. Host headers and record values are supported.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"limit": {
				Description: "The limit. It defaults to 100 and can be up to 3,000.",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"offset": {
				Description: "The offset. Default value: 0.",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"record_count_info": {
				Computed:    true,
				Description: "Count info of the queried record list.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"list_count": {
						Computed:    true,
						Description: "The count of records returned in the list.",
						Type:        schema.TypeInt,
					},
					"subdomain_count": {
						Computed:    true,
						Description: "The subdomain count.",
						Type:        schema.TypeInt,
					},
					"total_count": {
						Computed:    true,
						Description: "The total record count.",
						Type:        schema.TypeInt,
					},
				}},
				Type: schema.TypeList,
			},
			"record_line": {
				Description: "The split zone name.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"record_line_id": {
				Description: "The split zone ID. If `record_line_id` is passed in, the system will omit the parameter `record_line`.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"record_type": {
				Description: "The type of DNS record, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, or SPF record.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"result": {
				Computed:    true,
				Description: "The record list result.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"line": {
						Computed:    true,
						Description: "The record split zone.",
						Type:        schema.TypeString,
					},
					"line_id": {
						Computed:    true,
						Description: "The split zone ID.",
						Type:        schema.TypeString,
					},
					"monitor_status": {
						Computed:    true,
						Description: "The monitoring status of the record. Valid values: OK (normal), WARN (warning), and DOWN (downtime). It is empty if no monitoring is set or the monitoring is suspended.",
						Type:        schema.TypeString,
					},
					"mx": {
						Computed:    true,
						Description: "The MX value, applicable to the MX record only.\nNote: This field may return null, indicating that no valid values can be obtained.",
						Type:        schema.TypeInt,
					},
					"name": {
						Computed:    true,
						Description: "The host name.",
						Type:        schema.TypeString,
					},
					"record_id": {
						Computed:    true,
						Description: "Record ID.",
						Type:        schema.TypeInt,
					},
					"remark": {
						Computed:    true,
						Description: "The record remarks.",
						Type:        schema.TypeString,
					},
					"status": {
						Computed:    true,
						Description: "The record status. Valid values: ENABLE (enabled), DISABLE (disabled).",
						Type:        schema.TypeString,
					},
					"ttl": {
						Computed:    true,
						Description: "The record cache time.",
						Type:        schema.TypeInt,
					},
					"type": {
						Computed:    true,
						Description: "The record type.",
						Type:        schema.TypeString,
					},
					"updated_on": {
						Computed:    true,
						Description: "The update time.",
						Type:        schema.TypeString,
					},
					"value": {
						Computed:    true,
						Description: "The record value.",
						Type:        schema.TypeString,
					},
					"weight": {
						Computed:    true,
						Description: "The record weight, which is required for round-robin DNS records.",
						Type:        schema.TypeInt,
					},
				}},
				Type: schema.TypeList,
			},
			"result_output_file": {
				Description: "Used for store query result as JSON.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"sort_field": {
				Description: "The sorting field. Available values: name, line, type, value, weight, mx, and ttl,updated_on.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"sort_type": {
				Description: "The sorting type. Valid values: ASC (ascending, default), DESC (descending).",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"subdomain": {
				Description: "The host header of a DNS record. If this parameter is passed in, only the DNS record corresponding to this host header will be returned.",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}
}
func dataSourceTencentCloudDnspodRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("tencentcloud_dnspod_records.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := DnspodService{client}

	request := dnspod.NewDescribeRecordListRequest()
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("subdomain"); ok {
		request.Subdomain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_type"); ok {
		request.RecordType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_line"); ok {
		request.RecordLine = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_line"); ok {
		request.RecordLineId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("keyword"); ok {
		request.Keyword = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sort_field"); ok {
		request.SortField = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sort_type"); ok {
		request.SortType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.IntUint64(v.(int))
	}

	list, info, err := service.DescribeRecordList(ctx, request)

	if err != nil {
		return err
	}

	d.SetId("dnspod_records" + helper.DataResourceIdHash(request.ToJsonString()))

	result := make([]map[string]interface{}, 0, len(list))
	for i := range list {
		record := list[i]
		result = append(result, map[string]interface{}{
			"line":           record.Line,
			"line_id":        record.LineId,
			"monitor_status": record.MonitorStatus,
			"mx":             record.MX,
			"name":           record.Name,
			"record_id":      record.RecordId,
			"remark":         record.Remark,
			"status":         record.Status,
			"ttl":            record.TTL,
			"type":           record.Type,
			"updated_on":     record.UpdatedOn,
			"value":          record.Value,
			"weight":         record.Weight,
		})
	}

	err = helper.SetMapInterfaces(d, "record_count_info", map[string]interface{}{
		"list_count":      info.ListCount,
		"subdomain_count": info.SubdomainCount,
		"total_count":     info.TotalCount,
	})
	if err != nil {
		return err
	}

	err = d.Set("result", result)
	if err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		err = writeToFile(output.(string), map[string]interface{}{
			"record_count_info": info,
			"result":            result,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
