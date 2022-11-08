/*
Use this data source to query detailed information of tdmq rabbitmqClusters

Example Usage

```hcl
data "tencentcloud_tdmq_rabbitmq_clusters" "rabbitmqClusters" {
  id_keyword = ""
  name_keyword = ""
  cluster_id_list = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqRabbitmqClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRabbitmqClustersRead,
		Schema: map[string]*schema.Schema{
			"id_keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster id keyword.",
			},

			"name_keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster name keyword.",
			},

			"cluster_id_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "cluster ids.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "cluster info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cluster id.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cluster name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "region.",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "create time.",
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "remark.",
									},
									"public_end_point": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "public end point.",
									},
									"vpc_end_point": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "vpc end point.",
									},
								},
							},
						},
						"config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "cluster config info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_tps_per_vhost": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max tps per vhost.",
									},
									"max_conn_num_per_vhost": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max connection number per vhost.",
									},
									"max_vhost_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max vhost number.",
									},
									"max_exchange_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max exchange number.",
									},
									"max_queue_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max queue number.",
									},
									"max_retention_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max retention.",
									},
									"used_vhost_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "used vhost number.",
									},
									"used_exchange_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "used exchange number.",
									},
									"used_queue_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "used queue number.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "tags info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag value.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status.",
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

func dataSourceTencentCloudTdmqRabbitmqClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_rabbitmq_clusters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("id_keyword"); ok {
		paramMap["id_keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name_keyword"); ok {
		paramMap["name_keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id_list"); ok {
		cluster_id_listSet := v.(*schema.Set).List()
		cluster_id_list := []*string{}
		for i := range cluster_id_listSet {
			cluster_id := cluster_id_listSet[i].(string)
			cluster_id_list = append(cluster_id_list, &cluster_id)
		}
		paramMap["cluster_id_list"] = cluster_id_list
	}

	tdmqService := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterList []*tdmq.AMQPClusterDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := tdmqService.DescribeTdmqRabbitmqClustersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterList = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Tdmq clusterList failed, reason:%+v", logId, err)
		return err
	}

	clusterListList := []interface{}{}
	if clusterList != nil {
		for _, clusterList := range clusterList {
			clusterListMap := map[string]interface{}{}
			if clusterList.Info != nil {
				infoMap := map[string]interface{}{}
				if clusterList.Info.ClusterId != nil {
					infoMap["cluster_id"] = clusterList.Info.ClusterId
				}
				if clusterList.Info.ClusterName != nil {
					infoMap["cluster_name"] = clusterList.Info.ClusterName
				}
				if clusterList.Info.Region != nil {
					infoMap["region"] = clusterList.Info.Region
				}
				if clusterList.Info.CreateTime != nil {
					infoMap["create_time"] = clusterList.Info.CreateTime
				}
				if clusterList.Info.Remark != nil {
					infoMap["remark"] = clusterList.Info.Remark
				}
				if clusterList.Info.PublicEndPoint != nil {
					infoMap["public_end_point"] = clusterList.Info.PublicEndPoint
				}
				if clusterList.Info.VpcEndPoint != nil {
					infoMap["vpc_end_point"] = clusterList.Info.VpcEndPoint
				}

				clusterListMap["info"] = []interface{}{infoMap}
			}
			if clusterList.Config != nil {
				configMap := map[string]interface{}{}
				if clusterList.Config.MaxTpsPerVHost != nil {
					configMap["max_tps_per_vhost"] = clusterList.Config.MaxTpsPerVHost
				}
				if clusterList.Config.MaxConnNumPerVHost != nil {
					configMap["max_conn_num_per_vhost"] = clusterList.Config.MaxConnNumPerVHost
				}
				if clusterList.Config.MaxVHostNum != nil {
					configMap["max_vhost_num"] = clusterList.Config.MaxVHostNum
				}
				if clusterList.Config.MaxExchangeNum != nil {
					configMap["max_exchange_num"] = clusterList.Config.MaxExchangeNum
				}
				if clusterList.Config.MaxQueueNum != nil {
					configMap["max_queue_num"] = clusterList.Config.MaxQueueNum
				}
				if clusterList.Config.MaxRetentionTime != nil {
					configMap["max_retention_time"] = clusterList.Config.MaxRetentionTime
				}
				if clusterList.Config.UsedVHostNum != nil {
					configMap["used_vhost_num"] = clusterList.Config.UsedVHostNum
				}
				if clusterList.Config.UsedExchangeNum != nil {
					configMap["used_exchange_num"] = clusterList.Config.UsedExchangeNum
				}
				if clusterList.Config.UsedQueueNum != nil {
					configMap["used_queue_num"] = clusterList.Config.UsedQueueNum
				}

				clusterListMap["config"] = []interface{}{configMap}
			}
			if clusterList.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range clusterList.Tags {
					tagsMap := map[string]interface{}{}
					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}
					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}
				clusterListMap["tags"] = tagsList
			}
			if clusterList.Status != nil {
				clusterListMap["status"] = clusterList.Status
			}

			clusterListList = append(clusterListList, clusterListMap)
		}
		err := d.Set("list", clusterListList)
		if err != nil {
			log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
			return err
		}
	}

	d.SetId("tencentcloud_tdmq_rabbitmq_clusters")

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), clusterListList); e != nil {
			return e
		}
	}

	return nil
}
