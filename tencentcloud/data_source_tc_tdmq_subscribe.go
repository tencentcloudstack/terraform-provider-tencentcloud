/*
Use this data source to query detailed information of tdmq subscribe

Example Usage

```hcl
data "tencentcloud_tdmq_subscribe" "subscribe" {
  topic_name = "topic_name"
  offset = 0
  limit = 20
  subscription_name = &lt;nil&gt;
  total_count = &lt;nil&gt;
  subscription_set {
		subscription_name = &lt;nil&gt;
		subscription_id = &lt;nil&gt;
		topic_owner = &lt;nil&gt;
		msg_count = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		create_time = &lt;nil&gt;
		binding_key = &lt;nil&gt;
		endpoint = &lt;nil&gt;
		filter_tags = &lt;nil&gt;
		protocol = &lt;nil&gt;
		notify_strategy = &lt;nil&gt;
		notify_content_format = &lt;nil&gt;

  }
  request_id = &lt;nil&gt;
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

func dataSourceTencentCloudTdmqSubscribe() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqSubscribeRead,
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of the list of topics to be returned on the current page in case of paginated return. If a value is entered, limit is required. If this parameter is left empty, 0 will be used by default.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of topics to be returned per page in case of paginated return. If this parameter is not passed in, 20 will be used by default. Maximum value: 50.",
			},

			"subscription_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by SubscriptionName.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number.",
			},

			"subscription_set": {
				Type:        schema.TypeList,
				Description: "Set of subscription attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_name": {
							Type:        schema.TypeString,
							Description: "Subscription name, which must be unique in the same topic under the same account in the same region. It can contain up to 64 letters, digits, and hyphens and must begin with a letter.",
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Description: "Subscription ID, which will be used during monitoring data pull.",
						},
						"topic_owner": {
							Type:        schema.TypeInt,
							Description: "Subscription owner APPID.",
						},
						"msg_count": {
							Type:        schema.TypeInt,
							Description: "Number of messages to be delivered in the subscription.",
						},
						"last_modify_time": {
							Type:        schema.TypeInt,
							Description: "Time when the subscription attribute is last modified. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Description: "Subscription creation time. A Unix timestamp accurate down to the millisecond will be returned.",
						},
						"binding_key": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering policy for subscribing to and receiving messages.",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Description: "Endpoint that receives notifications, which varies by `protocol`: for HTTP, the endpoint must start with `http://`, and the `host` can be a domain or IP; for `queue`, `queueName` should be entered.",
						},
						"filter_tags": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering policy selected when a subscription is created:If `filterType` is 1, `filterTag` will be used for filtering. If `filterType` is 2, `bindingKey` will be used for filtering.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "Subscription protocol. Currently, two protocols are supported: HTTP and queue. To use the HTTP protocol, you need to build your own web server to receive messages. With the queue protocol, messages are automatically pushed to a CMQ queue and you can pull them concurrently.",
						},
						"notify_strategy": {
							Type:        schema.TypeString,
							Description: "CMQ push server retry policy in case an error occurs while pushing a message to `Endpoint`. Valid values: 1. `BACKOFF_RETRY`: backoff retry, which is to retry at a fixed interval, discard the message after a certain number of retries, and continue to push the next message; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry, which is to retry at an exponentially increasing interval, such as 1s, 2s, 4s, 8s, and so on. As a message can be retained in a topic for one day, failed messages will be discarded at most after one day of retry. Default value: `EXPONENTIAL_DECAY_RETRY`.",
						},
						"notify_content_format": {
							Type:        schema.TypeString,
							Description: "Push content format. Valid values: 1. `JSON`; 2. `SIMPLIFIED`, i.e., the raw format. If `Protocol` is `queue`, this value must be `SIMPLIFIED`. If `Protocol` is `http`, both options are acceptable, and the default value is `JSON`.",
						},
					},
				},
			},

			"request_id": {
				Type:        schema.TypeString,
				Description: "The unique request ID, which is returned for each request. RequestId is required for locating a problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_subscribe.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["TopicName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		paramMap["SubscriptionName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("subscription_set"); ok {
		subscriptionSetSet := v.([]interface{})
		tmpSet := make([]*tdmq.CmqSubscription, 0, len(subscriptionSetSet))

		for _, item := range subscriptionSetSet {
			cmqSubscription := tdmq.CmqSubscription{}
			cmqSubscriptionMap := item.(map[string]interface{})

			if v, ok := cmqSubscriptionMap["subscription_name"]; ok {
				cmqSubscription.SubscriptionName = helper.String(v.(string))
			}
			if v, ok := cmqSubscriptionMap["subscription_id"]; ok {
				cmqSubscription.SubscriptionId = helper.String(v.(string))
			}
			if v, ok := cmqSubscriptionMap["topic_owner"]; ok {
				cmqSubscription.TopicOwner = helper.IntUint64(v.(int))
			}
			if v, ok := cmqSubscriptionMap["msg_count"]; ok {
				cmqSubscription.MsgCount = helper.IntUint64(v.(int))
			}
			if v, ok := cmqSubscriptionMap["last_modify_time"]; ok {
				cmqSubscription.LastModifyTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqSubscriptionMap["create_time"]; ok {
				cmqSubscription.CreateTime = helper.IntUint64(v.(int))
			}
			if v, ok := cmqSubscriptionMap["binding_key"]; ok {
				bindingKeySet := v.(*schema.Set).List()
				cmqSubscription.BindingKey = helper.InterfacesStringsPoint(bindingKeySet)
			}
			if v, ok := cmqSubscriptionMap["endpoint"]; ok {
				cmqSubscription.Endpoint = helper.String(v.(string))
			}
			if v, ok := cmqSubscriptionMap["filter_tags"]; ok {
				filterTagsSet := v.(*schema.Set).List()
				cmqSubscription.FilterTags = helper.InterfacesStringsPoint(filterTagsSet)
			}
			if v, ok := cmqSubscriptionMap["protocol"]; ok {
				cmqSubscription.Protocol = helper.String(v.(string))
			}
			if v, ok := cmqSubscriptionMap["notify_strategy"]; ok {
				cmqSubscription.NotifyStrategy = helper.String(v.(string))
			}
			if v, ok := cmqSubscriptionMap["notify_content_format"]; ok {
				cmqSubscription.NotifyContentFormat = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &cmqSubscription)
		}
		paramMap["subscription_set"] = tmpSet
	}

	if v, ok := d.GetOk("request_id"); ok {
		paramMap["RequestId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var subscriptionSet []*tdmq.CmqSubscription

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqSubscribeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		subscriptionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(subscriptionSet))
	tmpList := make([]map[string]interface{}, 0, len(subscriptionSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
