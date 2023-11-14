/*
Provides a resource to create a redis param

Example Usage

```hcl
resource "tencentcloud_redis_param" "param" {
  instance_id = "crs-c1nl9rpv"
  instance_params {
		key = &lt;nil&gt;
		value = &lt;nil&gt;

  }
}
```

Import

redis param can be imported using the id, e.g.

```
terraform import tencentcloud_redis_param.param param_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"log"
	"time"
)

func resourceTencentCloudRedisParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisParamCreate,
		Read:   resourceTencentCloudRedisParamRead,
		Update: resourceTencentCloudRedisParamUpdate,
		Delete: resourceTencentCloudRedisParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"instance_params": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "A list of parameters modified by the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Sets the name of the parameter.For example, timeout.For more information about custom parameters, see(https://cloud.tencent.com/document/product/239/49925).",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Set the run value for the parameter name.For example, the corresponding running value of timeout can be set to 120 in seconds (s).Refers to the close of the client connection when the idle time reaches 120 s.For more information about parameter values, see(https://cloud.tencent.com/document/product/239/49925).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudRedisParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisParamUpdate(d, meta)
}

func resourceTencentCloudRedisParamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	paramId := d.Id()

	param, err := service.DescribeRedisParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if param == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if param.InstanceId != nil {
		_ = d.Set("instance_id", param.InstanceId)
	}

	if param.InstanceParams != nil {
		instanceParamsList := []interface{}{}
		for _, instanceParams := range param.InstanceParams {
			instanceParamsMap := map[string]interface{}{}

			if param.InstanceParams.Key != nil {
				instanceParamsMap["key"] = param.InstanceParams.Key
			}

			if param.InstanceParams.Value != nil {
				instanceParamsMap["value"] = param.InstanceParams.Value
			}

			instanceParamsList = append(instanceParamsList, instanceParamsMap)
		}

		_ = d.Set("instance_params", instanceParamsList)

	}

	return nil
}

func resourceTencentCloudRedisParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyInstanceParamsRequest()

	paramId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "instance_params"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyInstanceParams(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis param failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 30*readRetryTimeout, time.Second, service.RedisParamStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisParamRead(d, meta)
}

func resourceTencentCloudRedisParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
