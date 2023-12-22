package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceBackupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"backup_method": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup mode, currently supported: 0-logic backup, 1-physical backup, 2-all backups.The default is logical backup.",
			},

			"backup_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "backup list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"backup_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup mode type.",
						},
						"backup_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup mode name.",
						},
						"backup_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark of backup.",
						},
						"backup_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of backup(KN).",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time of backup.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end time of backup.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup status.",
						},
						"backup_method": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup method.",
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

func dataSourceTencentCloudMongodbInstanceBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_backups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("backup_method"); ok {
		paramMap["backup_method"] = helper.IntInt64(v.(int))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var backupList []*mongodb.BackupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceBackupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupList))
	tmpList := make([]map[string]interface{}, 0, len(backupList))

	if backupList != nil {
		for _, backupInfo := range backupList {
			backupInfoMap := map[string]interface{}{}

			if backupInfo.InstanceId != nil {
				backupInfoMap["instance_id"] = backupInfo.InstanceId
			}

			if backupInfo.BackupType != nil {
				backupInfoMap["backup_type"] = backupInfo.BackupType
			}

			if backupInfo.BackupName != nil {
				backupInfoMap["backup_name"] = backupInfo.BackupName
			}

			if backupInfo.BackupDesc != nil {
				backupInfoMap["backup_desc"] = backupInfo.BackupDesc
			}

			if backupInfo.BackupSize != nil {
				backupInfoMap["backup_size"] = backupInfo.BackupSize
			}

			if backupInfo.StartTime != nil {
				backupInfoMap["start_time"] = backupInfo.StartTime
			}

			if backupInfo.EndTime != nil {
				backupInfoMap["end_time"] = backupInfo.EndTime
			}

			if backupInfo.Status != nil {
				backupInfoMap["status"] = backupInfo.Status
			}

			if backupInfo.BackupMethod != nil {
				backupInfoMap["backup_method"] = backupInfo.BackupMethod
			}

			ids = append(ids, *backupInfo.InstanceId)
			tmpList = append(tmpList, backupInfoMap)
		}

		_ = d.Set("backup_list", tmpList)
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
