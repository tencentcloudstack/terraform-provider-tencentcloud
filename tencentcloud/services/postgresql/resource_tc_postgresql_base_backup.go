package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlBaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBaseBackupCreate,
		Read:   resourceTencentCloudPostgresqlBaseBackupRead,
		Update: resourceTencentCloudPostgresqlBaseBackupUpdate,
		Delete: resourceTencentCloudPostgresqlBaseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"base_backup_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Base backup ID.",
			},

			"new_expire_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "New expiration time.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlBaseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_base_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request      = postgresql.NewCreateBaseBackupRequest()
		response     = postgresql.NewCreateBaseBackupResponse()
		dBInstanceId string
		baseBackupId string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().CreateBaseBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql BaseBackup failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgresql BaseBackup failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.BaseBackupId == nil {
		return fmt.Errorf("BaseBackupId is nil.")
	}

	baseBackupId = *response.Response.BaseBackupId

	d.SetId(strings.Join([]string{dBInstanceId, baseBackupId}, tccommon.FILED_SP))

	// wait
	waitReq := postgresql.NewDescribeBaseBackupsRequest()
	waitReq.Filters = []*postgresql.Filter{
		{
			Name:   helper.String("db-instance-id"),
			Values: helper.Strings([]string{dBInstanceId}),
		},
		{
			Name:   helper.String("base-backup-id"),
			Values: helper.Strings([]string{baseBackupId}),
		},
	}

	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DescribeBaseBackupsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe base backups failed, Response is nil."))
		}

		if len(result.Response.BaseBackupSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("BaseBackupSet is nil."))
		}

		tmpObj := result.Response.BaseBackupSet[0]
		if tmpObj.State == nil {
			return resource.NonRetryableError(fmt.Errorf("State is nil."))
		}

		if *tmpObj.State == "finished" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Base backup is still running..."))
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe base backup failed, reason:%+v", logId, err)
		return err
	}

	if v, ok := d.GetOk("new_expire_time"); ok {
		request := postgresql.NewModifyBaseBackupExpireTimeRequest()
		request.NewExpireTime = helper.String(v.(string))
		request.DBInstanceId = &dBInstanceId
		request.BaseBackupId = &baseBackupId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyBaseBackupExpireTime(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update postgresql BaseBackup failed, reason:%+v", logId, err)
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresqlBaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_base_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	BaseBackup, err := service.DescribePostgresqlBaseBackupById(ctx, baseBackupId)
	if err != nil {
		return err
	}

	if BaseBackup == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgresql_base_backup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if BaseBackup.Id != nil {
		_ = d.Set("base_backup_id", BaseBackup.Id)
	}

	if BaseBackup.ExpireTime != nil {
		_ = d.Set("new_expire_time", BaseBackup.ExpireTime)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlBaseBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_base_backup.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	immutableArgs := []string{"db_instance_id", "base_backup_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("new_expire_time") {
		request := postgresql.NewModifyBaseBackupExpireTimeRequest()
		if v, ok := d.GetOk("new_expire_time"); ok {
			request.NewExpireTime = helper.String(v.(string))
		}

		request.DBInstanceId = &dBInstanceId
		request.BaseBackupId = &baseBackupId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyBaseBackupExpireTime(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update postgresql BaseBackup failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("postgres", "dbInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresqlBaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_base_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	if err := service.DeletePostgresqlBaseBackupById(ctx, dBInstanceId, baseBackupId); err != nil {
		return err
	}

	return nil
}
