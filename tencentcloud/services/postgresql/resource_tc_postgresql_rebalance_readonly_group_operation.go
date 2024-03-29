package postgresql

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlRebalanceReadonlyGroupOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationCreate,
		Read:   resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead,
		Delete: resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationDelete,
		Schema: map[string]*schema.Schema{
			"read_only_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "readonly Group ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = postgresql.NewRebalanceReadOnlyGroupRequest()
		readOnlyGroupId string
	)
	if v, ok := d.GetOk("read_only_group_id"); ok {
		request.ReadOnlyGroupId = helper.String(v.(string))
		readOnlyGroupId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().RebalanceReadOnlyGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql RebalanceReadonlyGroupOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(readOnlyGroupId)

	return resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlRebalanceReadonlyGroupOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_rebalance_readonly_group_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
