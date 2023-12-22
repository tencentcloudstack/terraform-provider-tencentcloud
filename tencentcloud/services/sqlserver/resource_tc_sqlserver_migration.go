package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverMigrationCreate,
		Read:   resourceTencentCloudSqlserverMigrationRead,
		Update: resourceTencentCloudSqlserverMigrationUpdate,
		Delete: resourceTencentCloudSqlserverMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"migrate_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the migration task.",
			},

			"migrate_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Migration type (1 structure migration 2 data migration 3 incremental synchronization).",
			},

			"source_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Type of migration source 1 TencentDB for SQLServer 2 Cloud server self-built SQLServer database 4 SQLServer backup and restore 5 SQLServer backup and restore (COS mode).",
			},

			"source": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration source.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the migration source instance, which is used when MigrateType=1 (TencentDB for SQLServers). The format is mssql-si2823jyl.",
						},
						"cvm_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the migration source Cvm, used when MigrateType=2 (cloud server self-built SQL Server database).",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The Vpc network ID of the migration source Cvm is used when MigrateType=2 (cloud server self-built SQL Server database). The format is as follows vpc-6ys9ont9.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The subnet ID under the Vpc of the source Cvm is used when MigrateType=2 (ECS self-built SQL Server database). The format is as follows subnet-h9extioi.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User name, MigrateType=1 or MigrateType=2.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password, MigrateType=1 or MigrateType=2.",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Migrate the intranet IP of the self-built database of the source Cvm, and use it when MigrateType=2 (self-built SQL Server database of the cloud server).",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port number of the self-built database of the migration source Cvm, which is used when MigrateType=2 (self-built SQL Server database of the cloud server).",
						},
						"url": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The source backup address for offline migration. MigrateType=4 or MigrateType=5.",
						},
						"url_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The source backup password for offline migration, MigrateType=4 or MigrateType=5.",
						},
					},
				},
			},

			"target": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the migration target instance, in the format mssql-si2823jyl.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User name of the migration target instance.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password of the migration target instance.",
						},
					},
				},
			},

			"migrate_db_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Migrate DB objects. Offline migration is not used (SourceType=4 or SourceType=5).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the migration database.",
						},
					},
				},
			},

			"rename_restore": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Restore and rename the database in ReNameRestoreDatabase. If it is not filled in, the restored database will be named by default and all databases will be restored. Valid if SourceType=5.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the library. If oldName does not exist, a failure is returned.It can be left blank when used for offline migration tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "When the new name of the library is used for offline migration, if it is not filled in, it will be named according to OldName. OldName and NewName cannot be filled in at the same time. OldName and NewName must be filled in and cannot be duplicate when used for cloning database.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = sqlserver.NewCreateMigrationRequest()
		response  = sqlserver.NewCreateMigrationResponse()
		migrateId string
	)
	if v, ok := d.GetOk("migrate_name"); ok {
		request.MigrateName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("migrate_type"); v != nil {
		request.MigrateType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("source_type"); v != nil {
		request.SourceType = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "source"); ok {
		migrateSource := sqlserver.MigrateSource{}
		if v, ok := dMap["instance_id"]; ok {
			migrateSource.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_id"]; ok {
			migrateSource.CvmId = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			migrateSource.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			migrateSource.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateSource.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateSource.Password = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			migrateSource.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			migrateSource.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["url"]; ok {
			urlSet := v.(*schema.Set).List()
			for i := range urlSet {
				url := urlSet[i].(string)
				migrateSource.Url = append(migrateSource.Url, &url)
			}
		}
		if v, ok := dMap["url_password"]; ok {
			migrateSource.UrlPassword = helper.String(v.(string))
		}
		request.Source = &migrateSource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target"); ok {
		migrateTarget := sqlserver.MigrateTarget{}
		if v, ok := dMap["instance_id"]; ok {
			migrateTarget.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateTarget.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateTarget.Password = helper.String(v.(string))
		}
		request.Target = &migrateTarget
	}

	if v, ok := d.GetOk("migrate_db_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			migrateDB := sqlserver.MigrateDB{}
			if v, ok := dMap["db_name"]; ok {
				migrateDB.DBName = helper.String(v.(string))
			}
			request.MigrateDBSet = append(request.MigrateDBSet, &migrateDB)
		}
	}

	if v, ok := d.GetOk("rename_restore"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			renameRestoreDatabase := sqlserver.RenameRestoreDatabase{}
			if v, ok := dMap["old_name"]; ok {
				renameRestoreDatabase.OldName = helper.String(v.(string))
			}
			if v, ok := dMap["new_name"]; ok {
				renameRestoreDatabase.NewName = helper.String(v.(string))
			}
			request.RenameRestore = append(request.RenameRestore, &renameRestoreDatabase)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateMigration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver migration failed, reason:%+v", logId, err)
		return err
	}

	migrateId = helper.Int64ToStr(*response.Response.MigrateId)
	d.SetId(migrateId)

	return resourceTencentCloudSqlserverMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	migrationId := d.Id()

	migration, err := service.DescribeSqlserverMigrationById(ctx, migrationId)
	if err != nil {
		return err
	}

	if migration == nil {
		d.SetId("")
		return fmt.Errorf("resource `SqlserverMigration` %s does not exist", d.Id())
	}

	if migration.MigrateName != nil {
		_ = d.Set("migrate_name", migration.MigrateName)
	}

	if migration.MigrateType != nil {
		_ = d.Set("migrate_type", migration.MigrateType)
	}

	if migration.SourceType != nil {
		_ = d.Set("source_type", migration.SourceType)
	}

	if migration.Source != nil {
		sourceMap := map[string]interface{}{}

		if migration.Source.InstanceId != nil {
			sourceMap["instance_id"] = migration.Source.InstanceId
		}

		if migration.Source.CvmId != nil {
			sourceMap["cvm_id"] = migration.Source.CvmId
		}

		if migration.Source.VpcId != nil {
			sourceMap["vpc_id"] = migration.Source.VpcId
		}

		if migration.Source.SubnetId != nil {
			sourceMap["subnet_id"] = migration.Source.SubnetId
		}

		if migration.Source.UserName != nil {
			sourceMap["user_name"] = migration.Source.UserName
		}

		if migration.Source.Password != nil {
			sourceMap["password"] = migration.Source.Password
		}

		if migration.Source.Ip != nil {
			sourceMap["ip"] = migration.Source.Ip
		}

		if migration.Source.Port != nil {
			sourceMap["port"] = migration.Source.Port
		}

		if migration.Source.Url != nil {
			sourceMap["url"] = migration.Source.Url
		}

		if migration.Source.UrlPassword != nil {
			sourceMap["url_password"] = migration.Source.UrlPassword
		}

		_ = d.Set("source", []interface{}{sourceMap})
	}

	if migration.Target != nil {
		targetMap := map[string]interface{}{}

		if migration.Target.InstanceId != nil {
			targetMap["instance_id"] = migration.Target.InstanceId
		}

		if migration.Target.UserName != nil {
			targetMap["user_name"] = migration.Target.UserName
		}

		if migration.Target.Password != nil {
			targetMap["password"] = migration.Target.Password
		}

		_ = d.Set("target", []interface{}{targetMap})
	}

	if migration.MigrateDBSet != nil {
		migrateDBSetList := []interface{}{}
		for _, migrateDB := range migration.MigrateDBSet {
			migrateDBSetMap := map[string]interface{}{}

			if migrateDB.DBName != nil {
				migrateDBSetMap["db_name"] = migrateDB.DBName
			}

			migrateDBSetList = append(migrateDBSetList, migrateDBSetMap)
		}

		_ = d.Set("migrate_db_set", migrateDBSetList)

	}

	// omit rename_restore because read api doesn't return it
	// _ = d.Set("rename_restore", renameRestoreList)

	return nil
}

func resourceTencentCloudSqlserverMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := sqlserver.NewModifyMigrationRequest()
	migrateId := d.Id()

	request.MigrateId = helper.StrToUint64Point(migrateId)
	if d.HasChange("rename_restore") {
		o, _ := d.GetChange("rename_restore")
		_ = d.Set("rename_restore", o)
		return fmt.Errorf("argument `%s` cannot be changed", d.Id())
	}
	if d.HasChange("migrate_name") {
		if v, ok := d.GetOk("migrate_name"); ok {
			request.MigrateName = helper.String(v.(string))
		}
	}

	if d.HasChange("migrate_type") {
		if v, _ := d.GetOk("migrate_type"); v != nil {
			request.MigrateType = helper.IntUint64(v.(int))
		}
		if v, _ := d.GetOk("source_type"); v != nil {
			request.SourceType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("source_type") {
		if v, _ := d.GetOk("source_type"); v != nil {
			request.SourceType = helper.IntUint64(v.(int))
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "source"); ok {
		migrateSource := sqlserver.MigrateSource{}
		if v, ok := dMap["instance_id"]; ok {
			migrateSource.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_id"]; ok {
			migrateSource.CvmId = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			migrateSource.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			migrateSource.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateSource.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateSource.Password = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			migrateSource.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			migrateSource.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["url"]; ok {
			urlSet := v.(*schema.Set).List()
			for i := range urlSet {
				url := urlSet[i].(string)
				migrateSource.Url = append(migrateSource.Url, &url)
			}
		}
		if v, ok := dMap["url_password"]; ok {
			migrateSource.UrlPassword = helper.String(v.(string))
		}
		request.Source = &migrateSource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target"); ok {
		migrateTarget := sqlserver.MigrateTarget{}
		if v, ok := dMap["instance_id"]; ok {
			migrateTarget.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateTarget.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateTarget.Password = helper.String(v.(string))
		}
		request.Target = &migrateTarget
	}

	if v, ok := d.GetOk("migrate_db_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			migrateDB := sqlserver.MigrateDB{}
			if v, ok := dMap["db_name"]; ok {
				migrateDB.DBName = helper.String(v.(string))
			}
			request.MigrateDBSet = append(request.MigrateDBSet, &migrateDB)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyMigration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver migration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	migrateId := d.Id()

	if err := service.DeleteSqlserverMigrationById(ctx, migrateId); err != nil {
		return err
	}

	return nil
}
