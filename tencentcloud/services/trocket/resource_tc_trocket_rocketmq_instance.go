package trocket

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

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

func ResourceTencentCloudTrocketRocketmqInstance() *schema.Resource {
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
				Description: "SKU code. Available specifications are as follows: experiment_500, basic_1k, basic_2k, basic_3k, basic_4k, basic_5k, basic_6k, basic_7k, basic_8k, basic_9k, basic_10k, pro_4k, pro_6k, pro_8k, pro_1w, pro_15k, pro_2w, pro_25k, pro_3w, pro_35k, pro_4w, pro_45k, pro_5w, pro_55k, pro_60k, pro_65k, pro_70k, pro_75k, pro_80k, pro_85k, pro_90k, pro_95k, pro_100k, platinum_1w, platinum_2w, platinum_3w, platinum_4w, platinum_5w, platinum_6w, platinum_7w, platinum_8w, platinum_9w, platinum_10w, platinum_12w, platinum_14w, platinum_16w, platinum_18w, platinum_20w, platinum_25w, platinum_30w, platinum_35w, platinum_40w, platinum_45w, platinum_50w, platinum_60w, platinum_70w, platinum_80w, platinum_90w, platinum_100w.",
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
				Description: "Whether to enable the public network. Must set `bandwidth` when `enable_public` equal true.",
			},

			"bandwidth": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Public network bandwidth. `bandwidth` must be greater than zero when `enable_public` equal true.",
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

			"public_end_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public network access address.",
			},

			"vpc_end_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC access address.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request      = trocket.NewCreateInstanceRequest()
		response     = trocket.NewCreateInstanceResponse()
		instanceId   string
		enablePublic bool
		bandwidth    int
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
		enablePublic = v.(bool)
		request.EnablePublic = helper.Bool(enablePublic)
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		bandwidth = v.(int)
		request.Bandwidth = helper.IntInt64(bandwidth)
	}

	if enablePublic && bandwidth <= 0 {
		return fmt.Errorf("`bandwidth` must be greater than zero when `enable_public` equal true.")
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().CreateInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create trocket rocketmqInstance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqInstance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	// wait
	service := TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"RUNNING"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::trocket:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqInstanceRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	rocketmqInstance, err := service.DescribeTrocketRocketmqInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if rocketmqInstance == nil {
		log.Printf("[WARN]%s resource `TrocketRocketmqInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
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

	var enablePublic bool
	for _, endpoint := range rocketmqInstance.EndpointList {
		endpointType := endpoint.Type
		if endpointType == nil {
			continue
		}

		if *endpointType == ENDPOINT_TYPE_PUBLIC {
			enablePublic = true
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

			_ = d.Set("public_end_point", endpoint.EndpointUrl)
		}

		if *endpointType == ENDPOINT_TYPE_VPC {
			if endpoint.VpcId != nil {
				_ = d.Set("vpc_id", endpoint.VpcId)
			}

			if endpoint.SubnetId != nil {
				_ = d.Set("subnet_id", endpoint.SubnetId)
			}

			_ = d.Set("vpc_end_point", endpoint.EndpointUrl)
		}

	}

	_ = d.Set("enable_public", enablePublic)

	if rocketmqInstance.MessageRetention != nil {
		_ = d.Set("message_retention", rocketmqInstance.MessageRetention)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tags, err := tagService.DescribeResourceTags(ctx, "trocket", "instance", tcClient.Region, instanceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTrocketRocketmqInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                      = tccommon.GetLogId(tccommon.ContextNil)
		ctx                        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request1                   = trocket.NewModifyInstanceRequest()
		request2                   = trocket.NewModifyInstanceEndpointRequest()
		instanceId                 = d.Id()
		needModifyInstance         bool
		needModifyInstanceEndpoint bool
	)

	immutableArgs := []string{"instance_type", "vpc_id", "subnet_id", "enable_public"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request1.Name = helper.String(v.(string))
		}

		needModifyInstance = true
	}

	if d.HasChange("sku_code") {
		if v, ok := d.GetOk("sku_code"); ok {
			request1.SkuCode = helper.String(v.(string))
		}

		needModifyInstance = true
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request1.Remark = helper.String(v.(string))
		}

		needModifyInstance = true
	}

	if d.HasChange("message_retention") {
		if v, ok := d.GetOkExists("message_retention"); ok {
			request1.MessageRetention = helper.IntInt64(v.(int))
		}

		needModifyInstance = true
	}

	if needModifyInstance {
		request1.InstanceId = &instanceId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().ModifyInstance(request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update trocket rocketmqInstance failed, reason:%+v", logId, err)
			return err
		}

		// wait
		service := TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"RUNNING"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("bandwidth") {
		if v, ok := d.GetOkExists("bandwidth"); ok {
			request2.Bandwidth = helper.IntInt64(v.(int))
		}

		needModifyInstanceEndpoint = true
	}

	if d.HasChange("ip_rules") {
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

				request2.IpRules = append(request2.IpRules, &ipRule)
			}
		}

		needModifyInstanceEndpoint = true
	}

	if needModifyInstanceEndpoint {
		if v, ok := d.GetOkExists("enable_public"); ok {
			if v.(bool) {
				request2.InstanceId = &instanceId
				request2.Type = helper.String("PUBLIC")
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().ModifyInstanceEndpoint(request2)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s update trocket rocketmqInstance failed, reason:%+v", logId, err)
					return err
				}

				// wait
				service := TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
				conf := tccommon.BuildStateChangeConf([]string{}, []string{"RUNNING"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(instanceId, []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
				}
			} else {
				return fmt.Errorf("Only instances with public network access can modify `bandwidth` or `ip_rules`.")
			}
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("trocket", "instance", tcClient.Region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTrocketRocketmqInstanceRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.DeleteTrocketRocketmqInstanceById(ctx, instanceId); err != nil {
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{""}, 10*tccommon.ReadRetryTimeout, time.Second, service.TrocketRocketmqInstanceStateRefreshFunc(d.Id(), []string{}))
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
