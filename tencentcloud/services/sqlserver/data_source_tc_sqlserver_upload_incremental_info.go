package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverUploadIncrementalInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverUploadIncrementalInfoRead,
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
			"incremental_migration_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the incremental import task.",
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

func dataSourceTencentCloudSqlserverUploadIncrementalInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_upload_incremental_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                 = tccommon.GetLogId(tccommon.ContextNil)
		ctx                   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service               = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		uploadIncrementalInfo *sqlserver.DescribeUploadIncrementalInfoResponseParams
		instanceId            string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		paramMap["BackupMigrationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("incremental_migration_id"); ok {
		paramMap["IncrementalMigrationId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverUploadIncrementalInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		uploadIncrementalInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	if uploadIncrementalInfo.BucketName != nil {
		_ = d.Set("bucket_name", uploadIncrementalInfo.BucketName)
	}

	if uploadIncrementalInfo.Region != nil {
		_ = d.Set("region", uploadIncrementalInfo.Region)
	}

	if uploadIncrementalInfo.Path != nil {
		_ = d.Set("path", uploadIncrementalInfo.Path)
	}

	if uploadIncrementalInfo.TmpSecretId != nil {
		_ = d.Set("tmp_secret_id", uploadIncrementalInfo.TmpSecretId)
	}

	if uploadIncrementalInfo.TmpSecretKey != nil {
		_ = d.Set("tmp_secret_key", uploadIncrementalInfo.TmpSecretKey)
	}

	if uploadIncrementalInfo.XCosSecurityToken != nil {
		_ = d.Set("x_cos_security_token", uploadIncrementalInfo.XCosSecurityToken)
	}

	if uploadIncrementalInfo.StartTime != nil {
		_ = d.Set("start_time", uploadIncrementalInfo.StartTime)
	}

	if uploadIncrementalInfo.ExpiredTime != nil {
		_ = d.Set("expired_time", uploadIncrementalInfo.ExpiredTime)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
