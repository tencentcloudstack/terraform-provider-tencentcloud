/*
Provides a resource to create a postgresql base_backup

Example Usage

```hcl
resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  tags = {
    "createdBy" = "terraform"
  }
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlBaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBaseBackupCreate,
		Read:   resourceTencentCloudPostgresqlBaseBackupRead,
		Update: resourceTencentCloudPostgresqlBaseBackupUpdate,
		Delete: resourceTencentCloudPostgresqlBaseBackupDelete,
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
	defer logElapsed("resource.tencentcloud_postgresql_base_backup.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgresql.NewCreateBaseBackupRequest()
		response     = postgresql.NewCreateBaseBackupResponse()
		dBInstanceId string
		baseBackupId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateBaseBackup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgresql BaseBackup failed, reason:%+v", logId, err)
		return err
	}

	baseBackupId = *response.Response.BaseBackupId

	d.SetId(strings.Join([]string{dBInstanceId, baseBackupId}, FILED_SP))

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::postgres:%s:uin/:dbInstanceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresqlBaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_base_backup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlBaseBackup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if BaseBackup.Id != nil {
		_ = d.Set("base_backup_id", BaseBackup.Id)
	}

	if BaseBackup.ExpireTime != nil {
		_ = d.Set("new_expire_time", BaseBackup.ExpireTime)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "postgres", "dbInstanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPostgresqlBaseBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_base_backup.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgresql.NewModifyBaseBackupExpireTimeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dBInstanceId := idSplit[0]
	baseBackupId := idSplit[1]

	request.DBInstanceId = &dBInstanceId
	request.BaseBackupId = &baseBackupId

	immutableArgs := []string{"db_instance_id", "base_backup_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("new_expire_time") {
		if v, ok := d.GetOk("new_expire_time"); ok {
			request.NewExpireTime = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyBaseBackupExpireTime(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgresql BaseBackup failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("postgres", "dbInstanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudPostgresqlBaseBackupRead(d, meta)
}

func resourceTencentCloudPostgresqlBaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_base_backup.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
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
