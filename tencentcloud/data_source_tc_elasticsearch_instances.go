/*
Use this data source to query elasticsearch instances.

Example Usage

```hcl
data "tencentcloud_elasticsearch_instances" "foo" {
  instance_id = "es-17634f05"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instance to be queried.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the instance to be queried.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag of the instance to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of elasticsearch instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the instance.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of a VPC network.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of a VPC subnetwork.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the instance.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of instance.",
						},
						"deploy_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster deployment mode.",
						},
						"multi_zone_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details of AZs in multi-AZ deployment mode.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The id of a VPC subnetwork.",
									},
								},
							},
						},
						"license_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "License type.",
						},
						"node_info_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node information list, which describe the specification information of various types of nodes in the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of nodes.",
									},
									"node_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node specification.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node type.",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node disk type.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Node disk size.",
									},
								},
							},
						},
						"basic_security_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable X-Pack security authentication in Basic Edition 6.8 and above.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A mapping of tags to assign to the instance.",
						},
						"elasticsearch_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elasticsearch domain name.",
						},
						"elasticsearch_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elasticsearch VIP.",
						},
						"elasticsearch_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elasticsearch port.",
						},
						"kibana_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kibana access URL.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance creation time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudElasticsearchInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	instanceId := ""
	instanceName := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		instanceName = v.(string)
	}
	tags := helper.GetTags(d, "tags")
	var instances []*es.InstanceInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = elasticsearchService.DescribeInstancesByFilter(ctx, instanceId, instanceName, tags)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return nil
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		tags := make(map[string]string, len(instance.TagList))
		for _, tag := range instance.TagList {
			tags[*tag.TagKey] = *tag.TagValue
		}

		mapping := map[string]interface{}{
			"instance_id":          instance.InstanceId,
			"instance_name":        instance.InstanceName,
			"availability_zone":    instance.Zone,
			"vpc_id":               instance.VpcUid,
			"subnet_id":            instance.SubnetUid,
			"version":              instance.EsVersion,
			"charge_type":          instance.ChargeType,
			"deploy_mode":          instance.DeployMode,
			"license_type":         instance.LicenseType,
			"basic_security_type":  instance.SecurityType,
			"tags":                 tags,
			"elasticsearch_domain": instance.EsDomain,
			"elasticsearch_vip":    instance.EsVip,
			"elasticsearch_port":   instance.EsPort,
			"kibana_url":           instance.KibanaUrl,
			"create_time":          instance.CreateTime,
		}
		if instance.MultiZoneInfo != nil && len(instance.MultiZoneInfo) > 0 {
			infos := make([]map[string]interface{}, 0, len(instance.MultiZoneInfo))
			for _, v := range instance.MultiZoneInfo {
				info := map[string]interface{}{
					"availability_zone": v.Zone,
					"subnet_id":         v.SubnetId,
				}
				infos = append(infos, info)
			}
			mapping["multi_zone_info"] = infos
		}
		if instance.NodeInfoList != nil && len(instance.NodeInfoList) > 0 {
			infos := make([]map[string]interface{}, 0, len(instance.NodeInfoList))
			for _, v := range instance.NodeInfoList {
				info := map[string]interface{}{
					"node_num":  v.NodeNum,
					"node_type": v.NodeType,
					"type":      v.Type,
					"disk_type": v.DiskType,
					"disk_size": v.DiskSize,
				}
				infos = append(infos, info)
			}
			mapping["node_info_list"] = infos
		}

		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set elasticsearch instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
