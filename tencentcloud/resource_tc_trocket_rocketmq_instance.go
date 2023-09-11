/*
Provides a resource to create a rocketmq 5.x instance

Example Usage

```hcl
resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name = "test"
  sku_code = "experiment_500"
  remark = "test"
  vpc_id = "vpc-qmvl8z4f"
  subnet_id = "subnet-ncef9v74"
  tags = {
    tag_key = "rocketmq"
    tag_value = "5.x"
  }
}
```

Import

trocket rocketmq_instance can be imported using the id, e.g.

```
terraform import tencentcloud_trocket_rocketmq_instance.rocketmq_instance rocketmq_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTrocketRocketmqInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTrocketRocketmqInstanceCreate,
		Read:   resourceTencentCloudTrocketRocketmqInstanceRead,
		Update: resourceTencentCloudTrocketRocketmqInstanceUpdate,
		Delete: resourceTencentCloudTrocketRocketmqInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance type. Valid values: `EXPERIMENT`, `BASIC`, `PRO`, `PLATINUM`.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"sku_code": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SKU code. Available specifications are as follows: experiment_500, basic_1k, basic_2k, basic_4k, basic_6k.",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Remark.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tag description list.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC id.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id.",
			},

			"enable_public": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the public network.",
			},

			"bandwidth": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Public network bandwidth.",
			},

			"ip_rules": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Public network access whitelist.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP.",
						},
						"allow": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to allow release or not.",
						},
						"remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remark.",
						},
					},
				},
			},

			"message_retention": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Message retention time in hours.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		request    = trocket.NewCreateInstanceRequest()
		response   = trocket.NewCreateInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sku_code"); ok {
		request.SkuCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	vpcInfo := trocket.VpcInfo{
		VpcId:    helper.String(d.Get("vpc_id").(string)),
		SubnetId: helper.String(d.Get("subnet_id").(string)),
	}

	request.VpcList = []*trocket.VpcInfo{&vpcInfo}

	if v, ok := d.GetOkExists("enable_public"); ok {
		request.EnablePublic = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ip_rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ipRule := trocket.IpRule{}
			if v, ok := dMap["ip"]; ok {
				ipRule.Ip = helper.String(v.(string))
			}
			if v, ok := dMap["allow"]; ok {
				ipRule.Allow = helper.Bool(v.(bool))
			}
			if v, ok := dMap["remark"]; ok {
				ipRule.Remark = helper.String(v.(string))
			}
			request.IpRules = append(request.IpRules, &ipRule)
		}
	}

	if v, ok := d.GetOkExists("message_retention"); ok {
		request.MessageRetention = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().CreateInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqInstance failed, reason:%+v", logId, err)
		return err
	}
	instanceId = *response.Response.InstanceId

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 10*readRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::trocket:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudTrocketRocketmqInstanceRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	rocketmqInstance, err := service.DescribeTrocketRocketmqInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rocketmqInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TrocketRocketmqInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rocketmqInstance.InstanceType != nil {
		_ = d.Set("instance_type", rocketmqInstance.InstanceType)
	}

	if rocketmqInstance.InstanceName != nil {
		_ = d.Set("name", rocketmqInstance.InstanceName)
	}

	if rocketmqInstance.SkuCode != nil {
		_ = d.Set("sku_code", rocketmqInstance.SkuCode)
	}

	if rocketmqInstance.Remark != nil {
		_ = d.Set("remark", rocketmqInstance.Remark)
	}

	if len(rocketmqInstance.EndpointList) != 0 {
		endpoint := rocketmqInstance.EndpointList[0]
		if endpoint.VpcId != nil {
			_ = d.Set("vpc_id", endpoint.VpcId)
		}

		if endpoint.SubnetId != nil {
			_ = d.Set("subnet_id", endpoint.SubnetId)
		}

		if len(endpoint.IpRules) != 0 {
			ipRuleList := make([]interface{}, 0)
			for _, ipRule := range endpoint.IpRules {
				ipRuleMap := make(map[string]interface{})
				ipRuleMap["ip"] = ipRule.Ip
				ipRuleMap["allow"] = ipRule.Allow
				ipRuleMap["remark"] = ipRule.Remark
				ipRuleList = append(ipRuleList, ipRuleMap)
			}
			_ = d.Set("ip_rules", ipRuleList)
		}
		if endpoint.Bandwidth != nil {
			_ = d.Set("bandwidth", endpoint.Bandwidth)
		}

		if endpoint.Type != nil {
			if *endpoint.Type == "PUBLIC" {
				_ = d.Set("enable_public", true)
			} else {
				_ = d.Set("enable_public", true)
			}
		}

	}

	if rocketmqInstance.MessageRetention != nil {
		_ = d.Set("message_retention", rocketmqInstance.MessageRetention)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	tags, err := tagService.DescribeResourceTags(ctx, "trocket", "instance", tcClient.Region, instanceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTrocketRocketmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := trocket.NewModifyInstanceRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_type", "vpc_id", "subnet_id", "enable_public", "bandwidth", "ip_rules"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("sku_code") {
		if v, ok := d.GetOk("sku_code"); ok {
			request.SkuCode = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("message_retention") {
		if v, ok := d.GetOkExists("message_retention"); ok {
			request.MessageRetention = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTrocketClient().ModifyInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update trocket rocketmqInstance failed, reason:%+v", logId, err)
		return err
	}

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf([]string{}, []string{"RUNNING"}, 10*readRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("trocket", "instance", tcClient.Region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}
	return resourceTencentCloudTrocketRocketmqInstanceRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_trocket_rocketmq_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TrocketService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteTrocketRocketmqInstanceById(ctx, instanceId); err != nil {
		return err
	}
	conf := BuildStateChangeConf([]string{}, []string{""}, 10*readRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err := conf.WaitForState(); err != nil {
		if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkerr.Code == "ResourceNotFound.Instance" {
				return nil
			}
		}
		return err
	}
	return nil
}
