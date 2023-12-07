package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainSqlFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSqlFiltersRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"filter_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "filter id list.",
			},

			"statuses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "status list.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "sql filter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "task id.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task status, optional value is RUNNING, FINISHED, TERMINATED.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sql type, optional value is SELECT, UPDATE, DELETE, INSERT, REPLACE.",
						},
						"origin_keys": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "origin keys.",
						},
						"origin_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "origin rule.",
						},
						"rejected_sql_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "rejected sql count.",
						},
						"current_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "current concurrency.",
						},
						"max_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "maxmum concurrency.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"current_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "current time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expire time.",
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

func dataSourceTencentCloudDbbrainSqlFiltersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_sql_filters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_ids"); ok {
		filter_idSet := v.(*schema.Set).List()
		tmpList := make([]*int64, 0, len(filter_idSet))
		for i := range filter_idSet {
			filter_id := filter_idSet[i].(int)
			tmpList = append(tmpList, helper.IntInt64(filter_id))
		}
		paramMap["filter_ids"] = tmpList
	}

	if v, ok := d.GetOk("statuses"); ok {
		statuseSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(statuseSet))
		for i := range statuseSet {
			status := statuseSet[i].(string)
			tmpList = append(tmpList, helper.String(status))
		}
		paramMap["statuses"] = tmpList
	}

	dbbrainService := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*dbbrain.SQLFilter
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dbbrainService.DescribeDbbrainSqlFiltersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dbbrain items failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(items))
	itemList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, item := range items {
			itemMap := map[string]interface{}{}
			if item.Id != nil {
				itemMap["id"] = item.Id
			}
			if item.Status != nil {
				itemMap["status"] = item.Status
			}
			if item.SqlType != nil {
				itemMap["sql_type"] = item.SqlType
			}
			if item.OriginKeys != nil {
				itemMap["origin_keys"] = item.OriginKeys
			}
			if item.OriginRule != nil {
				itemMap["origin_rule"] = item.OriginRule
			}
			if item.RejectedSqlCount != nil {
				itemMap["rejected_sql_count"] = item.RejectedSqlCount
			}
			if item.CurrentConcurrency != nil {
				itemMap["current_concurrency"] = item.CurrentConcurrency
			}
			if item.MaxConcurrency != nil {
				itemMap["max_concurrency"] = item.MaxConcurrency
			}
			if item.CreateTime != nil {
				itemMap["create_time"] = item.CreateTime
			}
			if item.CurrentTime != nil {
				itemMap["current_time"] = item.CurrentTime
			}
			if item.ExpireTime != nil {
				itemMap["expire_time"] = item.ExpireTime
			}
			ids = append(ids, helper.Int64ToStr(*item.Id))
			itemList = append(itemList, itemMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", itemList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), itemList); e != nil {
			return e
		}
	}

	return nil
}
