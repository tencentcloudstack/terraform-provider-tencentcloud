package cdwpg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdwpgNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdwpgNodesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},

			"instance_nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Node id.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type.",
						},
						"node_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
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

func dataSourceTencentCloudCdwpgNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdwpg_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	var respData *cdwpgv20201230.DescribeInstanceNodesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdwpgNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	instanceNodesList := make([]map[string]interface{}, 0, len(respData.InstanceNodes))
	if respData.InstanceNodes != nil {
		for _, instanceNodes := range respData.InstanceNodes {
			instanceNodesMap := map[string]interface{}{}

			if instanceNodes.NodeId != nil {
				instanceNodesMap["node_id"] = instanceNodes.NodeId
			}

			if instanceNodes.NodeType != nil {
				instanceNodesMap["node_type"] = instanceNodes.NodeType
			}

			if instanceNodes.NodeIp != nil {
				instanceNodesMap["node_ip"] = instanceNodes.NodeIp
			}

			instanceNodesList = append(instanceNodesList, instanceNodesMap)
		}

		_ = d.Set("instance_nodes", instanceNodesList)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
