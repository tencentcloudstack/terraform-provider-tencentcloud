package postgresql

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlReadonlyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlReadOnlyAttachmentCreate,
		Read:   resourceTencentCloudPostgresqlReadOnlyAttachmentRead,
		//Update: resourceTencentCloudPostgresqlReadOnlyAttachmentUpdate,
		Delete: resourceTencentCLoudPostgresqlReadOnlyAttachmentDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Read only instance ID.",
			},
			"read_only_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Read only group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlReadOnlyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewAddDBInstanceToReadOnlyGroupRequest()
		dbInstanceId string
		groupId      string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInstanceId = v.(string)
		request.DBInstanceId = helper.String(dbInstanceId)

	}
	if v, ok := d.GetOk("read_only_group_id"); ok {
		groupId = v.(string)
		request.ReadOnlyGroupId = helper.String(groupId)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().AddDBInstanceToReadOnlyGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}
	instanceId := helper.IdFormat(dbInstanceId, groupId)
	d.SetId(instanceId)

	return nil
}

func resourceTencentCloudPostgresqlReadOnlyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_attachment.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	postgresqlService := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	_, err := postgresqlService.DescribePostgresqlReadOnlyGroupById(ctx, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCLoudPostgresqlReadOnlyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_readonly_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := postgresql.NewRemoveDBInstanceFromReadOnlyGroupRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dbInstanceId := idSplit[0]
	groupId := idSplit[1]
	request.ReadOnlyGroupId = helper.String(groupId)
	request.DBInstanceId = helper.String(dbInstanceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().RemoveDBInstanceFromReadOnlyGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
