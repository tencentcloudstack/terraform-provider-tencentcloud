package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				ForceNew:    true,
				Description: "Database character encoding, such as UTF8, LATIN1, LATIN2, WIN1250, WIN1251, WIN1252, KOI8R, EUC_JP, EUC_KR. Default: UTF8.",
			},

			"collate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Database collation rule.",
			},

			"ctype": {
				Type:        schema.TypeString,
				Optional:    true,
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
		request      = postgresql.NewCreateDatabaseRequest()
		dbInstanceId string
		databaseName string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dbInstanceId = v.(string)
	}

	if v, ok := d.GetOk("database_name"); ok {
		request.DatabaseName = helper.String(v.(string))
		databaseName = v.(string)
	}

	if v, ok := d.GetOk("database_owner"); ok {
		request.DatabaseOwner = helper.String(v.(string))
	}

	if v, ok := d.GetOk("encoding"); ok {
		request.Encoding = helper.String(v.(string))
	}

	if v, ok := d.GetOk("collate"); ok {
		request.Collate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ctype"); ok {
		request.Ctype = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgres database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create postgres database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create postgres database success, dbInstanceId=%s, databaseName=%s", logId, dbInstanceId, databaseName)
	if dbInstanceId == "" || databaseName == "" {
		return fmt.Errorf("db_instance_id or database_name is empty, dbInstanceId=%s, databaseName=%s", dbInstanceId, databaseName)
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

	request := postgresql.NewDescribeDatabasesRequest()
	request.DBInstanceId = helper.String(dbInstanceId)
	request.Limit = helper.Uint64(100)
	request.Offset = helper.Uint64(0)
	request.Filters = []*postgresql.Filter{
		{
			Name:   helper.String("database-name"),
			Values: []*string{helper.String(databaseName)},
		},
	}

	var databaseInfo *postgresql.Database
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeDatabasesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe postgres database failed, Response is nil."))
		}

		for _, item := range result.Response.Databases {
			if item != nil && item.DatabaseName != nil && *item.DatabaseName == databaseName {
				databaseInfo = item
				break
			}
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read postgres database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if databaseInfo == nil {
		log.Printf("[CRUD] postgres database id=%s", d.Id())
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

	needChange := false
	mutableArgs := []string{"database_owner"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := postgresql.NewModifyDatabaseOwnerRequest()
		request.DBInstanceId = helper.String(dbInstanceId)
		request.DatabaseName = helper.String(databaseName)

		if v, ok := d.GetOk("database_owner"); ok {
			request.DatabaseOwner = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDatabaseOwnerWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify postgres database owner failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update postgres database failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudPostgresqlDatabaseRead(d, meta)
}

func resourceTencentCloudPostgresqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_database.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = postgresql.NewDeleteDatabaseRequest()
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dbInstanceId := idSplit[0]
	databaseName := idSplit[1]

	request.DBInstanceId = helper.String(dbInstanceId)
	request.DatabaseName = helper.String(databaseName)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DeleteDatabaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete postgres database failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete postgres database failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
