/*
Use this data source to query detailed information of private dns records

Example Usage

```hcl
data "tencentcloud_private_dns_records" "private_dns_record" {
  zone_id = "zone-xxxxxx"
  filters {
	name = "Value"
	values = ["8.8.8.8"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPrivateDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPrivateDnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private zone id: zone-xxxxxx.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter parameters (Value and RecordType filtering are supported).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Parameter values.",
						},
					},
				},
			},

			"record_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parse record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record sid.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private zone id: zone-xxxxxx.",
						},
						"sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdomain name.",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record type, optional record type are: A, AAAA, CNAME, MX, TXT, PTR.",
						},
						"record_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record value.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record cache time, the smaller the value, the faster it takes effect. The value is 1-86400s. The default is 600.",
						},
						"mx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX priority: required if the record type is MX. Value range: 5,10,15,20,30,40,50.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record status.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record weight, value is 1-100.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record creation time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record update time.",
						},
						"extra": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Additional information.",
						},
						"enabled": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Enabled. 0 meaning paused, 1 meaning senabled.",
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

func dataSourceTencentCloudPrivateDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_private_dns_records.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	zoneId := d.Get("zone_id").(string)
	filterList := make([]*privatedns.Filter, 0)

	if v, ok := d.GetOk("filters"); ok {
		filters := v.([]interface{})

		for _, item := range filters {
			filter := privatedns.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			filterList = append(filterList, &filter)
		}
	}

	service := PrivateDnsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var recordSet []*privatedns.PrivateZoneRecord

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePrivateDnsRecordByFilter(ctx, zoneId, filterList)
		if e != nil {
			return retryError(e)
		}
		recordSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(recordSet))
	tmpList := make([]map[string]interface{}, 0, len(recordSet))

	if recordSet != nil {
		for _, privateZoneRecord := range recordSet {
			privateZoneRecordMap := map[string]interface{}{}

			if privateZoneRecord.RecordId != nil {
				privateZoneRecordMap["record_id"] = privateZoneRecord.RecordId
			}

			if privateZoneRecord.ZoneId != nil {
				privateZoneRecordMap["zone_id"] = privateZoneRecord.ZoneId
			}

			if privateZoneRecord.SubDomain != nil {
				privateZoneRecordMap["sub_domain"] = privateZoneRecord.SubDomain
			}

			if privateZoneRecord.RecordType != nil {
				privateZoneRecordMap["record_type"] = privateZoneRecord.RecordType
			}

			if privateZoneRecord.RecordValue != nil {
				privateZoneRecordMap["record_value"] = privateZoneRecord.RecordValue
			}

			if privateZoneRecord.TTL != nil {
				privateZoneRecordMap["ttl"] = privateZoneRecord.TTL
			}

			if privateZoneRecord.MX != nil {
				privateZoneRecordMap["mx"] = privateZoneRecord.MX
			}

			if privateZoneRecord.Status != nil {
				privateZoneRecordMap["status"] = privateZoneRecord.Status
			}

			if privateZoneRecord.Weight != nil {
				privateZoneRecordMap["weight"] = privateZoneRecord.Weight
			}

			if privateZoneRecord.CreatedOn != nil {
				privateZoneRecordMap["created_on"] = privateZoneRecord.CreatedOn
			}

			if privateZoneRecord.UpdatedOn != nil {
				privateZoneRecordMap["updated_on"] = privateZoneRecord.UpdatedOn
			}

			if privateZoneRecord.Extra != nil {
				privateZoneRecordMap["extra"] = privateZoneRecord.Extra
			}

			if privateZoneRecord.Enabled != nil {
				privateZoneRecordMap["enabled"] = privateZoneRecord.Enabled
			}

			ids = append(ids, *privateZoneRecord.RecordId)
			tmpList = append(tmpList, privateZoneRecordMap)
		}

		_ = d.Set("record_set", tmpList)
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
