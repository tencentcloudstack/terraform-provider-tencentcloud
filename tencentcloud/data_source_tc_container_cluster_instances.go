/*
Get all instances of the specific cluster.

Use this data source to get all instances in a specific cluster.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.

Example Usage

```hcl
data "tencentcloud_container_cluster_instances" "foo_instance" {
  cluster_id = "cls-abcdefg"
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudContainerClusterInstances() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use 'tencentcloud_kubernetes_clusters' instead.",
		Read:               dataSourceTencentCloudContainerClusterInstancesRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An id identify the cluster, like cls-xxxxxx.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An int variable describe how many instances in return at most.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of instances.",
			},
			"nodes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "An information list of kubernetes instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"abnormal_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe the reason when node is in abnormal state(if it was).",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe the cpu of the node.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe the memory of the node.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An id identify the node, provided by cvm.",
						},
						"is_normal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe whether the node is normal.",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe the wan ip of the node.",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe the lan ip of the node.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_container_cluster_instances.read")()

	logId := getLogId(contextNil)

	request := tke.NewDescribeClusterInstancesRequest()
	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		request.ClusterId = common.StringPtr(clusterId.(string))
	}

	if limit, ok := d.GetOkExists("limit"); ok {
		request.Limit = common.Int64Ptr(limit.(int64))
	}

	var response *tke.DescribeClusterInstancesResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().DescribeClusterInstances(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s DescribeClusterInstances failed, reason:%s\n ", logId, err.Error())
		return err
	}

	nodes := make([]map[string]interface{}, 0, *response.Response.TotalCount)
	ids := make([]string, 0, *response.Response.TotalCount)
	for _, node := range response.Response.InstanceSet {
		ids = append(ids, *node.InstanceId)

		nodeInfo := make(map[string]interface{})
		nodeInfo["instance_id"] = *node.InstanceId
		nodeInfo["abnormal_reason"] = *node.FailedReason
		nodeInfo["wan_ip"] = ""
		nodeInfo["lan_ip"] = ""
		nodeInfo["cpu"] = 0
		nodeInfo["mem"] = 0
		if *node.InstanceState == "failed" {
			nodeInfo["is_normal"] = 0
		} else {
			nodeInfo["is_normal"] = 1
		}

		describeInstancesreq := cvm.NewDescribeInstancesRequest()
		describeInstancesreq.InstanceIds = []*string{node.InstanceId}
		var describeInstancesResponse *cvm.DescribeInstancesResponse
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().DescribeInstances(describeInstancesreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeInstancesreq.GetAction(), describeInstancesreq.ToJsonString(), e.Error())
				return retryError(e)
			}
			describeInstancesResponse = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s DescribeInstances failed, reason:%s\n ", logId, err.Error())
			return err
		}

		if len(describeInstancesResponse.Response.InstanceSet) > 0 {
			nodeInfo["cpu"] = *describeInstancesResponse.Response.InstanceSet[0].CPU
			nodeInfo["mem"] = *describeInstancesResponse.Response.InstanceSet[0].Memory
			if len(describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses) > 0 {
				nodeInfo["wan_ip"] = *describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses[0]
			}
			if len(describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses) > 0 {
				nodeInfo["lan_ip"] = *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0]
			}
		}

		nodes = append(nodes, nodeInfo)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("nodes", nodes)
	_ = d.Set("total_count", *response.Response.TotalCount)

	return nil
}
