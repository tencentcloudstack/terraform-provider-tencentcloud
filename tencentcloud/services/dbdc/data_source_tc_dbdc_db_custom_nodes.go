package dbdc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomNodesRead,
		Schema: map[string]*schema.Schema{
			"node_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query by one or more Node IDs. Maximum 100 IDs per request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported filter names: cluster-id, node-name (exact match), status (Creating, Running, Isolating, Isolated, Activating, Destroying), zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag Key and Value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"node_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"ssh_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSH endpoint for accessing the node, format: IP:Port.",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal network IP address of the node.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID the node belongs to.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone the node belongs to.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type/spec.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node CPU size in cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node memory size in GiB.",
						},
						"system_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System disk info. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk type.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Disk size in GiB.",
									},
								},
							},
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data disk list. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk type.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Disk size in GiB.",
									},
									"disk_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk name.",
									},
								},
							},
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node OS name.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node OS image ID.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID the node SSH endpoint belongs to.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID the node SSH endpoint belongs to.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node status.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payment type.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node expiration time.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node creation time.",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node isolation time.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node tag information. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto-renew flag. 1=auto renew, 0=no auto renew.",
						},
						"switch_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch ID.",
						},
						"rack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rack ID.",
						},
						"host_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host IP.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("node_ids"); ok {
		nodeIdsSet := v.([]interface{})
		tmpSet := make([]*string, 0, len(nodeIdsSet))
		for _, item := range nodeIdsSet {
			nodeId := item.(string)
			tmpSet = append(tmpSet, helper.String(nodeId))
		}
		paramMap["NodeIds"] = tmpSet
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dbdcv20201029.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := dbdcv20201029.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.([]interface{})
				for i := range valuesSet {
					value := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.([]interface{})
		tmpSet := make([]*dbdcv20201029.Tag, 0, len(tagsSet))
		for _, item := range tagsSet {
			tagsMap := item.(map[string]interface{})
			tag := dbdcv20201029.Tag{}
			if v, ok := tagsMap["key"].(string); ok && v != "" {
				tag.Key = helper.String(v)
			}

			if v, ok := tagsMap["value"].(string); ok {
				tag.Value = helper.String(v)
			}
			tmpSet = append(tmpSet, &tag)
		}
		paramMap["Tags"] = tmpSet
	}

	var respData []*dbdcv20201029.DBCustomNode
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeDBCustomNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return reqErr
	}

	nodeSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, node := range respData {
			nodeMap := map[string]interface{}{}
			if node.NodeId != nil {
				nodeMap["node_id"] = node.NodeId
			}

			if node.NodeName != nil {
				nodeMap["node_name"] = node.NodeName
			}

			if node.SSHEndpoint != nil {
				nodeMap["ssh_endpoint"] = node.SSHEndpoint
			}

			if node.LanIP != nil {
				nodeMap["lan_ip"] = node.LanIP
			}

			if node.ClusterId != nil {
				nodeMap["cluster_id"] = node.ClusterId
			}

			if node.Zone != nil {
				nodeMap["zone"] = node.Zone
			}

			if node.NodeType != nil {
				nodeMap["node_type"] = node.NodeType
			}

			if node.CPU != nil {
				nodeMap["cpu"] = node.CPU
			}

			if node.Memory != nil {
				nodeMap["memory"] = node.Memory
			}

			if node.SystemDisk != nil {
				systemDiskList := make([]map[string]interface{}, 0, 1)
				systemDiskMap := map[string]interface{}{}
				if node.SystemDisk.DiskType != nil {
					systemDiskMap["disk_type"] = node.SystemDisk.DiskType
				}
				if node.SystemDisk.DiskSize != nil {
					systemDiskMap["disk_size"] = node.SystemDisk.DiskSize
				}
				systemDiskList = append(systemDiskList, systemDiskMap)
				nodeMap["system_disk"] = systemDiskList
			}

			if node.DataDisks != nil {
				dataDisksList := make([]map[string]interface{}, 0, len(node.DataDisks))
				for _, dataDisk := range node.DataDisks {
					dataDiskMap := map[string]interface{}{}
					if dataDisk.DiskType != nil {
						dataDiskMap["disk_type"] = dataDisk.DiskType
					}
					if dataDisk.DiskSize != nil {
						dataDiskMap["disk_size"] = dataDisk.DiskSize
					}
					if dataDisk.DiskName != nil {
						dataDiskMap["disk_name"] = dataDisk.DiskName
					}
					dataDisksList = append(dataDisksList, dataDiskMap)
				}
				nodeMap["data_disks"] = dataDisksList
			}

			if node.OsName != nil {
				nodeMap["os_name"] = node.OsName
			}

			if node.ImageId != nil {
				nodeMap["image_id"] = node.ImageId
			}

			if node.VpcId != nil {
				nodeMap["vpc_id"] = node.VpcId
			}

			if node.SubnetId != nil {
				nodeMap["subnet_id"] = node.SubnetId
			}

			if node.Status != nil {
				nodeMap["status"] = node.Status
			}

			if node.ChargeType != nil {
				nodeMap["charge_type"] = node.ChargeType
			}

			if node.ExpireTime != nil {
				nodeMap["expire_time"] = node.ExpireTime
			}

			if node.CreatedTime != nil {
				nodeMap["created_time"] = node.CreatedTime
			}

			if node.IsolatedTime != nil {
				nodeMap["isolated_time"] = node.IsolatedTime
			}

			if node.Tags != nil {
				tagsList := make([]map[string]interface{}, 0, len(node.Tags))
				for _, tag := range node.Tags {
					tagMap := map[string]interface{}{}
					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}
					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}
					tagsList = append(tagsList, tagMap)
				}
				nodeMap["tags"] = tagsList
			}

			if node.AutoRenew != nil {
				nodeMap["auto_renew"] = node.AutoRenew
			}

			if node.SwitchId != nil {
				nodeMap["switch_id"] = node.SwitchId
			}

			if node.RackId != nil {
				nodeMap["rack_id"] = node.RackId
			}

			if node.HostIp != nil {
				nodeMap["host_ip"] = node.HostIp
			}

			nodeSetList = append(nodeSetList, nodeMap)
		}

		_ = d.Set("node_set", nodeSetList)
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
