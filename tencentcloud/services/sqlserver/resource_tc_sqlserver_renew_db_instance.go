package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverRenewDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRenewDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverRenewDBInstanceRead,
		Update: resourceTencentCloudSqlserverRenewDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRenewDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     1,
				Description: "How many months to renew, the value range is 1-48, the default is 1.",
			},
		},
	}
}

func resourceTencentCloudSqlserverRenewDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
		period     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("period"); ok {
		period = strconv.Itoa(v.(int))
	} else {
		period = "1"
	}

	d.SetId(strings.Join([]string{instanceId, period}, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverRenewDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRenewDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	period := idSplit[1]

	renewDBInstance, err := service.DescribeSqlserverRenewDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if renewDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRenewDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if renewDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", renewDBInstance.InstanceId)
	}

	tmpPeriod, _ := strconv.Atoi(period)
	_ = d.Set("period", tmpPeriod)

	return nil
}

func resourceTencentCloudSqlserverRenewDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = sqlserver.NewRenewDBInstanceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	request.InstanceId = &instanceId

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().RenewDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver renewDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRenewDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRenewDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_renew_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
