/*
Provides a resource to create a cynosdb backup

Example Usage

```hcl
resource "tencentcloud_cynosdb_backup" "backup" {
  cluster_id = "cynosdbmysql-bws8h88b"
  backup_type = "logic"
  backup_name = "testname"
}
```

Import

cynosdb backup can be imported using the id, e.g. cynosdbmysql-bws8h88b#297272

```
terraform import tencentcloud_cynosdb_backup.backup {cluster_id}#${backup_id}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
							Description: "&quot;The name of database.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.",
						},
						"tables": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "&quot;Table list.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.",
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request   = cynosdb.NewCreateBackupRequest()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
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

	result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateBackup(request)
	currentTimeUnix := time.Now().Unix()
	if e != nil {
		return e
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	params := make(map[string]interface{})
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		backups, e := service.DescribeCynosdbBackup(ctx, clusterId, params)

		if e != nil {
			return retryError(e)
		}
		for _, backup := range backups {
			backupSnapshotTime, _ := time.ParseInLocation("2006-01-02 15:04:05", *backup.SnapshotTime, time.Local)
			backupSnapshotTimeUnix := backupSnapshotTime.Unix()
			log.Printf("*backup.SnapshotTime: %v, backupSnapshotTimeUnix: %v, currentTimeUnix: %v", *backup.SnapshotTime, backupSnapshotTimeUnix, currentTimeUnix)

			if backupSnapshotTimeUnix-currentTimeUnix > 0 && backupSnapshotTimeUnix-currentTimeUnix <= 5 {
				d.SetId(clusterId + FILED_SP + strconv.FormatInt(*backup.BackupId, 10))
				return nil
			}
		}
		return resource.RetryableError(fmt.Errorf("Backup not finished"))
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb backup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbBackupRead(d, meta)
}

func resourceTencentCloudCynosdbBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	backupId := idSplit[1]

	backup, err := service.DescribeCynosdbBackupById(ctx, clusterId, backupId)
	if err != nil {
		return err
	}

	if backup == nil {
		d.SetId("")
		return fmt.Errorf("resource `CynosdbBackup` %s does not exist", d.Id())
	}
	_ = d.Set("cluster_id", backupId)

	if backup.BackupType != nil {
		_ = d.Set("backup_type", backup.BackupType)
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	backupId := idSplit[1]

	request.ClusterId = &clusterId
	backupIdInt64, err := strconv.ParseInt(backupId, 10, 64)
	if err != nil {
		return err
	}
	request.BackupId = &backupIdInt64

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
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
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	backupId := idSplit[1]

	if err := service.DeleteCynosdbBackupById(ctx, clusterId, backupId); err != nil {
		return err
	}

	return nil
}
