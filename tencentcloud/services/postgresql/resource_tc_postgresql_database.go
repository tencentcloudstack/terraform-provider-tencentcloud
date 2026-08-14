package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudPostgresqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlDatabaseCreate,
		Read:   resourceTencentCloudPostgresqlDatabaseRead,
		Update: resourceTencentCloudPostgresqlDatabaseUpdate,
		Delete: resourceTencentCloudPostgresqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DB instance ID, such as postgres-6fego161.",
			},

			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database name.",
			},

			"database_owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database owner account.",
			},

			"encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Database character encoding, such as UTF8, LATIN1, LATIN2, WIN1250, WIN1251, WIN1252, KOI8R, EUC_JP, EUC_KR. Default: UTF8.",
			},

			"collate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Database collation rule.",
			},

			"ctype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Database character classification.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_database.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dbInstanceId string
		databaseName string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("database_name"); ok {
		databaseName = v.(string)
	}

	if dbInstanceId == "" || databaseName == "" {
		return fmt.Errorf("db_instance_id or database_name is empty, dbInstanceId=%s, databaseName=%s", dbInstanceId, databaseName)
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.CreatePostgresqlDatabase(ctx, dbInstanceId, databaseName, d.Get("database_owner").(string), d.Get("encoding").(string), d.Get("collate").(string), d.Get("ctype").(string)); err != nil {
		log.Printf("[CRITAL]%s create postgresql database failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dbInstanceId, databaseName}, tccommon.FILED_SP))
	return resourceTencentCloudPostgresqlDatabaseRead(d, meta)
}

func resourceTencentCloudPostgresqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_database.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	databaseName := idSplit[1]

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	databaseInfo, err := service.DescribePostgresqlDatabaseById(ctx, dbInstanceId, databaseName)
	if err != nil {
		log.Printf("[CRITAL]%s read postgresql database failed, reason:%+v", logId, err)
		return err
	}

	if databaseInfo == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgresql_database` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("db_instance_id", dbInstanceId)

	if databaseInfo.DatabaseName != nil {
		_ = d.Set("database_name", databaseInfo.DatabaseName)
	}

	if databaseInfo.DatabaseOwner != nil {
		_ = d.Set("database_owner", databaseInfo.DatabaseOwner)
	}

	if databaseInfo.Encoding != nil {
		_ = d.Set("encoding", databaseInfo.Encoding)
	}

	if databaseInfo.Collate != nil {
		_ = d.Set("collate", databaseInfo.Collate)
	}

	if databaseInfo.Ctype != nil {
		_ = d.Set("ctype", databaseInfo.Ctype)
	}

	return nil
}

func resourceTencentCloudPostgresqlDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_database.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	databaseName := idSplit[1]

	if d.HasChange("database_owner") {
		service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := service.ModifyPostgresqlDatabaseOwner(ctx, dbInstanceId, databaseName, d.Get("database_owner").(string)); err != nil {
			log.Printf("[CRITAL]%s update postgresql database owner failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudPostgresqlDatabaseRead(d, meta)
}

func resourceTencentCloudPostgresqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_database.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	databaseName := idSplit[1]

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.DeletePostgresqlDatabaseById(ctx, dbInstanceId, databaseName); err != nil {
		log.Printf("[CRITAL]%s delete postgresql database failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
