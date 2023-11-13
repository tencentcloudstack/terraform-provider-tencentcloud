/*
Provides a resource to create a redis modify_network_config

Example Usage

```hcl
resource "tencentcloud_redis_modify_network_config" "modify_network_config" {
  instance_id = "crs-c1nl9rpv"
  operation = "changeVip"
  vip = "10.1.1.2"
  vpc_id = "vpc-hu6khgap"
  subnet_id = "subnet-6mt7lir6"
  recycle = 7
  v_port = 6379
}
```

Import

redis modify_network_config can be imported using the id, e.g.

```
terraform import tencentcloud_redis_modify_network_config.modify_network_config modify_network_config_id
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

func resourceTencentCloudRedisModifyNetworkConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisModifyNetworkConfigCreate,
		Read:   resourceTencentCloudRedisModifyNetworkConfigRead,
		Update: resourceTencentCloudRedisModifyNetworkConfigUpdate,
		Delete: resourceTencentCloudRedisModifyNetworkConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"operation": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Refers to the category of pre-modified networks, including:- changeVip: refers to switching a VPC, including its intranet IPv4 address and port.- changeVpc: Switches the subnet to which the VPC belongs.- changeBaseToVpc: refers to the switch of the basic network to the private network.- changeVPort: refers to modifying only the network port of the instance.",
			},

			"vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Refers to the IPv4 address of the private network of the instance.If the Operation parameter is changeVip, you need to configure this parameter.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Refers to the modified VPC ID, which must be configured when the Operation parameter is changeVpc or changeBaseToVpc.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Refers to the subnet ID of the modified VPC, which must be configured if the Operation parameter is changeVpc or changeBaseToVpc.",
			},

			"recycle": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "How long the IPv4 address of the original private network is retained.- Unit: days.- Value range: 0, 1, 2, 3, 7, 15.Note: To set the original address retention period, the latest version of the SDK is required, otherwise the original address will be released immediately, see [SDK Center](https://cloud.tencent.com/document/sdk)。.",
			},

			"v_port": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Refers to the modified network port.If the Operation parameter is changeVPort or changeVip, you need to configure this parameter.The value range is [1024,65535].",
			},
		},
	}
}

func resourceTencentCloudRedisModifyNetworkConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modify_network_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisModifyNetworkConfigUpdate(d, meta)
}

func resourceTencentCloudRedisModifyNetworkConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modify_network_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	modifyNetworkConfigId := d.Id()

	modifyNetworkConfig, err := service.DescribeRedisModifyNetworkConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if modifyNetworkConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisModifyNetworkConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if modifyNetworkConfig.InstanceId != nil {
		_ = d.Set("instance_id", modifyNetworkConfig.InstanceId)
	}

	if modifyNetworkConfig.Operation != nil {
		_ = d.Set("operation", modifyNetworkConfig.Operation)
	}

	if modifyNetworkConfig.Vip != nil {
		_ = d.Set("vip", modifyNetworkConfig.Vip)
	}

	if modifyNetworkConfig.VpcId != nil {
		_ = d.Set("vpc_id", modifyNetworkConfig.VpcId)
	}

	if modifyNetworkConfig.SubnetId != nil {
		_ = d.Set("subnet_id", modifyNetworkConfig.SubnetId)
	}

	if modifyNetworkConfig.Recycle != nil {
		_ = d.Set("recycle", modifyNetworkConfig.Recycle)
	}

	if modifyNetworkConfig.VPort != nil {
		_ = d.Set("v_port", modifyNetworkConfig.VPort)
	}

	return nil
}

func resourceTencentCloudRedisModifyNetworkConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modify_network_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyNetworkConfigRequest()

	modifyNetworkConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "operation", "vip", "vpc_id", "subnet_id", "recycle", "v_port"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyNetworkConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis modifyNetworkConfig failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"succeed"}, 60*readRetryTimeout, time.Second, service.RedisModifyNetworkConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudRedisModifyNetworkConfigRead(d, meta)
}

func resourceTencentCloudRedisModifyNetworkConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_modify_network_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
