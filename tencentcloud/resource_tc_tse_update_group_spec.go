/*
Provides a resource to create a tse update_group_spec

Example Usage

```hcl
resource "tencentcloud_tse_update_group_spec" "update_group_spec" {
  gateway_id = ""
  group_id = ""
  node_config {
		specification = ""
		number =

  }
}
```

Import

tse update_group_spec can be imported using the id, e.g.

```
terraform import tencentcloud_tse_update_group_spec.update_group_spec update_group_spec_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudTseUpdateGroupSpec() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseUpdateGroupSpecCreate,
		Read:   resourceTencentCloudTseUpdateGroupSpecRead,
		Update: resourceTencentCloudTseUpdateGroupSpecUpdate,
		Delete: resourceTencentCloudTseUpdateGroupSpecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway IDonly postpaid gateway supported.",
			},

			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway group ID.",
			},

			"node_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Group node config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Group specification, 1c2g|2c4g|4c8g|8c16g.",
						},
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Group node number, 2-50.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseUpdateGroupSpecCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_update_group_spec.create")()
	defer inconsistentCheck(d, meta)()

	var groupId string
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(groupId)

	return resourceTencentCloudTseUpdateGroupSpecUpdate(d, meta)
}

func resourceTencentCloudTseUpdateGroupSpecRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_update_group_spec.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	updateGroupSpecId := d.Id()

	updateGroupSpec, err := service.DescribeTseUpdateGroupSpecById(ctx, groupId)
	if err != nil {
		return err
	}

	if updateGroupSpec == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseUpdateGroupSpec` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if updateGroupSpec.GatewayId != nil {
		_ = d.Set("gateway_id", updateGroupSpec.GatewayId)
	}

	if updateGroupSpec.GroupId != nil {
		_ = d.Set("group_id", updateGroupSpec.GroupId)
	}

	if updateGroupSpec.NodeConfig != nil {
		nodeConfigMap := map[string]interface{}{}

		if updateGroupSpec.NodeConfig.Specification != nil {
			nodeConfigMap["specification"] = updateGroupSpec.NodeConfig.Specification
		}

		if updateGroupSpec.NodeConfig.Number != nil {
			nodeConfigMap["number"] = updateGroupSpec.NodeConfig.Number
		}

		_ = d.Set("node_config", []interface{}{nodeConfigMap})
	}

	return nil
}

func resourceTencentCloudTseUpdateGroupSpecUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_update_group_spec.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewUpdateCloudNativeAPIGatewaySpecRequest()

	updateGroupSpecId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"gateway_id", "group_id", "node_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("node_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
			cloudNativeAPIGatewayNodeConfig := tse.CloudNativeAPIGatewayNodeConfig{}
			if v, ok := dMap["specification"]; ok {
				cloudNativeAPIGatewayNodeConfig.Specification = helper.String(v.(string))
			}
			if v, ok := dMap["number"]; ok {
				cloudNativeAPIGatewayNodeConfig.Number = helper.IntInt64(v.(int))
			}
			request.NodeConfig = &cloudNativeAPIGatewayNodeConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().UpdateCloudNativeAPIGatewaySpec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse updateGroupSpec failed, reason:%+v", logId, err)
		return err
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 5*readRetryTimeout, time.Second, service.TseUpdateGroupSpecStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTseUpdateGroupSpecRead(d, meta)
}

func resourceTencentCloudTseUpdateGroupSpecDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_update_group_spec.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
