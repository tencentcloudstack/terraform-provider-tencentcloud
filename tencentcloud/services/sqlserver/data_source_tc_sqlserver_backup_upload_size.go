package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverBackupUploadSize() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverBackupUploadSizeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of imported target instance.",
			},
			"backup_migration_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup import task ID, which is returned through the API CreateBackupMigration.",
			},
			"incremental_migration_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Incremental import task ID.",
			},
			"cos_upload_backup_file_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of uploaded backups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup size.",
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

func dataSourceTencentCloudSqlserverBackupUploadSizeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_backup_upload_size.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		cosUploadBackupFileSet []*sqlserver.CosUploadBackupFile
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_migration_id"); ok {
		paramMap["BackupMigrationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("incremental_migration_id"); ok {
		paramMap["IncrementalMigrationId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverBackupUploadSizeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		cosUploadBackupFileSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(cosUploadBackupFileSet))
	tmpList := make([]map[string]interface{}, 0, len(cosUploadBackupFileSet))

	if cosUploadBackupFileSet != nil {
		for _, cosUploadBackupFile := range cosUploadBackupFileSet {
			cosUploadBackupFileMap := map[string]interface{}{}

			if cosUploadBackupFile.FileName != nil {
				cosUploadBackupFileMap["file_name"] = cosUploadBackupFile.FileName
			}

			if cosUploadBackupFile.Size != nil {
				cosUploadBackupFileMap["size"] = cosUploadBackupFile.Size
			}

			ids = append(ids, *cosUploadBackupFile.FileName)
			tmpList = append(tmpList, cosUploadBackupFileMap)
		}

		_ = d.Set("cos_upload_backup_file_set", tmpList)
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
