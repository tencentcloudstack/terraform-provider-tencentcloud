/*
Provides a resource to create a dcdb modify_real_server_access_strategy

Example Usage

```hcl
resource "tencentcloud_dcdb_modify_real_server_access_strategy" "modify_real_server_access_strategy" {
  instance_id = ""
  rs_access_strategy =
}
```

Import

dcdb modify_real_server_access_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_modify_real_server_access_strategy.modify_real_server_access_strategy modify_real_server_access_strategy_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbModifyRealServerAccessStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbModifyRealServerAccessStrategyCreate,
		Read:   resourceTencentCloudDcdbModifyRealServerAccessStrategyRead,
		Delete: resourceTencentCloudDcdbModifyRealServerAccessStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"rs_access_strategy": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "RS nearest access mode, 0-no policy, 1-nearest access.",
			},
		},
	}
}

func resourceTencentCloudDcdbModifyRealServerAccessStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_real_server_access_strategy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewModifyRealServerAccessStrategyRequest()
		response   = dcdb.NewModifyRealServerAccessStrategyResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("rs_access_strategy"); v != nil {
		request.RsAccessStrategy = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyRealServerAccessStrategy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb modifyRealServerAccessStrategy failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbModifyRealServerAccessStrategyRead(d, meta)
}

func resourceTencentCloudDcdbModifyRealServerAccessStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_real_server_access_strategy.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbModifyRealServerAccessStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_modify_real_server_access_strategy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
