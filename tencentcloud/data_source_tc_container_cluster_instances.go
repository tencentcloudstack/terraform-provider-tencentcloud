package tencentcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	ccs "github.com/zqfan/tencentcloud-sdk-go/services/ccs/unversioned"
)

func dataSourceTencentCloudContainerClusterInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudContainerClusterInstancesRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed values
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nodes": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"abnormal_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_normal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"wan_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lan_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).ccsConn
	describeClusterInstancesReq := ccs.NewDescribeClusterInstancesRequest()

	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		describeClusterInstancesReq.ClusterId = common.StringPtr(clusterId.(string))
	}

	if limit, ok := d.GetOkExists("limit"); ok {
		describeClusterInstancesReq.Limit = common.IntPtr(limit.(int))
	}

	response, err := client.DescribeClusterInstances(describeClusterInstancesReq)
	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("data_source_tencent_cloud_container_cluster_instances got error, no code response")
	}

	if *response.Code != 0 {
		return fmt.Errorf("data_source_tencent_cloud_container_cluster_instances got error, code %v , message %v", *response.Code, *response.CodeDesc)
	}

	id := fmt.Sprintf("%d", time.Now().Unix())
	nodes := make([]map[string]interface{}, 0)
	for _, node := range response.Data.Nodes {
		nodeInfo := make(map[string]interface{}, 0)
		if node.AbnormalReason != nil {
			nodeInfo["abnormal_reason"] = *node.AbnormalReason
		}
		if node.CPU != nil {
			nodeInfo["cpu"] = *node.CPU
		}
		if node.Mem != nil {
			nodeInfo["mem"] = *node.Mem
		}
		if node.InstanceId != nil {
			nodeInfo["instance_id"] = *node.InstanceId
		}
		if node.IsNormal != nil {
			nodeInfo["is_normal"] = *node.IsNormal
		}
		if node.WanIp != nil {
			nodeInfo["wan_ip"] = *node.WanIp
		}
		if node.LanIp != nil {
			nodeInfo["lan_ip"] = *node.LanIp
		}
		nodes = append(nodes, nodeInfo)
	}

	d.Set("nodes", nodes)
	d.SetId(id)
	if response.Data.TotalCount != nil {
		d.Set("total_count", *response.Data.TotalCount)
	} else {
		d.Set("total_count", 0)
	}

	return nil
}
