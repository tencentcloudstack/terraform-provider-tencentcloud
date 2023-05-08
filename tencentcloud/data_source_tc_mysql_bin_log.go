/*
Use this data source to query detailed information of mysql bin_log

Example Usage

```hcl
data "tencentcloud_mysql_bin_log" "bin_log" {
  instance_id = "cdb-fitq5t9h"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlBinLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBinLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details of binary log files that meet the query conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "binlog log backup file name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup file size, unit: Byte.",
						},
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File storage time, time format: 2016-03-17 02:10:37.",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "download link.",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "download link.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specific log type, possible values are: binlog - binary log.",
						},
						"binlog_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binlog file start time.",
						},
						"binlog_finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "binlog file deadline.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the local binlog file is located.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup task status. Possible values are `SUCCESS`: backup succeeded, `FAILED`: backup failed, `RUNNING`: backup in progress.",
						},
						"remote_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Binlog remote backup details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub_backup_id": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "The ID of the remote backup subtask.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region where remote backup is located.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backup task status. Possible values are `SUCCESS`: backup succeeded, `FAILED`: backup failed, `RUNNING`: backup in progress.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start time of remote backup task.",
									},
									"finish_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time of remote backup task.",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "download link.",
									},
								},
							},
						},
						"cos_storage_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage method, 0-regular storage, 1-archive storage, the default is 0.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
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

func dataSourceTencentCloudMysqlBinLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_bin_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var binLog []*cdb.BinlogInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBinLogByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		binLog = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(binLog))
	tmpList := make([]map[string]interface{}, 0, len(binLog))
	if binLog != nil {
		for _, binlogInfo := range binLog {
			binlogInfoMap := map[string]interface{}{}

			if binlogInfo.Name != nil {
				binlogInfoMap["name"] = binlogInfo.Name
			}

			if binlogInfo.Size != nil {
				binlogInfoMap["size"] = binlogInfo.Size
			}

			if binlogInfo.Date != nil {
				binlogInfoMap["date"] = binlogInfo.Date
			}

			if binlogInfo.IntranetUrl != nil {
				binlogInfoMap["intranet_url"] = binlogInfo.IntranetUrl
			}

			if binlogInfo.InternetUrl != nil {
				binlogInfoMap["internet_url"] = binlogInfo.InternetUrl
			}

			if binlogInfo.Type != nil {
				binlogInfoMap["type"] = binlogInfo.Type
			}

			if binlogInfo.BinlogStartTime != nil {
				binlogInfoMap["binlog_start_time"] = binlogInfo.BinlogStartTime
			}

			if binlogInfo.BinlogFinishTime != nil {
				binlogInfoMap["binlog_finish_time"] = binlogInfo.BinlogFinishTime
			}

			if binlogInfo.Region != nil {
				binlogInfoMap["region"] = binlogInfo.Region
			}

			if binlogInfo.Status != nil {
				binlogInfoMap["status"] = binlogInfo.Status
			}

			if binlogInfo.RemoteInfo != nil {
				remoteInfoList := []interface{}{}
				for _, remoteInfo := range binlogInfo.RemoteInfo {
					remoteInfoMap := map[string]interface{}{}

					if remoteInfo.SubBackupId != nil {
						remoteInfoMap["sub_backup_id"] = remoteInfo.SubBackupId
					}

					if remoteInfo.Region != nil {
						remoteInfoMap["region"] = remoteInfo.Region
					}

					if remoteInfo.Status != nil {
						remoteInfoMap["status"] = remoteInfo.Status
					}

					if remoteInfo.StartTime != nil {
						remoteInfoMap["start_time"] = remoteInfo.StartTime
					}

					if remoteInfo.FinishTime != nil {
						remoteInfoMap["finish_time"] = remoteInfo.FinishTime
					}

					if remoteInfo.Url != nil {
						remoteInfoMap["url"] = remoteInfo.Url
					}

					remoteInfoList = append(remoteInfoList, remoteInfoMap)
				}

				binlogInfoMap["remote_info"] = remoteInfoList
			}

			if binlogInfo.CosStorageType != nil {
				binlogInfoMap["cos_storage_type"] = binlogInfo.CosStorageType
			}

			if binlogInfo.InstanceId != nil {
				binlogInfoMap["instance_id"] = binlogInfo.InstanceId
			}

			ids = append(ids, *binlogInfo.Name)
			tmpList = append(tmpList, binlogInfoMap)
		}

		err = d.Set("items", tmpList)
		if err != nil {
			return err
		}
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
