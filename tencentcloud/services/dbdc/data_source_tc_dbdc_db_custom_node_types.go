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

func DataSourceTencentCloudDbdcDbCustomNodeTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomNodeTypesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported filter names: region, zone, node-family, node-type.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"node_type_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom node type list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone, such as ap-guangzhou-6.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type, such as DB.SA5.2XLARGE32.",
						},
						"node_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node family, such as DB.AT5, DB.SA5.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in GiB.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type sell status. Values: SELL, SOLD_OUT.",
						},
						"system_disk_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System disk types supported by this node type. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"data_disk_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data disk types supported by this node type. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomNodeTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_node_types.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
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

	var respData []*dbdcv20201029.DBCustomNodeTypeInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeDBCustomNodeTypesByFilter(ctx, paramMap)
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

	nodeTypeSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, nodeType := range respData {
			nodeTypeMap := map[string]interface{}{}
			if nodeType.Zone != nil {
				nodeTypeMap["zone"] = nodeType.Zone
			}

			if nodeType.NodeType != nil {
				nodeTypeMap["node_type"] = nodeType.NodeType
			}

			if nodeType.NodeFamily != nil {
				nodeTypeMap["node_family"] = nodeType.NodeFamily
			}

			if nodeType.CPU != nil {
				nodeTypeMap["cpu"] = nodeType.CPU
			}

			if nodeType.Memory != nil {
				nodeTypeMap["memory"] = nodeType.Memory
			}

			if nodeType.Status != nil {
				nodeTypeMap["status"] = nodeType.Status
			}

			if nodeType.SystemDiskTypes != nil {
				systemDiskTypesList := make([]string, 0, len(nodeType.SystemDiskTypes))
				for _, systemDiskType := range nodeType.SystemDiskTypes {
					if systemDiskType != nil {
						systemDiskTypesList = append(systemDiskTypesList, *systemDiskType)
					}
				}
				nodeTypeMap["system_disk_types"] = systemDiskTypesList
			}

			if nodeType.DataDiskTypes != nil {
				dataDiskTypesList := make([]string, 0, len(nodeType.DataDiskTypes))
				for _, dataDiskType := range nodeType.DataDiskTypes {
					if dataDiskType != nil {
						dataDiskTypesList = append(dataDiskTypesList, *dataDiskType)
					}
				}
				nodeTypeMap["data_disk_types"] = dataDiskTypesList
			}

			nodeTypeSetList = append(nodeTypeSetList, nodeTypeMap)
		}

		_ = d.Set("node_type_set", nodeTypeSetList)
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
