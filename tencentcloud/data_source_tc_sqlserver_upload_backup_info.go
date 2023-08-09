/*
Use this data source to query detailed information of sqlserver upload_backup_info

Example Usage

```hcl
data "tencentcloud_sqlserver_upload_backup_info" "example" {
  instance_id         = "mssql-qelbzgwf"
  backup_migration_id = "mssql-backup-migration-8a0f3eht"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverUploadBackupInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverUploadBackupInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"backup_migration_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup import task ID, which is returned through the API CreateBackupMigration.",
			},
			"bucket_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Bucket name.",
			},
			"region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Bucket location information.",
			},
			"path": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Storage path.",
			},
			"tmp_secret_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary key ID.",
			},
			"tmp_secret_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary key (Key).",
			},
			"x_cos_security_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary key (Token).",
			},
			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary key start time.",
			},
			"expired_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Temporary key expiration time.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverUploadBackupInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_upload_backup_info.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		service          = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		uploadBackupInfo *sqlserver.DescribeUploadBackupInfoResponseParams
		instanceId       string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		paramMap["BackupMigrationId"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverUploadBackupInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		uploadBackupInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	if uploadBackupInfo.BucketName != nil {
		_ = d.Set("bucket_name", uploadBackupInfo.BucketName)
	}

	if uploadBackupInfo.Region != nil {
		_ = d.Set("region", uploadBackupInfo.Region)
	}

	if uploadBackupInfo.Path != nil {
		_ = d.Set("path", uploadBackupInfo.Path)
	}

	if uploadBackupInfo.TmpSecretId != nil {
		_ = d.Set("tmp_secret_id", uploadBackupInfo.TmpSecretId)
	}

	if uploadBackupInfo.TmpSecretKey != nil {
		_ = d.Set("tmp_secret_key", uploadBackupInfo.TmpSecretKey)
	}

	if uploadBackupInfo.XCosSecurityToken != nil {
		_ = d.Set("x_cos_security_token", uploadBackupInfo.XCosSecurityToken)
	}

	if uploadBackupInfo.StartTime != nil {
		_ = d.Set("start_time", uploadBackupInfo.StartTime)
	}

	if uploadBackupInfo.ExpiredTime != nil {
		_ = d.Set("expired_time", uploadBackupInfo.ExpiredTime)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
