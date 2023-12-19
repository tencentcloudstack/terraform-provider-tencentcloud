package cdb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlRoGroupLoadOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRoGroupLoadOperationCreate,
		Read:   resourceTencentCloudMysqlRoGroupLoadOperationRead,
		Delete: resourceTencentCloudMysqlRoGroupLoadOperationDelete,

		Schema: map[string]*schema.Schema{
			"ro_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the RO group, in the format: cdbrg-c1nl9rpv.",
			},
		},
	}
}

func resourceTencentCloudMysqlRoGroupLoadOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group_load_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = mysql.NewBalanceRoGroupLoadRequest()
		roGroupId string
	)
	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
		request.RoGroupId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().BalanceRoGroupLoad(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql roGroupLoadOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(roGroupId)

	return resourceTencentCloudMysqlRoGroupLoadOperationRead(d, meta)
}

func resourceTencentCloudMysqlRoGroupLoadOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group_load_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlRoGroupLoadOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group_load_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
