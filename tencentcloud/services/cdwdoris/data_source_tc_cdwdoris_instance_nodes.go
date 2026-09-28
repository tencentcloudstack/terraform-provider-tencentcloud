package cdwdoris

import (
	"context"
	"log"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwdorisv20211228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdwdorisInstanceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdwdorisInstanceNodesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster instance ID, `cdw-xxxx` string type.",
			},
			"node_role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster role type, default is `data` data node.",
			},
			"display_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display policy, `All` to display all.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// computed
			"instance_nodes_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cluster node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Model, such as S1. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"core": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of CPU cores. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the clickhouse cluster to which it belongs. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"rip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rip. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"fe_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "FE node role. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"virtual_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual availability zone. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"compute_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compute group ID. Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},
			"node_roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cluster-supported node role types.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceTencentCloudCdwdorisInstanceNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdwdoris_instance_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = CdwdorisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		respData  []*cdwdorisv20211228.InstanceNode
		respRoles []*string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_role"); ok {
		paramMap["NodeRole"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display_policy"); ok {
		paramMap["DisplayPolicy"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, roles, e := service.DescribeCdwdorisInstanceNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(result) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe cdwdoris_instance_nodes failed, response is empty."))
		}

		respData = result
		respRoles = roles
		return nil
	})

	if err != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return err
	}

	tmpInstanceNodesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instanceNode := range respData {
			instanceNodeMap := map[string]interface{}{}
			if instanceNode.Ip != nil {
				instanceNodeMap["ip"] = instanceNode.Ip
			}

			if instanceNode.Spec != nil {
				instanceNodeMap["spec"] = instanceNode.Spec
			}

			if instanceNode.Core != nil {
				instanceNodeMap["core"] = instanceNode.Core
			}

			if instanceNode.Memory != nil {
				instanceNodeMap["memory"] = instanceNode.Memory
			}

			if instanceNode.DiskType != nil {
				instanceNodeMap["disk_type"] = instanceNode.DiskType
			}

			if instanceNode.DiskSize != nil {
				instanceNodeMap["disk_size"] = instanceNode.DiskSize
			}

			if instanceNode.Role != nil {
				instanceNodeMap["role"] = instanceNode.Role
			}

			if instanceNode.Status != nil {
				instanceNodeMap["status"] = instanceNode.Status
			}

			if instanceNode.Rip != nil {
				instanceNodeMap["rip"] = instanceNode.Rip
			}

			if instanceNode.FeRole != nil {
				instanceNodeMap["fe_role"] = instanceNode.FeRole
			}

			if instanceNode.UUID != nil {
				instanceNodeMap["uuid"] = instanceNode.UUID
			}

			if instanceNode.Zone != nil {
				instanceNodeMap["zone"] = instanceNode.Zone
			}

			if instanceNode.VirtualZone != nil {
				instanceNodeMap["virtual_zone"] = instanceNode.VirtualZone
			}

			if instanceNode.CreateTime != nil {
				instanceNodeMap["create_time"] = instanceNode.CreateTime
			}

			if instanceNode.ComputeGroupId != nil {
				instanceNodeMap["compute_group_id"] = instanceNode.ComputeGroupId
			}

			tmpInstanceNodesList = append(tmpInstanceNodesList, instanceNodeMap)
		}

		_ = d.Set("instance_nodes_list", tmpInstanceNodesList)
	}

	if respRoles != nil {
		tmpNodeRoles := make([]string, 0, len(respRoles))
		for _, role := range respRoles {
			if role == nil {
				continue
			}
			tmpNodeRoles = append(tmpNodeRoles, *role)
		}
		_ = d.Set("node_roles", tmpNodeRoles)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpInstanceNodesList); e != nil {
			return e
		}
	}

	return nil
}
