/*
Use this resource to create ckafka topic.

Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "foo" {
	instance_id                     = "ckafka-f9ife4zz"
	topic_name                      = "example"
	note                            = "topic note"
	replica_num                     = 2
	partition_num                   = 1
	enable_white_list               = true
	ip_white_list                   = ["ip1","ip2"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 3600000
	retention                       = 60000
	max_message_bytes               = 0
}
```

Import

ckafka topic can be imported using the instance_id#topic_name, e.g.

```
$ terraform import tencentcloud_ckafka_topic.foo ckafka-f9ife4zz#example
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

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
				ForceNew:    true,
				Description: "Instance name.",
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Available zone id.",
			},
			"period": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 24),
				Description:  "Prepaid purchase time, such as '1m', is one month.",
			},
			"instance_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 3),
				Description: "Instance specifications, professional version fills in by default 1." +
					" 1: entry type, 2: standard type, 3: advanced type, 4: capacity type,  5: advanced type 1," +
					" 6: advanced type 2, 7: advanced type 3 , 8: advanced type 4, 9: exclusive type.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc id.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet id.",
			},
			"msg_retention_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The maximum retention time of instance logs, in minutes." +
					" the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled," +
					" which means that the log retention time recovery policy is not enabled.",
			},
			"cluster_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Cluster id.",
			},
			"renew_flag": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  false,
				Description: "Prepaid automatic renewal mark, 0 means the default state, the initial state," +
					" 1 means automatic renewal, 2 means clear no automatic renewal (user setting).",
			},
			"kafka_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "Kafka version (0.10.2/1.1.1/2.4.1).",
			},
			"spec_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "Default is profession.",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     false,
				Description: "Disk Size.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     false,
				Description: "Whether to open the ip whitelist, `true`: open, `false`: close.",
			},
			"partition": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     false,
				Description: "Partition size, the professional version does not need set.",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "Type of disk.",
			},
		},
	}
}

func resourceTencentCloudCkafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_instance.create")()
	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		ckafkaService = CkafkaService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		request = ckafka.NewCreateInstancePreRequest()
	)
	instanceName := d.Get("instance_name").(string)
	request.InstanceName = &instanceName

	zoneId := d.Get("zone_id").(int64)
	request.ZoneId = &zoneId

	period := d.Get("period").(string)
	request.Period = &period

	instanceType := d.Get("instance_type").(int64)
	request.InstanceType = &instanceType

	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId := v.(string)
		request.VpcId = helper.String(vpcId)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId := v.(string)
		request.SubnetId = helper.String(subnetId)
	}

	if v, ok := d.GetOk("msg_retention_time"); ok {
		retentionTime := v.(int64)
		request.MsgRetentionTime = helper.Int64(retentionTime)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId := v.(int64)
		request.ClusterId = helper.Int64(clusterId)
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		renewFlag := v.(int64)
		request.RenewFlag = helper.Int64(renewFlag)
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		kafkaVersion := v.(string)
		request.KafkaVersion = helper.String(kafkaVersion)
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		kafkaVersion := v.(string)
		request.KafkaVersion = helper.String(kafkaVersion)
	}

	if v, ok := d.GetOk("spec_type"); ok {
		specType := v.(string)
		request.SpecificationsType = helper.String(specType)
	}

	if v, ok := d.GetOk("disk_size"); ok {
		diskSize := v.(int64)
		request.DiskSize = helper.Int64(diskSize)
	}

	if v, ok := d.GetOk("band_width"); ok {
		bandWidth := v.(int64)
		request.BandWidth = helper.Int64(bandWidth)
	}

	if v, ok := d.GetOk("partition"); ok {
		partition := v.(int64)
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

	result, err := ckafkaService.client.



	return resourceTencentCloudCkafkaInstanceRead(d, meta)
}

func resourceTencentCloudCkafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.read")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		d.SetId("")
		return nil
	}
	instanceId := items[0]
	topicName := items[1]
	topicListInfo, hasExist, e := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if e != nil {
		return e
	}
	if !hasExist {
		d.SetId("")
		return nil
	}

	var topicinfo *ckafka.TopicAttributesResponse
	errInfo := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		topicDetail, e := ckafkcService.DescribeCkafkaTopicAttributes(ctx, instanceId, topicName)
		if e != nil {
			return retryError(e)
		}
		if topicDetail == nil {
			d.SetId("")
			return nil
		}
		topicinfo = topicDetail
		return nil
	})
	if errInfo != nil {
		return fmt.Errorf("[API]Describe kafka topic fail,reason:%s", errInfo.Error())
	}
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("note", topicinfo.Note)
	_ = d.Set("ip_white_list", topicinfo.IpWhiteList)
	_ = d.Set("ip_white_list_count", topicListInfo.IpWhiteListCount)
	_ = d.Set("enable_white_list", *topicinfo.EnableWhiteList == 1)
	_ = d.Set("replica_num", topicListInfo.ReplicaNum)
	_ = d.Set("create_time", topicinfo.CreateTime)
	_ = d.Set("partition_num", topicinfo.PartitionNum)
	_ = d.Set("topic_name", topicListInfo.TopicName)
	_ = d.Set("forward_interval", topicListInfo.ForwardInterval)
	_ = d.Set("forward_cos_bucket", topicListInfo.ForwardCosBucket)
	_ = d.Set("forward_status", topicListInfo.ForwardStatus)
	_ = d.Set("clean_up_policy", topicinfo.Config.CleanUpPolicy)
	_ = d.Set("max_message_bytes", topicinfo.Config.MaxMessageBytes)
	_ = d.Set("sync_replica_min_num", topicinfo.Config.MinInsyncReplicas)
	_ = d.Set("retention", topicinfo.Config.Retention)
	_ = d.Set("segment_bytes", topicinfo.Config.SegmentBytes)
	_ = d.Set("segment", topicinfo.Config.SegmentMs)
	if topicinfo.Config.UncleanLeaderElectionEnable != nil {
		_ = d.Set("unclean_leader_election_enable", *topicinfo.Config.UncleanLeaderElectionEnable == 1)
	}
	return nil
}

func resourceTencentCloudCkafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	request := ckafka.NewModifyTopicAttributesRequest()
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)
	whiteListSwitch := d.Get("enable_white_list").(bool)
	cleanUpPolicy := d.Get("clean_up_policy").(string)
	retention := d.Get("retention").(int)
	var note string
	if v, ok := d.GetOk("note"); ok {
		note = v.(string)
		request.Note = &note
	}
	if v, ok := d.GetOk("segment"); ok {
		if v.(int) != 0 {
			request.SegmentMs = helper.IntInt64(v.(int))
		}
	}
	request.InstanceId = &instanceId
	request.TopicName = &topicName
	request.EnableWhiteList = helper.BoolToInt64Ptr(whiteListSwitch)
	request.MinInsyncReplicas = helper.IntInt64(d.Get("sync_replica_min_num").(int))
	request.UncleanLeaderElectionEnable = helper.BoolToInt64Ptr(d.Get("unclean_leader_election_enable").(bool))
	request.CleanUpPolicy = &cleanUpPolicy
	request.RetentionMs = helper.IntInt64(retention)
	if d.Get("max_message_bytes").(int) != 0 {
		request.MaxMessageBytes = helper.IntInt64(d.Get("max_message_bytes").(int))
	}
	//Update ip white List
	if whiteListSwitch {
		oldInterface, newInterface := d.GetChange("ip_white_list")
		oldIpWhiteListInterface := oldInterface.([]interface{})
		newIpWhiteListInterface := newInterface.([]interface{})
		var oldIpWhiteList, newIpWhiteList []*string
		for _, value := range oldIpWhiteListInterface {
			oldIpWhiteList = append(oldIpWhiteList, helper.String(value.(string)))
		}
		for _, value := range newIpWhiteListInterface {
			newIpWhiteList = append(newIpWhiteList, helper.String(value.(string)))
		}
		if len(oldIpWhiteList) > 0 {
			error := ckafkcService.RemoveCkafkaTopicIpWhiteList(ctx, instanceId, topicName, oldIpWhiteList)
			if error != nil {
				return fmt.Errorf("IP whitelist Modification failed, reason[%s]\n", error.Error())
			}
		}
		if len(newIpWhiteList) == 0 {
			return fmt.Errorf("this Topic %s Create Failed, reason: ip whitelist switch is on, ip whitelist cannot be empty", topicName)
		}
		error := ckafkcService.AddCkafkaTopicIpWhiteList(ctx, instanceId, topicName, newIpWhiteList)
		if error != nil {
			return fmt.Errorf("IP whitelist Modification failed, reason[%s]\n", error.Error())
		}
	} else {
		//IP whiteList Switch not turned on, and the ip whitelist cannot be modified
		if d.HasChange("ip_white_list") {
			return fmt.Errorf("this Topic %s IP whitelist Modification failed, reason: The Ip Whitelist Switch is not turned on", topicName)
		}
	}
	//Update partition num
	oldPartitionNum, newPartitionNum := d.GetChange("partition_num")
	if newPartitionNum.(int) < oldPartitionNum.(int) {
		return fmt.Errorf("this Topic %s partition Modification failed, reason: The partitonNum must more than current partitionNum", topicName)
	} else {
		if newPartitionNum.(int) > oldPartitionNum.(int) {
			err := ckafkcService.AddCkafkaTopicPartition(ctx, instanceId, topicName, int64(newPartitionNum.(int)))
			if err != nil {
				return err
			}
		}
	}
	err := ckafkcService.ModifyCkafkaTopicAttribute(ctx, request)
	if err != nil {
		return err
	}
	_, has, errDes := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if errDes != nil {
		return errDes
	}
	if !has {
		errDes = fmt.Errorf("this Topic %s Update Failed", topicName)
		return errDes
	}

	return resourceTencentCloudCkafkaTopicRead(d, meta)
}

func resourceTencentCLoudCkafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)

	err := ckafkcService.DeleteCkafkaTopic(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	return nil
}
