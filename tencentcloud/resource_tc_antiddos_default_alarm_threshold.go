/*
Provides a resource to create a antiddos default alarm threshold

Example Usage

```hcl
resource "tencentcloud_antiddos_default_alarm_threshold" "default_alarm_threshold" {
  default_alarm_config {
	alarm_type = 1
	alarm_threshold = 2000
  }
  instance_type = "bgp"
}
```

Import

antiddos default_alarm_threshold can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold ${instanceType}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAntiddosDefaultAlarmThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosDefaultAlarmThresholdCreate,
		Read:   resourceTencentCloudAntiddosDefaultAlarmThresholdRead,
		Update: resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate,
		Delete: resourceTencentCloudAntiddosDefaultAlarmThresholdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"default_alarm_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Alarm threshold configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm threshold type, value [1 (incoming traffic alarm threshold) 2 (attack cleaning traffic alarm threshold)].",
						},
						"alarm_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm threshold, in Mbps, with a value of&gt;=0; When used as an input parameter, setting 0 will delete the alarm threshold configuration;.",
						},
					},
				},
			},

			"instance_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Product type, value [bgp (represents advanced defense package product) bgpip (represents advanced defense IP product)].",
			},
		},
	}
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_type").(string))

	return resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate(d, meta)
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceType := d.Id()

	defaultAlarmThreshold, err := service.DescribeAntiddosDefaultAlarmThresholdById(ctx, instanceType, 1)
	if err != nil {
		return err
	}

	if defaultAlarmThreshold == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDefaultAlarmThreshold` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if defaultAlarmThreshold != nil {
		defaultAlarmConfigMap := map[string]interface{}{}

		if defaultAlarmThreshold.AlarmType != nil {
			defaultAlarmConfigMap["alarm_type"] = defaultAlarmThreshold.AlarmType
		}

		if defaultAlarmThreshold.AlarmThreshold != nil {
			defaultAlarmConfigMap["alarm_threshold"] = defaultAlarmThreshold.AlarmThreshold
		}

		_ = d.Set("default_alarm_config", []interface{}{defaultAlarmConfigMap})
	}

	_ = d.Set("instance_type", instanceType)

	return nil
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := antiddos.NewCreateDefaultAlarmThresholdRequest()

	instanceType := d.Id()

	request.InstanceType = &instanceType

	if d.HasChange("default_alarm_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "default_alarm_config"); ok {
			defaultAlarmThreshold := antiddos.DefaultAlarmThreshold{}
			if v, ok := dMap["alarm_type"]; ok {
				defaultAlarmThreshold.AlarmType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["alarm_threshold"]; ok {
				defaultAlarmThreshold.AlarmThreshold = helper.IntUint64(v.(int))
			}
			request.DefaultAlarmConfig = &defaultAlarmThreshold
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreateDefaultAlarmThreshold(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos defaultAlarmThreshold failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosDefaultAlarmThresholdRead(d, meta)
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
