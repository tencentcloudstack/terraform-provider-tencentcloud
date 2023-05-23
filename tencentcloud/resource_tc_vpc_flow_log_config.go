/*
Provides a resource to create a vpc flow_log_config

Example Usage

```hcl
resource "tencentcloud_vpc_flow_log_config" "flow_log_config" {
  flow_log_id = "fl-geg2keoj"
  enable = false
}
```

Import

vpc flow_log_config can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_flow_log_config.flow_log_config flow_log_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudVpcFlowLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcFlowLogConfigCreate,
		Read:   resourceTencentCloudVpcFlowLogConfigRead,
		Update: resourceTencentCloudVpcFlowLogConfigUpdate,
		Delete: resourceTencentCloudVpcFlowLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_log_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow log ID.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "If enable snapshot policy.",
			},
		},
	}
}

func resourceTencentCloudVpcFlowLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log_config.create")()
	defer inconsistentCheck(d, meta)()

	flowLogId := d.Get("flow_log_id").(string)

	d.SetId(flowLogId)

	return resourceTencentCloudVpcFlowLogConfigUpdate(d, meta)
}

func resourceTencentCloudVpcFlowLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	flowLogId := d.Id()

	request := vpc.NewDescribeFlowLogsRequest()
	request.FlowLogId = &flowLogId

	flowLogs, err := service.DescribeFlowLogs(ctx, request)
	if err != nil {
		return err
	}

	if len(flowLogs) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcFlowLogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("flow_log_id", flowLogId)

	flowLogConfig := flowLogs[0]

	if flowLogConfig.Enable != nil {
		_ = d.Set("enable", flowLogConfig.Enable)
	}

	return nil
}

func resourceTencentCloudVpcFlowLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enable         bool
		enableRequest  = vpc.NewEnableFlowLogsRequest()
		disableRequest = vpc.NewDisableFlowLogsRequest()
	)

	flowLogId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.FlowLogIds = []*string{&flowLogId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().EnableFlowLogs(enableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc flowLogConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.FlowLogIds = []*string{&flowLogId}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DisableFlowLogs(disableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc flowLogConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudVpcFlowLogConfigRead(d, meta)
}

func resourceTencentCloudVpcFlowLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_flow_log_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
