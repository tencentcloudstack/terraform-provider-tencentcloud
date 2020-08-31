/*
Use this resource to create ckafka topic instance.

Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "foo" {
	instance_id						= "ckafka-f9ife4zz"
	topic_name						= "example"
	note							= "topic note"
	replica_num						= 2
	partition_num					= 1
	enable_white_list				= 1
	ip_white_list    				= ["ip1","ip2"]
	clean_up_policy					= "delete"
	sync_replica_min_num			= 1
	unclean_leader_election_enable  = false
	segment							= 3600000
	retention						= 60000
	max_message_bytes				= 0
}
```

Import

ckafka topic instance can be imported using the instance_id#topic_name, e.g.

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
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCkafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaTopicCreate,
		Read:   resourceTencentCloudCkafkaTopicRead,
		Update: resourceTencentCloudCkafkaTopicUpdate,
		Delete: resourceTencentCLoudCkafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ckafka instance ID.",
			},
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-). The length range is from 1 to 64.",
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 24),
				Description:  "The number of partition.",
			},
			"replica_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 3),
				Description:  "The number of replica, the maximum is 3.",
			},
			"enable_white_list": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "IP Whitelist switch, 1: open; 0: close.",
			},
			"ip_white_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ip whitelist, quota limit, required when enableWhileList=1.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"note": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 64),
				Description:  "The subject note is a string of no more than 64 characters. It must start with a letter, and the remaining part can contain letters, numbers and dashes (-).",
			},
			"retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60000,
				ValidateFunc: validateIntegerMin(60000),
				Description:  "Message can be selected. Retention time, unit ms, the current minimum value is 60000ms.",
			},
			"sync_replica_min_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Min number of sync replicas,Default is 1.",
			},
			"clean_up_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "delete",
				Description: "Clear log policy, log clear mode, the default is delete. delete: logs are deleted according to the storage time, compact: logs are compressed according to the key, compact, delete: logs are compressed according to the key and will be deleted according to the storage time.",
			},
			"unclean_leader_election_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to allow unsynchronized replicas to be selected as leader, false: not allowed, true: allowed, not allowed by default.",
			},
			"max_message_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 8388608),
				Description:  "Max message bytes.",
			},
			"segment": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerMin(3600000),
				Description:  "Segment scrolling time, in ms, the current minimum is 3600000ms.",
			},
			"message_storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message storage location.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the topic instance.",
			},
			"forward_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Periodic frequency of data backup to cos.",
			},
			"forward_cos_bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data backup cos bucket: the bucket address that is dumped to cos.",
			},
			"forward_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Data backup cos status: 1 do not open data backup, 0 open data backup.",
			},
			"segment_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of bytes rolled by shard.",
			},
		},
	}
}

func resourceTencentCloudCkafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	request, error := resourceTencentCLoudGetCreateTopicRequest(d)
	if error != nil {
		return error
	}
	instanceId := *request.InstanceId
	topicName := *request.TopicName
	ipWhiteList := request.IpWhiteList
	//Before create topic,Check if kafka exists
	_, has, error := ckafkcService.DescribeInstanceById(ctx, instanceId)
	if error != nil {
		return error
	}
	if !has {
		return fmt.Errorf("ckafka %s does not exist", instanceId)
	}
	errCreate := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, err := ckafkcService.CreateCkafkaTopic(ctx, request)
		if err != nil {
			return retryError(err)
		}
		return resource.NonRetryableError(err)
	})
	if errCreate != nil {
		return errCreate
	}
	_, hasExist, err := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if err != nil {
		return err
	}
	if !hasExist {
		err = fmt.Errorf("this Topic %s Create Failed, timeout", topicName)
	}
	if err != nil {
		return err
	}

	if len(ipWhiteList) > 0 {
		err = ckafkcService.AddCkafkaTopicIpWhiteList(ctx, instanceId, topicName, ipWhiteList)
		if err != nil {
			return err
		}
	}
	resourceId := instanceId + FILED_SP + topicName
	d.SetId(resourceId)
	return resourceTencentCloudCkafkaTopicRead(d, meta)
}

func resourceTencentCloudCkafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.read")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	items := strings.Split(d.Id(), FILED_SP)
	instanceId := items[0]
	topicName := items[1]
	topicListInfo, hasExist, e := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if e != nil {
		return fmt.Errorf("[API]Describe kafka topic fail,reason:%s", e.Error())
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
	_ = d.Set("enable_white_list", topicinfo.EnableWhiteList)
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
	_ = d.Set("unclean_leader_election_enable", topicinfo.Config.UncleanLeaderElectionEnable)
	return nil
}

func resourceTencentCloudCkafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	request, error := resourceTencentCLoudGetModifyTopicRequest(d)
	if error != nil {
		return error
	}
	instanceId := *request.InstanceId
	topicName := *request.TopicName

	unsupportedUpdateFields := []string{
		"topic_name",
		"partition_num",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("Template resource_ckafka_topic update on %s is not supportted yet. Please renew it on controller web page.", field)
		}
	}
	//update ip white List
	if *request.EnableWhiteList != 0 {
		if d.HasChange("ip_white_list") {
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
			error := ckafkcService.RemoveCkafkaTopicIpWhiteList(ctx, instanceId, topicName, oldIpWhiteList)
			if error != nil {
				return fmt.Errorf("IP whitelist Modification failed")
			}
			error = ckafkcService.AddCkafkaTopicIpWhiteList(ctx, instanceId, topicName, newIpWhiteList)
			if error != nil {
				return fmt.Errorf("IP whitelist Modification failed")
			}
		}
	}

	//Before update topic,Check if kafka exists
	has, error := ckafkcService.DescribeCkafkaById(ctx, instanceId)
	if error != nil {
		return error
	}
	if !has {
		return fmt.Errorf("ckafka %s does not exist", instanceId)
	}
	errRet, err := ckafkcService.ModifyCkafkaTopicAttribute(ctx, request)
	if err != nil {
		return err
	}
	if errRet != nil {
		return errRet
	}

	_, has, errDes := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if errDes != nil {
		return errDes
	}
	if !has {
		errDes = fmt.Errorf("this Topic %s Update Failed, timeout", topicName)
		return errDes
	}

	resourceId := instanceId + FILED_SP + topicName
	d.SetId(resourceId)
	return resourceTencentCloudCkafkaTopicRead(d, meta)
}

func resourceTencentCLoudCkafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ckafka_topic.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)

	//Before delete topic,Check if kafka exists
	has, error := ckafkcService.DescribeCkafkaById(ctx, instanceId)
	if error != nil {
		return error
	}
	if !has {
		return fmt.Errorf("ckafka %s does not exist", instanceId)
	}
	//Check if topic exists
	_, hasExist, e := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if e != nil {
		return e
	}
	if !hasExist {
		return fmt.Errorf("ckafka %s, topic %s does not exist", instanceId, topicName)
	}

	err := ckafkcService.DeleteCkafkaTopic(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCLoudGetCreateTopicRequest(d *schema.ResourceData) (request *ckafka.CreateTopicRequest, err error) {
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)
	partitionNum := *helper.IntInt64(d.Get("partition_num").(int))
	replicaNum := *helper.IntInt64(d.Get("replica_num").(int))
	syncReplicaMinNum := *helper.IntInt64(d.Get("sync_replica_min_num").(int))
	uncleanLeaderElectionEnable := d.Get("unclean_leader_election_enable").(bool)
	whiteListSwitchFlag := d.Get("enable_white_list").(int)
	//var ipWhiteList []*string
	var cleanUpPolicy, note string
	var retentionMs, segmentMs int64
	var whiteListSwitch bool
	ipWhiteLists := d.Get("ip_white_list").([]interface{})
	var ipWhiteList = make([]*string, 0, len(ipWhiteLists))
	for _, value := range ipWhiteLists {
		ipWhiteList = append(ipWhiteList, helper.String(value.(string)))
	}
	if v, ok := d.GetOk("clean_up_policy"); ok {
		cleanUpPolicy = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		note = v.(string)
	}
	if v, ok := d.GetOk("retention"); ok {
		retentionMs = *helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("segment"); ok {
		segmentMs = *helper.IntInt64(v.(int))
	}
	if whiteListSwitchFlag != 0 {
		whiteListSwitch = true
	}
	request = ckafka.NewCreateTopicRequest()
	request.InstanceId = &instanceId
	request.TopicName = &topicName
	request.PartitionNum = &partitionNum
	request.ReplicaNum = &replicaNum
	request.EnableWhiteList = helper.BoolToInt64Ptr(whiteListSwitch)
	if whiteListSwitch {
		request.IpWhiteList = ipWhiteList
	}
	request.CleanUpPolicy = &cleanUpPolicy
	if note != "" {
		request.Note = &note
	}
	request.MinInsyncReplicas = &syncReplicaMinNum
	request.UncleanLeaderElectionEnable = helper.BoolToInt64Ptr(uncleanLeaderElectionEnable)
	if retentionMs != 0 {
		request.RetentionMs = &retentionMs
	}
	if segmentMs != 0 {
		request.SegmentMs = &segmentMs
	}
	return request, err
}
func resourceTencentCLoudGetModifyTopicRequest(d *schema.ResourceData) (request *ckafka.ModifyTopicAttributesRequest, err error) {
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)
	syncReplicaMinNum := *helper.IntInt64(d.Get("sync_replica_min_num").(int))
	uncleanLeaderElectionEnable := d.Get("unclean_leader_election_enable").(bool)
	whiteListSwitchFlag := d.Get("enable_white_list").(int)
	maxMessageBytes := *helper.IntInt64(d.Get("max_message_bytes").(int))
	var cleanUpPolicy, note string
	var retentionMs, segmentMs int64
	var whiteListSwitch bool
	if v, ok := d.GetOk("clean_up_policy"); ok {
		cleanUpPolicy = v.(string)
	}
	if v, ok := d.GetOk("note"); ok {
		note = v.(string)
	}
	if v, ok := d.GetOk("retention"); ok {
		retentionMs = *helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("segment"); ok {
		segmentMs = *helper.IntInt64(v.(int))
	}
	if whiteListSwitchFlag != 0 {
		whiteListSwitch = true
	}
	request = ckafka.NewModifyTopicAttributesRequest()
	request.InstanceId = &instanceId
	request.TopicName = &topicName
	request.EnableWhiteList = helper.BoolToInt64Ptr(whiteListSwitch)
	request.CleanUpPolicy = &cleanUpPolicy
	if note != "" {
		request.Note = &note
	}
	request.MinInsyncReplicas = &syncReplicaMinNum
	request.UncleanLeaderElectionEnable = helper.BoolToInt64Ptr(uncleanLeaderElectionEnable)
	if retentionMs != 0 {
		request.RetentionMs = &retentionMs
	}
	if segmentMs != 0 {
		request.SegmentMs = &segmentMs
	}
	if maxMessageBytes != 0 {
		request.MaxMessageBytes = &maxMessageBytes
	}
	return request, err
}
