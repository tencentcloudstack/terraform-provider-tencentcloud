package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRedisBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisBackupRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "start time, such as 2017-02-08 19:09:26.Query the list of backups that the instance started backing up during the [beginTime, endTime] time period.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time, such as 2017-02-08 19:09:26.Query the list of backups that the instance started backing up during the [beginTime, endTime] time period.",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Status of the backup task:1: Backup is in the process.2: The backup is normal.3: Backup to RDB file processing.4: RDB conversion completed.-1: The backup has expired.-2: Backup deleted.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance name, which supports fuzzy search based on instance name.",
			},

			"backup_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "An array of backups for the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup start time.",
						},
						"backup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup ID.",
						},
						"backup_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup type.1: User-initiated manual backup.0: System-initiated backup in the early morning.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup status.1: The backup is locked by another process.2: The backup is normal and not locked by any process.-1: The backup has expired.3: The backup is being exported.4: The backup export is successful.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notes information for the backup.",
						},
						"locked": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the backup is locked.0: Not locked.1: Has been locked.",
						},
						"backup_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal fields, which can be ignored by the user.",
						},
						"full_backup": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal fields, which can be ignored by the user.",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal fields, which can be ignored by the user.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of instance.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the backup is located.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup end time.",
						},
						"file_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up file types.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup file expiration time.",
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

func dataSourceTencentCloudRedisBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("begin_time"); ok {
		paramMap["begin_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		statusList := []*int64{}
		for i := range statusSet {
			status := statusSet[i].(int)
			statusList = append(statusList, helper.IntInt64(status))
		}
		paramMap["status"] = statusList
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var backupSet []*redis.RedisBackupSet

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisBackupByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		backupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupSet))
	tmpList := make([]map[string]interface{}, 0, len(backupSet))

	if backupSet != nil {
		for _, redisBackupSet := range backupSet {
			redisBackupSetMap := map[string]interface{}{}

			if redisBackupSet.StartTime != nil {
				redisBackupSetMap["start_time"] = redisBackupSet.StartTime
			}

			if redisBackupSet.BackupId != nil {
				redisBackupSetMap["backup_id"] = redisBackupSet.BackupId
			}

			if redisBackupSet.BackupType != nil {
				redisBackupSetMap["backup_type"] = redisBackupSet.BackupType
			}

			if redisBackupSet.Status != nil {
				redisBackupSetMap["status"] = redisBackupSet.Status
			}

			if redisBackupSet.Remark != nil {
				redisBackupSetMap["remark"] = redisBackupSet.Remark
			}

			if redisBackupSet.Locked != nil {
				redisBackupSetMap["locked"] = redisBackupSet.Locked
			}

			if redisBackupSet.BackupSize != nil {
				redisBackupSetMap["backup_size"] = redisBackupSet.BackupSize
			}

			if redisBackupSet.FullBackup != nil {
				redisBackupSetMap["full_backup"] = redisBackupSet.FullBackup
			}

			if redisBackupSet.InstanceType != nil {
				redisBackupSetMap["instance_type"] = redisBackupSet.InstanceType
			}

			if redisBackupSet.InstanceId != nil {
				redisBackupSetMap["instance_id"] = redisBackupSet.InstanceId
			}

			if redisBackupSet.InstanceName != nil {
				redisBackupSetMap["instance_name"] = redisBackupSet.InstanceName
			}

			if redisBackupSet.Region != nil {
				redisBackupSetMap["region"] = redisBackupSet.Region
			}

			if redisBackupSet.EndTime != nil {
				redisBackupSetMap["end_time"] = redisBackupSet.EndTime
			}

			if redisBackupSet.FileType != nil {
				redisBackupSetMap["file_type"] = redisBackupSet.FileType
			}

			if redisBackupSet.ExpireTime != nil {
				redisBackupSetMap["expire_time"] = redisBackupSet.ExpireTime
			}

			ids = append(ids, *redisBackupSet.InstanceId)
			tmpList = append(tmpList, redisBackupSetMap)
		}

		_ = d.Set("backup_set", tmpList)
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
