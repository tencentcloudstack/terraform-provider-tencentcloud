/*
Use this resource to create ckafka instance.

~> **NOTE:** It only support create profession ckafka instance.

Example Usage

```hcl
resource "tencentcloud_ckafka_instance" "foo" {
  band_width         = 40
  disk_size          = 500
  disk_type          = "CLOUD_BASIC"
  period             = 1
  instance_name      = "ckafka-instance-tf-test"
  kafka_version      = "1.1.1"
  msg_retention_time = 1300
  multi_zone_flag    = true
  partition          = 800
  public_network     = 3
  renew_flag         = 0
  subnet_id          = "subnet-4vwihrzk"
  vpc_id             = "vpc-82p1t1nv"
  zone_id            = 100006
  zone_ids           = [
    100006,
    100007,
  ]

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    bottom_retention        = 0
    disk_quota_percentage   = 0
    enable                  = 1
    step_forward_percentage = 0
  }
}
```

Import

ckafka instance can be imported using the instance_id, e.g.

```
$ terraform import tencentcloud_ckafka_instance.foo ckafka-f9ife4zz
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCkafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaInstanceCreate,
		Read:   resourceTencentCloudCkafkaInstanceRead,
		Update: resourceTencentCloudCkafkaInstanceUpdate,
		Delete: resourceTencentCLoudCkafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Available zone id.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					multiZone := d.Get("multi_zone_flag").(bool)
					zoneId := d.Get("zone_id").(int)
					v, ok := d.GetOk("zone_ids")

					if !multiZone || !ok || old == "" {
						return old == new
					}

					zoneIds := v.(*schema.Set)
					return zoneIds.Contains(zoneId)
				},
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Prepaid purchase time, such as 1, is one month.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Vpc id.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet id.",
			},
			"msg_retention_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "The maximum retention time of instance logs, in minutes." +
					" the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled," +
					" which means that the log retention time recovery policy is not enabled.",
			},
			"renew_flag": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Prepaid automatic renewal mark, 0 means the default state, the initial state," +
					" 1 means automatic renewal, 2 means clear no automatic renewal (user setting).",
			},
			"kafka_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Kafka version (0.10.2/1.1.1/2.4.1).",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Instance bandwidth in MBps.",
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Disk Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. " +
					"If it is not within the interval, the plan will cause a change when first created.",
			},
			"partition": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Partition Size. Its interval varies with bandwidth, and the input must be within the interval, which can be viewed through the control. " +
					"If it is not within the interval, the plan will cause a change when first created.",
			},
			"multi_zone_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the instance is multi zones. NOTE: if set to `true`, `zone_ids` must set together.",
			},
			"zone_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "List of available zone id. NOTE: this argument must set together with `multi_zone_flag`.",
				RequiredWith: []string{"multi_zone_flag"},
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Partition size, the professional version does not need tag.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Type of disk.",
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_topic_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Automatic creation. true: enabled, false: not enabled.",
						},
						"default_num_partitions": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "If auto.create.topic.enable is set to true and this value is not set, " +
								"3 will be used by default.",
						},
						"default_replication_factor": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "If auto.create.topic.enable is set to true but this value is not set, " +
								"2 will be used by default.",
						},
					},
				},
				Description: "Instance configuration.",
			},
			"dynamic_retention_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Whether the dynamic message retention time configuration is " +
								"enabled. 0: disabled; 1: enabled.",
						},
						"disk_quota_percentage": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Disk quota threshold (in percentage) for triggering " +
								"the message retention time change event.",
						},
						"step_forward_percentage": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Percentage by which the message retention " +
								"time is shortened each time.",
						},
						"bottom_retention": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Minimum retention time, in minutes.",
						},
					},
				},
				Description: "Dynamic message retention policy configuration.",
			},
			"rebalance_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Modification of the rebalancing time after upgrade.",
			},
			"public_network": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp.",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vip of instance.",
			},
			"vport": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of instance.",
			},
		},
	}
}

func resourceTencentCloudCkafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_instance.create")()
	var (
		logId   = getLogId(contextNil)
		service = CkafkaService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		request  = ckafka.NewCreateInstancePreRequest()
		response = ckafka.NewCreateInstancePreResponse()
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
	)
	instanceName := d.Get("instance_name").(string)
	request.InstanceName = &instanceName

	zoneId := int64(d.Get("zone_id").(int))
	request.ZoneId = &zoneId

	period := int64(d.Get("period").(int))
	request.Period = helper.String(fmt.Sprintf("%dm", period))
	// only support create profession instance
	request.InstanceType = helper.Int64(1)
	request.SpecificationsType = helper.String("profession")

	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId := v.(string)
		request.VpcId = helper.String(vpcId)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId := v.(string)
		request.SubnetId = helper.String(subnetId)
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		renewFlag := int64(v.(int))
		request.RenewFlag = helper.Int64(renewFlag)
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		kafkaVersion := v.(string)
		request.KafkaVersion = helper.String(kafkaVersion)
	}

	if v, ok := d.GetOk("disk_size"); ok {
		diskSize := int64(v.(int))
		request.DiskSize = helper.Int64(diskSize)
	}

	if v, ok := d.GetOk("band_width"); ok {
		bandWidth := int64(v.(int))
		request.BandWidth = helper.Int64(bandWidth)
	}

	if v, ok := d.GetOk("partition"); ok {
		partition := int64(v.(int))
		request.Partition = helper.Int64(partition)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagSet := make([]*ckafka.Tag, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := ckafka.Tag{
				TagKey:   helper.String(m["key"].(string)),
				TagValue: helper.String(m["value"].(string)),
			}
			tagSet = append(tagSet, &tagInfo)
		}
		request.Tags = tagSet
	}

	if v, ok := d.GetOk("disk_type"); ok {
		diskType := v.(string)
		request.DiskType = helper.String(diskType)
	}

	if flag := d.Get("multi_zone_flag").(bool); flag {
		request.MultiZoneFlag = helper.Bool(flag)
		ids := d.Get("zone_ids").(*schema.Set).List()
		for _, v := range ids {
			request.ZoneIds = append(request.ZoneIds, helper.IntInt64(v.(int)))
		}
	}

	result, err := service.client.UseCkafkaClient().CreateInstancePre(request)
	response = result

	if err != nil {
		log.Printf("[CRITAL]%s create ckafka instance failed, reason:%s\n", logId, err.Error())
		return err
	}

	instanceId := response.Response.Result.Data.InstanceId

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		has, ready, err := service.CheckCkafkaInstanceReady(ctx, *instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("ckafka instance not exists."))
		}
		if ready {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("create ckafka instance task is processing"))
	})
	if err != nil {
		return err
	}
	d.SetId(*instanceId)

	// modify instance attributes
	var (
		needModify    = false
		modifyRequest = ckafka.NewModifyInstanceAttributesRequest()
	)
	modifyRequest.InstanceId = instanceId

	if v, ok := d.GetOk("msg_retention_time"); ok {
		retentionTime := int64(v.(int))
		modifyRequest.MsgRetentionTime = helper.Int64(retentionTime)
	}

	if v, ok := d.GetOk("config"); ok {
		needModify = true
		config := make([]*ckafka.ModifyInstanceAttributesConfig, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			configInfo := ckafka.ModifyInstanceAttributesConfig{}
			if autoCreateTopicEnable, ok := dMap["auto_create_topic_enable"]; ok {
				configInfo.AutoCreateTopicEnable = helper.Bool(autoCreateTopicEnable.(bool))
			}
			if defaultNumPartitions, ok := dMap["default_num_partitions"]; ok {
				configInfo.DefaultNumPartitions = helper.Int64(int64(defaultNumPartitions.(int)))
			}
			if defaultReplicationFactor, ok := dMap["default_replication_factor"]; ok {
				configInfo.DefaultReplicationFactor = helper.Int64(int64(defaultReplicationFactor.(int)))
			}
			config = append(config, &configInfo)
		}
		modifyRequest.Config = config[0]
	}

	if v, ok := d.GetOk("dynamic_retention_config"); ok {
		needModify = true
		dynamic := make([]*ckafka.DynamicRetentionTime, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dynamicInfo := ckafka.DynamicRetentionTime{}
			if enable, ok := dMap["enable"]; ok {
				dynamicInfo.Enable = helper.Int64(int64(enable.(int)))
			}
			if diskQuotaPercentage, ok := dMap["disk_quota_percentage"]; ok {
				dynamicInfo.DiskQuotaPercentage = helper.Int64(int64(diskQuotaPercentage.(int)))
			}
			if stepForwardPercentage, ok := dMap["step_forward_percentage"]; ok {
				dynamicInfo.StepForwardPercentage = helper.Int64(int64(stepForwardPercentage.(int)))
			}
			if bottomRetention, ok := dMap["bottom_retention"]; ok {
				dynamicInfo.BottomRetention = helper.Int64(int64(bottomRetention.(int)))
			}
			dynamic = append(dynamic, &dynamicInfo)
		}
		modifyRequest.DynamicRetentionConfig = dynamic[0]
	}

	if v, ok := d.GetOk("rebalance_time"); ok {
		needModify = true
		modifyRequest.RebalanceTime = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("public_network"); ok {
		needModify = true
		modifyRequest.PublicNetwork = helper.Int64(int64(v.(int)))
	}

	if needModify {
		error := service.ModifyCkafkaInstanceAttributes(ctx, modifyRequest)
		if error != nil {
			return fmt.Errorf("[API]Set kafka instance attributes fail, reason:%s", error.Error())
		}
	}
	return resourceTencentCloudCkafkaInstanceRead(d, meta)
}

func resourceTencentCloudCkafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var service = CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId := d.Id()

	var info *ckafka.InstanceDetail

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		res, has, e := service.DescribeCkafkaInstanceById(ctx, instanceId)
		if e != nil {
			return retryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}
		info = res
		return nil
	})
	if err != nil {
		return fmt.Errorf("[API]Describe kafka instance fail, reason:%s", err.Error())
	}
	_ = d.Set("instance_name", info.InstanceName)
	_ = d.Set("zone_id", info.ZoneId)
	// calculate period
	//createTime := *info.CreateTime
	//expireTime := *info.ExpireTime
	//period := (expireTime - createTime) / (3600 * 24 * 31)
	//_ = d.Set("period", &period)
	_ = d.Set("vpc_id", info.VpcId)
	_ = d.Set("subnet_id", info.SubnetId)
	_ = d.Set("renew_flag", info.RenewFlag)
	_ = d.Set("kafka_version", info.Version)
	_ = d.Set("disk_size", info.DiskSize)
	bandWidth := info.Bandwidth
	_ = d.Set("vip", info.Vip)
	_ = d.Set("vport", info.Vport)
	_ = d.Set("band_width", *bandWidth/8)
	_ = d.Set("partition", info.MaxPartitionNumber)

	if len(info.ZoneIds) > 0 {
		_ = d.Set("multi_zone_flag", true)
		ids := helper.Int64sInterfaces(info.ZoneIds)
		idSet := schema.NewSet(func(i interface{}) int {
			return i.(int)
		}, ids)
		_ = d.Set("zone_ids", idSet)
	}

	tagSets := make([]map[string]interface{}, 0, len(info.Tags))
	for _, item := range info.Tags {
		tagSets = append(tagSets, map[string]interface{}{
			"key":   item.TagKey,
			"value": item.TagValue,
		})
	}
	_ = d.Set("tags", tagSets)
	_ = d.Set("disk_type", info.DiskType)
	_ = d.Set("rebalance_time", info.RebalanceTime)

	// query msg_retention_time
	var (
		request  = ckafka.NewDescribeInstanceAttributesRequest()
		response = ckafka.NewDescribeInstanceAttributesResponse()
	)
	request.InstanceId = &instanceId
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.client.UseCkafkaClient().DescribeInstanceAttributes(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		attr := response.Response.Result
		_ = d.Set("msg_retention_time", attr.MsgRetentionTime)

		if attr.Config != nil {
			config := make([]map[string]interface{}, 0)
			config = append(config, map[string]interface{}{
				"auto_create_topic_enable":   attr.Config.AutoCreateTopicsEnable,
				"default_num_partitions":     attr.Config.DefaultNumPartitions,
				"default_replication_factor": attr.Config.DefaultReplicationFactor,
			})
			_ = d.Set("config", config)
		}

		dynamicConfig := make([]map[string]interface{}, 0)
		dynamicConfig = append(dynamicConfig, map[string]interface{}{
			"enable":                  attr.RetentionTimeConfig.Enable,
			"disk_quota_percentage":   attr.RetentionTimeConfig.DiskQuotaPercentage,
			"step_forward_percentage": attr.RetentionTimeConfig.StepForwardPercentage,
			"bottom_retention":        attr.RetentionTimeConfig.BottomRetention,
		})
		_ = d.Set("dynamic_retention_config", dynamicConfig)
		_ = d.Set("public_network", attr.PublicNetwork)

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Ckafka Instance Attributes failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudCkafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_instance.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if d.HasChange("zone_id") || d.HasChange("period") || d.HasChange("vpc_id") || d.HasChange("subnet_id") ||
		d.HasChange("renew_flag") || d.HasChange("kafka_version") || d.HasChange("multi_zone_flag") || d.HasChange("zone_ids") ||
		d.HasChange("tags") || d.HasChange("disk_type") {

		return fmt.Errorf("parms like 'zone_id | period | vpc_id | subnet_id | renew_flag | " +
			"kafka_version | multi_zone_flag | zone_ids | tags | disk_type', do not support change now.")
	}

	instanceId := d.Id()
	request := ckafka.NewModifyInstanceAttributesRequest()
	request.InstanceId = &instanceId
	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	}

	if d.HasChange("msg_retention_time") {
		if v, ok := d.GetOk("msg_retention_time"); ok {
			request.MsgRetentionTime = helper.Int64(int64(v.(int)))
		}
	}

	if d.HasChange("config") {
		if v, ok := d.GetOk("config"); ok {
			items := v.([]interface{})
			dMap := items[0].(map[string]interface{})
			configInfo := ckafka.ModifyInstanceAttributesConfig{}
			if autoCreateTopicEnable, ok := dMap["auto_create_topic_enable"]; ok {
				configInfo.AutoCreateTopicEnable = helper.Bool(autoCreateTopicEnable.(bool))
			}
			if defaultNumPartitions, ok := dMap["default_num_partitions"]; ok {
				configInfo.DefaultNumPartitions = helper.Int64(int64(defaultNumPartitions.(int)))
			}
			if defaultReplicationFactor, ok := dMap["default_replication_factor"]; ok {
				configInfo.DefaultReplicationFactor = helper.Int64(int64(defaultReplicationFactor.(int)))
			}
			request.Config = &configInfo
		}
	}

	if d.HasChange("dynamic_retention_config") {
		if v, ok := d.GetOk("dynamic_retention_config"); ok {
			items := v.([]interface{})
			dMap := items[0].(map[string]interface{})
			dynamicInfo := ckafka.DynamicRetentionTime{}
			if enable, ok := dMap["enable"]; ok {
				dynamicInfo.Enable = helper.Int64(int64(enable.(int)))
			}
			if diskQuotaPercentage, ok := dMap["disk_quota_percentage"]; ok {
				dynamicInfo.DiskQuotaPercentage = helper.Int64(int64(diskQuotaPercentage.(int)))
			}
			if stepForwardPercentage, ok := dMap["step_forward_percentage"]; ok {
				dynamicInfo.StepForwardPercentage = helper.Int64(int64(stepForwardPercentage.(int)))
			}
			if bottomRetention, ok := dMap["bottom_retention"]; ok {
				dynamicInfo.BottomRetention = helper.Int64(int64(bottomRetention.(int)))
			}
			request.DynamicRetentionConfig = &dynamicInfo
		}
	}

	if d.HasChange("rebalance_time") {
		if v, ok := d.GetOk("rebalance_time"); ok {
			request.RebalanceTime = helper.Int64(int64(v.(int)))
		}
	}

	if d.HasChange("public_network") {
		if v, ok := d.GetOk("public_network"); ok {
			request.PublicNetwork = helper.Int64(int64(v.(int)))
		}
	}

	error := service.ModifyCkafkaInstanceAttributes(ctx, request)
	if error != nil {
		return fmt.Errorf("[API]Set kafka instance attributes fail, reason:%s", error.Error())
	}

	if d.HasChange("band_width") || d.HasChange("disk_size") || d.HasChange("partition") {
		request := ckafka.NewModifyInstancePreRequest()
		request.InstanceId = helper.String(instanceId)
		if v, ok := d.GetOk("band_width"); ok {
			request.BandWidth = helper.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("disk_size"); ok {
			request.DiskSize = helper.Int64(int64(v.(int)))
		}
		if v, ok := d.GetOk("partition"); ok {
			request.Partition = helper.Int64(int64(v.(int)))
		}

		_, err := service.client.UseCkafkaClient().ModifyInstancePre(request)
		if err != nil {
			return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId,
				request.GetAction(), request.ToJsonString(), err.Error())
		}

		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			_, ready, err := service.CheckCkafkaInstanceReady(ctx, instanceId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if ready {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("upgrade ckafka instance task is processing"))
		})

		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaInstanceRead(d, meta)
}

func resourceTencentCLoudCkafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_instance.delete")()
	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CkafkaService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		request = ckafka.NewDeleteInstancePreRequest()
	)
	instanceId := d.Id()
	request.InstanceId = &instanceId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := service.client.UseCkafkaClient().DeleteInstancePre(request)
		if err != nil {
			return retryError(err, "UnsupportedOperation")
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		has, _, err := service.CheckCkafkaInstanceReady(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete ckafka instance task is processing"))
	})
	if err != nil {
		return err
	}
	return nil
}
