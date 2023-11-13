/*
Provides a resource to create a cynosdb backup

Example Usage

```hcl
resource "tencentcloud_cynosdb_backup" "backup" {
  cluster_id = &lt;nil&gt;
  backup_type = &lt;nil&gt;
  backup_databases = &lt;nil&gt;
  backup_tables {
		database = ""
		tables = &lt;nil&gt;

  }
  backup_name = &lt;nil&gt;
}
```

Import

cynosdb backup can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_backup.backup backup_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbBackupCreate,
		Read:   resourceTencentCloudCynosdbBackupRead,
		Update: resourceTencentCloudCynosdbBackupUpdate,
		Delete: resourceTencentCloudCynosdbBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"backup_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup type, optional value: logic, logical backup; snapshot, physical backup.",
			},

			"backup_databases": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Backup database list, only valid when BackupType is logic.",
			},

			"backup_tables": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Backup table list, only valid when BackupType is logic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of database.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tables": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Table list.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"backup_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The name of backup.",
			},
		},
	}
}

func resourceTencentCloudCynosdbBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewCreateBackupRequest()
		response  = cynosdb.NewCreateBackupResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_type"); ok {
		request.BackupType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_databases"); ok {
		backupDatabasesSet := v.(*schema.Set).List()
		for i := range backupDatabasesSet {
			backupDatabases := backupDatabasesSet[i].(string)
			request.BackupDatabases = append(request.BackupDatabases, &backupDatabases)
		}
	}

	if v, ok := d.GetOk("backup_tables"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			databaseTables := cynosdb.DatabaseTables{}
			if v, ok := dMap["database"]; ok {
				databaseTables.Database = helper.String(v.(string))
			}
			if v, ok := dMap["tables"]; ok {
				tablesSet := v.(*schema.Set).List()
				for i := range tablesSet {
					tables := tablesSet[i].(string)
					databaseTables.Tables = append(databaseTables.Tables, &tables)
				}
			}
			request.BackupTables = append(request.BackupTables, &databaseTables)
		}
	}

	if v, ok := d.GetOk("backup_name"); ok {
		request.BackupName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb backup failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudCynosdbBackupRead(d, meta)
}

func resourceTencentCloudCynosdbBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	backupId := d.Id()

	backup, err := service.DescribeCynosdbBackupById(ctx, clusterId)
	if err != nil {
		return err
	}

	if backup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backup.ClusterId != nil {
		_ = d.Set("cluster_id", backup.ClusterId)
	}

	if backup.BackupType != nil {
		_ = d.Set("backup_type", backup.BackupType)
	}

	if backup.BackupDatabases != nil {
		_ = d.Set("backup_databases", backup.BackupDatabases)
	}

	if backup.BackupTables != nil {
		backupTablesList := []interface{}{}
		for _, backupTables := range backup.BackupTables {
			backupTablesMap := map[string]interface{}{}

			if backup.BackupTables.Database != nil {
				backupTablesMap["database"] = backup.BackupTables.Database
			}

			if backup.BackupTables.Tables != nil {
				backupTablesMap["tables"] = backup.BackupTables.Tables
			}

			backupTablesList = append(backupTablesList, backupTablesMap)
		}

		_ = d.Set("backup_tables", backupTablesList)

	}

	if backup.BackupName != nil {
		_ = d.Set("backup_name", backup.BackupName)
	}

	return nil
}

func resourceTencentCloudCynosdbBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyBackupNameRequest()

	backupId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "backup_type", "backup_databases", "backup_tables", "backup_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyBackupName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb backup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbBackupRead(d, meta)
}

func resourceTencentCloudCynosdbBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	backupId := d.Id()

	if err := service.DeleteCynosdbBackupById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
