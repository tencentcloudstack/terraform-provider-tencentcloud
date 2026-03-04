package ckafka

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaInstancesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaInstancesV2Read,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance ID.",
			},

			"search_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance name, instance ID, availability zone, VPC ID or subnet ID. Fuzzy search is supported.",
			},

			"status": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter by instance status. -1: creation failed, 0: creating, 1: running, 2: deleting, 3: deleted, 4: deletion failed, 5: isolated, 7: upgrading.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by tag key.",
			},

			"instance_id_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter by instance ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag list (intersection).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported fields: Ip, VpcId, SubNetId, InstanceType, InstanceId. Note: filter.Values can contain up to 10 values.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name. Supported: Ip, VpcId, SubNetId, InstanceType, InstanceId.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values (up to 10 values).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance VIP.",
						},
						"vport": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance port.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status. 0: creating, 1: running, 2: deleting, 5: isolated, -1: creation failed.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance bandwidth in Mbps.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance disk size in GB.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone ID.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"healthy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance health status. 1: healthy, 2: alarm, 3: abnormal.",
						},
						"healthy_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance health information.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance creation time (Unix timestamp).",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance expiration time (Unix timestamp).",
						},
						"is_internal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is an internal customer. 1: internal, 0: external.",
						},
						"topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current number of topics.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka version number.",
						},
						"zone_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cross-availability zone.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"cvm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CKafka sale type.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type.",
						},
						"max_topic_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of topics.",
						},
						"max_partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of partitions.",
						},
						"rebalance_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Planned upgrade configuration time.",
						},
						"partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current partition number.",
						},
						"public_network_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network bandwidth billing mode.",
						},
						"public_network": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public network bandwidth.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster type.",
						},
						"features": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dynamic message retention policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto-renewal flag. 0: default state (user has not set, the initial state is auto-renewal), 1: auto-renewal, 2: explicit no auto-renewal (user has set).",
						},
						"vip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual IP list.",
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
										Description: "Virtual port.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCkafkaInstancesV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_instances_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	// Handle instance_id parameter
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = v.(string)
	}

	// Handle search_word parameter
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = v.(string)
	}

	// Handle status parameter
	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		statusList := make([]interface{}, 0, len(statusSet))
		for _, status := range statusSet {
			statusList = append(statusList, status.(int))
		}
		paramMap["Status"] = statusList
	}

	// Handle tag_key parameter
	if v, ok := d.GetOk("tag_key"); ok {
		paramMap["TagKey"] = v.(string)
	}

	// Handle instance_id_list parameter
	if v, ok := d.GetOk("instance_id_list"); ok {
		instanceIdSet := v.(*schema.Set).List()
		instanceIdList := make([]interface{}, 0, len(instanceIdSet))
		for _, id := range instanceIdSet {
			instanceIdList = append(instanceIdList, id.(string))
		}
		paramMap["InstanceIdList"] = instanceIdList
	}

	// Handle tag_list parameter
	if v, ok := d.GetOk("tag_list"); ok {
		tagList := v.([]interface{})
		tagListParam := make([]interface{}, 0, len(tagList))
		for _, tag := range tagList {
			tagMap := tag.(map[string]interface{})
			tagListParam = append(tagListParam, map[string]interface{}{
				"tag_key":   tagMap["tag_key"].(string),
				"tag_value": tagMap["tag_value"].(string),
			})
		}
		paramMap["TagList"] = tagListParam
	}

	// Handle filters parameter
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		filters := make([]interface{}, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := make(map[string]interface{})

			if name, ok := filtersMap["name"].(string); ok && name != "" {
				filter["name"] = name
			}

			if values, ok := filtersMap["values"]; ok {
				valueSet := values.(*schema.Set).List()
				valueList := make([]interface{}, 0, len(valueSet))
				for _, value := range valueSet {
					valueList = append(valueList, value.(string))
				}
				filter["values"] = valueList
			}

			if len(filter) > 0 {
				filters = append(filters, filter)
			}
		}
		if len(filters) > 0 {
			paramMap["Filters"] = filters
		}
	}

	var respData []*ckafka.InstanceDetail
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeInstancesDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	instanceList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instance := range respData {
			instanceMap := map[string]interface{}{}

			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}

			if instance.InstanceName != nil {
				instanceMap["instance_name"] = instance.InstanceName
			}

			if instance.Vip != nil {
				instanceMap["vip"] = instance.Vip
			}

			if instance.Vport != nil {
				instanceMap["vport"] = instance.Vport
			}

			if instance.Status != nil {
				instanceMap["status"] = instance.Status
			}

			if instance.Bandwidth != nil {
				instanceMap["bandwidth"] = instance.Bandwidth
			}

			if instance.DiskSize != nil {
				instanceMap["disk_size"] = instance.DiskSize
			}

			if instance.ZoneId != nil {
				instanceMap["zone_id"] = instance.ZoneId
			}

			if instance.VpcId != nil {
				instanceMap["vpc_id"] = instance.VpcId
			}

			if instance.SubnetId != nil {
				instanceMap["subnet_id"] = instance.SubnetId
			}

			if instance.Healthy != nil {
				instanceMap["healthy"] = instance.Healthy
			}

			if instance.HealthyMessage != nil {
				instanceMap["healthy_message"] = instance.HealthyMessage
			}

			if instance.CreateTime != nil {
				instanceMap["create_time"] = instance.CreateTime
			}

			if instance.ExpireTime != nil {
				instanceMap["expire_time"] = instance.ExpireTime
			}

			if instance.IsInternal != nil {
				instanceMap["is_internal"] = instance.IsInternal
			}

			if instance.TopicNum != nil {
				instanceMap["topic_num"] = instance.TopicNum
			}

			if instance.Tags != nil {
				tagsList := make([]map[string]interface{}, 0, len(instance.Tags))
				for _, tag := range instance.Tags {
					tagMap := map[string]interface{}{}
					if tag.TagKey != nil {
						tagMap["tag_key"] = tag.TagKey
					}
					if tag.TagValue != nil {
						tagMap["tag_value"] = tag.TagValue
					}
					tagsList = append(tagsList, tagMap)
				}
				instanceMap["tags"] = tagsList
			}

			if instance.Version != nil {
				instanceMap["version"] = instance.Version
			}

			if instance.ZoneIds != nil {
				instanceMap["zone_ids"] = instance.ZoneIds
			}

			if instance.Cvm != nil {
				instanceMap["cvm"] = instance.Cvm
			}

			if instance.InstanceType != nil {
				instanceMap["instance_type"] = instance.InstanceType
			}

			if instance.DiskType != nil {
				instanceMap["disk_type"] = instance.DiskType
			}

			if instance.MaxTopicNumber != nil {
				instanceMap["max_topic_number"] = instance.MaxTopicNumber
			}

			if instance.MaxPartitionNumber != nil {
				instanceMap["max_partition_number"] = instance.MaxPartitionNumber
			}

			if instance.RebalanceTime != nil {
				instanceMap["rebalance_time"] = instance.RebalanceTime
			}

			if instance.PartitionNumber != nil {
				instanceMap["partition_number"] = instance.PartitionNumber
			}

			if instance.PublicNetworkChargeType != nil {
				instanceMap["public_network_charge_type"] = instance.PublicNetworkChargeType
			}

			if instance.PublicNetwork != nil {
				instanceMap["public_network"] = instance.PublicNetwork
			}

			if instance.ClusterType != nil {
				instanceMap["cluster_type"] = instance.ClusterType
			}

			if instance.Features != nil {
				instanceMap["features"] = instance.Features
			}

			if instance.RenewFlag != nil {
				instanceMap["renew_flag"] = instance.RenewFlag
			}

			if instance.VipList != nil {
				vipList := make([]map[string]interface{}, 0, len(instance.VipList))
				for _, vipItem := range instance.VipList {
					vipMap := map[string]interface{}{}
					if vipItem.Vip != nil {
						vipMap["vip"] = vipItem.Vip
					}
					if vipItem.Vport != nil {
						vipMap["vport"] = vipItem.Vport
					}
					vipList = append(vipList, vipMap)
				}
				instanceMap["vip_list"] = vipList
			}

			instanceList = append(instanceList, instanceMap)
		}

		_ = d.Set("instance_list", instanceList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
