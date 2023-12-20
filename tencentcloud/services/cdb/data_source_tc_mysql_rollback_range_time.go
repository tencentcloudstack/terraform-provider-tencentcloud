package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlRollbackRangeTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlRollbackRangeTimeRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance IDs, the format of a single instance ID is: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"is_remote_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether the clone instance is in the same zone as the source instance, yes: `false`, no: `true`.",
			},

			"backup_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If the clone instance is not in the same region as the source instance, fill in the region where the clone instance is located, for example: ap-guangzhou.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Returned parameter information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Query database error code.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query database error information.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A list of instance IDs. The format of a single instance ID is: cdb-c1nl9rpv. Same as the instance ID displayed in the cloud database console page.",
						},
						"times": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Retrievable time range.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance rollback start time, time format: 2016-10-29 01:06:04.",
									},
									"end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time of instance rollback, time format: 2016-11-02 11:44:47.",
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

func dataSourceTencentCloudMysqlRollbackRangeTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_rollback_range_time.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("is_remote_zone"); ok {
		paramMap["IsRemoteZone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_region"); ok {
		paramMap["BackupRegion"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.InstanceRollbackRangeTime
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlRollbackRangeTimeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, instanceRollbackRangeTime := range items {
			instanceRollbackRangeTimeMap := map[string]interface{}{}

			if instanceRollbackRangeTime.Code != nil {
				instanceRollbackRangeTimeMap["code"] = instanceRollbackRangeTime.Code
			}

			if instanceRollbackRangeTime.Message != nil {
				instanceRollbackRangeTimeMap["message"] = instanceRollbackRangeTime.Message
			}

			if instanceRollbackRangeTime.InstanceId != nil {
				instanceRollbackRangeTimeMap["instance_id"] = instanceRollbackRangeTime.InstanceId
			}

			if instanceRollbackRangeTime.Times != nil {
				timesList := []interface{}{}
				for _, times := range instanceRollbackRangeTime.Times {
					timesMap := map[string]interface{}{}

					if times.Begin != nil {
						timesMap["begin"] = times.Begin
					}

					if times.End != nil {
						timesMap["end"] = times.End
					}

					timesList = append(timesList, timesMap)
				}

				instanceRollbackRangeTimeMap["times"] = timesList
			}

			ids = append(ids, *instanceRollbackRangeTime.InstanceId)
			tmpList = append(tmpList, instanceRollbackRangeTimeMap)
		}

		_ = d.Set("items", tmpList)
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
