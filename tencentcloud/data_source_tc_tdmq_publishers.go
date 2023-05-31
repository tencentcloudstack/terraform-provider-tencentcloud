/*
Use this data source to query detailed information of tdmq publishers

Example Usage

```hcl
data "tencentcloud_tdmq_publishers" "publishers" {
  cluster_id = "pulsar-9n95ax58b9vn"
  namespace  = "keep-ns"
  topic      = "keep-topic"
  filters {
    name   = "ProducerName"
    values = ["test"]
  }
  sort {
    name  = "ProducerName"
    order = "DESC"
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqPublishers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqPublishersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},
			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "topic name.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameter filter, support ProducerName, Address field.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the filter parameter.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "value.",
						},
					},
				},
			},
			"sort": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "sorter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "sorter.",
						},
						"order": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ascending ASC, descending DESC.",
						},
					},
				},
			},
			// computed
			"publishers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Producer Information ListNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"producer_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "producer idNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"producer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "producer nameNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "producer addressNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"client_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "client versionNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"msg_rate_in": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Message production rate (articles/second)Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"msg_throughput_in": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Message production throughput rate (bytes/second)Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"average_msg_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average message size (bytes)Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"connected_since": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "connection timeNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The topic partition number of the producer connectionNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTdmqPublishersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_publishers.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		publishers []*tdmq.Publisher
		clusterId  string
		Namespace  string
		Topic      string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
		Namespace = v.(string)
	}

	if v, ok := d.GetOk("topic"); ok {
		paramMap["Topic"] = helper.String(v.(string))
		Topic = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := tdmq.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "sort"); ok {
		sort := tdmq.Sort{}
		if v, ok := dMap["name"]; ok {
			sort.Name = helper.String(v.(string))
		}
		if v, ok := dMap["order"]; ok {
			sort.Order = helper.String(v.(string))
		}
		paramMap["sort"] = &sort
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqPublishersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		publishers = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0, len(publishers))

	if publishers != nil {
		for _, publisher := range publishers {
			publisherMap := map[string]interface{}{}

			if publisher.ProducerId != nil {
				publisherMap["producer_id"] = publisher.ProducerId
			}

			if publisher.ProducerName != nil {
				publisherMap["producer_name"] = publisher.ProducerName
			}

			if publisher.Address != nil {
				publisherMap["address"] = publisher.Address
			}

			if publisher.ClientVersion != nil {
				publisherMap["client_version"] = publisher.ClientVersion
			}

			if publisher.MsgRateIn != nil {
				publisherMap["msg_rate_in"] = publisher.MsgRateIn
			}

			if publisher.MsgThroughputIn != nil {
				publisherMap["msg_throughput_in"] = publisher.MsgThroughputIn
			}

			if publisher.AverageMsgSize != nil {
				publisherMap["average_msg_size"] = publisher.AverageMsgSize
			}

			if publisher.ConnectedSince != nil {
				publisherMap["connected_since"] = publisher.ConnectedSince
			}

			if publisher.Partition != nil {
				publisherMap["partition"] = publisher.Partition
			}

			tmpList = append(tmpList, publisherMap)
		}

		_ = d.Set("publishers", tmpList)
	}

	ids = append(ids, clusterId)
	ids = append(ids, Namespace)
	ids = append(ids, Topic)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
