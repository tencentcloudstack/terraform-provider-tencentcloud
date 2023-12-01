/*
Use this data source to query detailed information of tdmqRocketmq namespace

Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  name_keyword = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
  remark         = "remark."
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqRocketmqNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRocketmqNamespaceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"name_keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by name.",
			},

			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of namespaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace name, which can contain 3-64 letters, digits, hyphens, and underscores.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention time of unconsumed messages in milliseconds. Value range: 60 seconds-15 days.",
						},
						"retention_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention time of persisted messages in milliseconds.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks (up to 128 characters).",
						},
						"public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network access point address.",
						},
						"vpc_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC access point address.",
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

func dataSourceTencentCloudTdmqRocketmqNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmqRocketmq_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["cluster_id"] = v.(string)
	}

	if v, ok := d.GetOk("name_keyword"); ok {
		paramMap["name_keyword"] = v.(string)
	}

	tdmqRocketmqService := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	ids := make([]string, 0)
	var namespaces []*tdmqRocketmq.RocketMQNamespace
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := tdmqRocketmqService.DescribeTdmqRocketmqNamespaceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		namespaces = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read TdmqRocketmq namespaces failed, reason:%+v", logId, err)
		return err
	}

	namespaceList := []interface{}{}

	for _, namespace := range namespaces {
		namespaceMap := map[string]interface{}{}
		ids = append(ids, *namespace.NamespaceId)
		namespaceMap["namespace_id"] = namespace.NamespaceId
		if namespace.Ttl != nil {
			namespaceMap["ttl"] = namespace.Ttl
		}
		if namespace.RetentionTime != nil {
			namespaceMap["retention_time"] = namespace.RetentionTime
		}
		if namespace.Remark != nil {
			namespaceMap["remark"] = namespace.Remark
		}
		if namespace.PublicEndpoint != nil {
			namespaceMap["public_endpoint"] = namespace.PublicEndpoint
		}
		if namespace.VpcEndpoint != nil {
			namespaceMap["vpc_endpoint"] = namespace.VpcEndpoint
		}

		namespaceList = append(namespaceList, namespaceMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("namespaces", namespaceList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), namespaceList); e != nil {
			return e
		}
	}

	return nil
}
