package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by instance ID.",
			},
			"search_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance name, support fuzzy query.",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches the tag key value.",
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "(Filter Criteria) The status of the instance. 0: Create, 1: Run, 2: Delete, do not fill the default return all.",
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The field that needs to be filtered.",
						},
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "The filtered value of the field.",
						},
					},
				},
				Description: "Filter. filter.name supports ('Ip', 'VpcId', 'SubNetId', 'InstanceType','InstanceId'), filter.values can pass up to 10 values.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The page start offset, default is `0`.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "The number of pages, default is `10`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ckafka users. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance name.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP.",
						},
						"vport": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual PORT.",
						},
						"vip_list": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual IP.",
									},
									"vport": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual PORT.",
									},
								},
							},
							Description: "Virtual IP entities.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the instance. 0: Created, 1: Running, 2: Delete: 5 Quarantined, -1 Creation failed.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance bandwidth, in Mbps.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The storage size of the instance, in GB.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Availability Zone ID.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId, if empty, indicates that it is the underlying network.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id.",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the instance is renewed, the int enumeration value: 1 indicates auto-renewal, and 2 indicates that it is not automatically renewed.",
						},
						"healthy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status int: 1 indicates health, 2 indicates alarm, and 3 indicates abnormal instance status.",
						},
						"healthy_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status information.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time when the instance was created.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instance expiration time.",
						},
						"is_internal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is an internal customer. A value of 1 indicates an internal customer.",
						},
						"topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of topics.",
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Value.",
									},
								},
							},
							Description: "Tag information.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka version information. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"zone_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Across Availability Zones. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"cvm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ckafka sale type. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ckafka instance type. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk Type. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"max_topic_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of topics in the current specifications. Note: This field may return null, indicating that a valid value could not be retrieved..",
						},
						"max_partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of Partitions for the current specifications. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"rebalance_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schedule the upgrade configuration time. Note: This field may return null, indicating that a valid value could not be retrieved..",
						},
						"partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The current number of instances. Note: This field may return null, indicating that a valid value could not be retrieved..",
						},
						"public_network_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of Internet bandwidth. Note: This field may return null, indicating that a valid value could not be retrieved..",
						},
						"public_network": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Internet bandwidth value. Note: This field may return null, indicating that a valid value could not be retrieved..",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_instances.read")()

	ckafkaService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	request := ckafka.NewDescribeInstancesDetailRequest()
	if v, ok := d.GetOk("instance_ids"); ok {
		request.InstanceIdList = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("search_word"); ok {
		request.SearchWord = helper.String(v.(string))
	}
	if v, ok := d.GetOk("tag_key"); ok {
		request.TagKey = helper.String(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.InterfacesIntInt64Point(v.([]interface{}))
	}
	if v, ok := d.GetOk("filters"); ok {
		filterParams := v.([]interface{})
		filters := make([]*ckafka.Filter, 0)
		for _, filterParam := range filterParams {
			filterParamMap := filterParam.(map[string]interface{})
			filters = append(filters, &ckafka.Filter{
				Name:   helper.String(filterParamMap["name"].(string)),
				Values: helper.InterfacesStringsPoint(filterParamMap["values"].([]interface{})),
			})
		}
		request.Filters = filters
	}
	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.IntInt64(v.(int))
	}

	response, err := ckafkaService.client.UseCkafkaClient().DescribeInstancesDetail(request)
	if err != nil {
		return err
	}
	var kafkaInstanceDetails []*ckafka.InstanceDetail
	if response.Response.Result != nil {
		kafkaInstanceDetails = response.Response.Result.InstanceList
	}
	result := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, kafkaInstanceDetail := range kafkaInstanceDetails {
		kafkaInstanceDetailMap := make(map[string]interface{})
		ids = append(ids, *kafkaInstanceDetail.InstanceId)
		kafkaInstanceDetailMap["instance_id"] = kafkaInstanceDetail.InstanceId
		kafkaInstanceDetailMap["instance_name"] = kafkaInstanceDetail.InstanceName
		kafkaInstanceDetailMap["vip"] = kafkaInstanceDetail.Vip
		kafkaInstanceDetailMap["vport"] = kafkaInstanceDetail.Vport
		kafkaInstanceDetailMap["status"] = kafkaInstanceDetail.Status
		kafkaInstanceDetailMap["bandwidth"] = kafkaInstanceDetail.Bandwidth
		kafkaInstanceDetailMap["disk_size"] = kafkaInstanceDetail.DiskSize
		kafkaInstanceDetailMap["zone_id"] = kafkaInstanceDetail.ZoneId
		kafkaInstanceDetailMap["vpc_id"] = kafkaInstanceDetail.VpcId
		kafkaInstanceDetailMap["subnet_id"] = kafkaInstanceDetail.SubnetId
		kafkaInstanceDetailMap["renew_flag"] = kafkaInstanceDetail.RenewFlag
		kafkaInstanceDetailMap["healthy"] = kafkaInstanceDetail.Healthy
		kafkaInstanceDetailMap["healthy_message"] = kafkaInstanceDetail.HealthyMessage
		kafkaInstanceDetailMap["create_time"] = kafkaInstanceDetail.CreateTime
		kafkaInstanceDetailMap["expire_time"] = kafkaInstanceDetail.ExpireTime
		kafkaInstanceDetailMap["is_internal"] = kafkaInstanceDetail.IsInternal
		kafkaInstanceDetailMap["topic_num"] = kafkaInstanceDetail.TopicNum
		kafkaInstanceDetailMap["version"] = kafkaInstanceDetail.Version
		kafkaInstanceDetailMap["cvm"] = kafkaInstanceDetail.Cvm
		kafkaInstanceDetailMap["instance_type"] = kafkaInstanceDetail.InstanceType
		kafkaInstanceDetailMap["max_topic_number"] = kafkaInstanceDetail.MaxTopicNumber
		kafkaInstanceDetailMap["max_partition_number"] = kafkaInstanceDetail.MaxPartitionNumber
		kafkaInstanceDetailMap["rebalance_time"] = kafkaInstanceDetail.RebalanceTime
		kafkaInstanceDetailMap["partition_number"] = kafkaInstanceDetail.PartitionNumber
		kafkaInstanceDetailMap["public_network_charge_type"] = kafkaInstanceDetail.PublicNetworkChargeType
		kafkaInstanceDetailMap["public_network"] = kafkaInstanceDetail.PublicNetwork

		vipList := make([]map[string]string, 0)
		for _, vip := range kafkaInstanceDetail.VipList {
			vipList = append(vipList, map[string]string{
				"vip":   *vip.Vip,
				"vport": *vip.Vport,
			})
		}
		kafkaInstanceDetailMap["vip_list"] = vipList

		tags := make([]map[string]string, 0)
		for _, tag := range kafkaInstanceDetail.Tags {
			tags = append(tags, map[string]string{
				"tag_key":   *tag.TagKey,
				"tag_value": *tag.TagValue,
			})
		}
		kafkaInstanceDetailMap["tags"] = tags

		zoneIds := make([]int64, 0)
		for _, zoneId := range kafkaInstanceDetail.ZoneIds {
			zoneIds = append(zoneIds, *zoneId)
		}
		kafkaInstanceDetailMap["zone_ids"] = zoneIds

		result = append(result, kafkaInstanceDetailMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("instance_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
