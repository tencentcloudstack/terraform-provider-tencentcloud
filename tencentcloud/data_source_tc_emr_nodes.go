/*
Provides an available EMR for the user.

The EMR data source obtain the hardware node information by using the emr cluster ID.

Example Usage

```hcl
data "tencentcloud_emr_nodes" "my_emr_nodes" {
  node_flag="master"
  instance_id="emr-rnzqrleq"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEmrNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrNodesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster instance ID, the instance ID is as follows: emr-xxxxxxxx.",
			},
			"node_flag": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Node ID, the value is:
				- all: Means to get all type nodes, except cdb information.
				- master: Indicates that the master node information is obtained.
				- core: Indicates that the core node information is obtained.
				- task: indicates obtaining task node information.
				- common: means to get common node information.
				- router: Indicates obtaining router node information.
				- db: Indicates that the cdb information for the normal state is obtained.
				- recyle: Indicates that the node information in the Recycle Bin isolation, including the cdb information, is obtained.
				- renew: Indicates that all node information to be renewed, including cddb information, is obtained, and the auto-renewal node will not be returned.
				
				Note: Only the above values are now supported, entering other values will cause an error.`,
			},
			"hardware_resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				Description: "Resource type: Support all/host/pod, default is all.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Page number, with a default value of 0, represents the first page.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The number returned per page, the default value is 100, and the maximum value is 100.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of node details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User APPID.",
						},
						"serial_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Serial number.",
						},
						"order_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine instance ID.",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The master node is bound to the Internet IP address.",
						},
						"flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node type. 0: common node; 1: master node; 2: core node; 3: task node.",
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node specifications.",
						},
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of node cores.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node memory.",
						},
						"mem_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node memory description.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The node is located in the region.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone where the node is located.",
						},
						"apply_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application time.",
						},
						"free_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Release time.",
						},
						"disk_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hard disk size.",
						},
						"name_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node description.",
						},
						"services": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node deployment service.",
						},
						"storage_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk type.",
						},
						"root_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the system disk.",
						},
						"charge_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The type of payment.",
						},
						"cdb_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database IP.",
						},
						"cdb_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database port.",
						},
						"hw_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hard disk capacity.",
						},
						"hw_disk_size_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hard disk capacity description.",
						},
						"hw_mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory capacity.",
						},
						"hw_mem_size_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Memory capacity description.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time.",
						},
						"emr_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node resource ID.",
						},
						"is_auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Renewal logo.",
						},
						"device_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device identity.",
						},
						"mutable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Supports variations.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet IP.",
						},
						"destroyable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether this node is destroyable, 1 can be destroyed, 0 is not destroyable.",
						},
						"auto_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is an autoscaling node, 0 is a normal node, and 1 is an autoscaling node.",
						},
						"hardware_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type, host/pod.",
						},
						"is_dynamic_spec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Floating specifications, 1 yes, 0 no.",
						},
						"dynamic_pod_spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating specification value json string.",
						},
						"support_modify_pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to support change billing type 1 Yes and 0 No.",
						},
						"cdb_node_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Database information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DB instance.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database IP.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database port.",
									},
									"mem_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database memory specifications.",
									},
									"volume": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database disk specifications.",
									},
									"service": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service identity.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expiration time.",
									},
									"apply_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application time.",
									},
									"pay_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The type of payment.",
									},
									"expire_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Expired id.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database status.",
									},
									"is_auto_renew": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Renewal identity.",
									},
									"serial_no": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database string.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Zone Id.",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Region id.",
									},
								},
							},
						},
						"mc_multi_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Multi-cloud disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of cloud disks of this type.",
									},
									"type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Disk type.",
									},
									"volume": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the cloud disk.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The label of the node binding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEmrNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_emr_nodes.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	nodeFlag := d.Get("node_flag").(string)
	offset := d.Get("offset").(int)
	limit := d.Get("limit").(int)
	hardwareResourceType := d.Get("hardware_resource_type").(string)

	emrServer := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var nodes []*emr.NodeHardwareInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var errRet error
		nodes, errRet = emrServer.DescribeClusterNodes(ctx, instanceId, nodeFlag, hardwareResourceType, offset, limit)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	emrNodes := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, node := range nodes {
		mcMultiDisks := node.MCMultiDisk
		cdbNodeInfo := node.CdbNodeInfo
		tags := node.Tags
		mcMultiDiskList := make([]map[string]interface{}, 0)
		cdbNodeInfoMap := make(map[string]interface{})
		tagList := make([]map[string]interface{}, 0)
		for _, mcMultiDisk := range mcMultiDisks {
			tmpMCMultiDisk := make(map[string]interface{})
			tmpMCMultiDisk["count"] = mcMultiDisk.Count
			tmpMCMultiDisk["type"] = mcMultiDisk.Type
			tmpMCMultiDisk["volume"] = mcMultiDisk.Volume
			mcMultiDiskList = append(mcMultiDiskList, tmpMCMultiDisk)
		}
		for _, tag := range tags {
			tmpTag := make(map[string]interface{})
			tmpTag["tag_key"] = tag.TagKey
			tmpTag["tag_value"] = tag.TagValue
			tagList = append(tagList, tmpTag)
		}

		if cdbNodeInfo != nil {
			cdbNodeInfoMap["instance_name"] = cdbNodeInfo.InstanceName
			cdbNodeInfoMap["ip"] = cdbNodeInfo.Ip
			cdbNodeInfoMap["port"] = cdbNodeInfo.Port
			cdbNodeInfoMap["mem_size"] = cdbNodeInfo.MemSize
			cdbNodeInfoMap["volume"] = cdbNodeInfo.Volume
			cdbNodeInfoMap["service"] = cdbNodeInfo.Service
			cdbNodeInfoMap["expire_time"] = cdbNodeInfo.ExpireTime
			cdbNodeInfoMap["apply_time"] = cdbNodeInfo.ApplyTime
			cdbNodeInfoMap["pay_type"] = cdbNodeInfo.PayType
			cdbNodeInfoMap["expire_flag"] = cdbNodeInfo.ExpireFlag
			cdbNodeInfoMap["status"] = cdbNodeInfo.Status
			cdbNodeInfoMap["is_auto_renew"] = cdbNodeInfo.IsAutoRenew
			cdbNodeInfoMap["serial_no"] = cdbNodeInfo.SerialNo
			cdbNodeInfoMap["zone_id"] = cdbNodeInfo.ZoneId
			cdbNodeInfoMap["region_id"] = cdbNodeInfo.RegionId
		}

		nodeMap := map[string]interface{}{
			"app_id":                  node.AppId,
			"serial_no":               node.SerialNo,
			"order_no":                node.OrderNo,
			"wan_ip":                  node.WanIp,
			"flag":                    node.Flag,
			"spec":                    node.Spec,
			"cpu_num":                 node.CpuNum,
			"mem_size":                node.MemSize,
			"mem_desc":                node.MemDesc,
			"region_id":               node.RegionId,
			"zone_id":                 node.ZoneId,
			"apply_time":              node.ApplyTime,
			"free_time":               node.FreeTime,
			"disk_size":               node.DiskSize,
			"name_tag":                node.NameTag,
			"services":                node.Services,
			"storage_type":            node.StorageType,
			"root_size":               node.RootSize,
			"charge_type":             node.ChargeType,
			"cdb_ip":                  node.CdbIp,
			"cdb_port":                node.CdbPort,
			"hw_disk_size":            node.HwDiskSize,
			"hw_disk_size_desc":       node.HwDiskSizeDesc,
			"hw_mem_size":             node.HwMemSize,
			"hw_mem_size_desc":        node.HwMemSizeDesc,
			"expire_time":             node.ExpireTime,
			"emr_resource_id":         node.EmrResourceId,
			"is_auto_renew":           node.IsAutoRenew,
			"device_class":            node.DeviceClass,
			"mutable":                 node.Mutable,
			"ip":                      node.Ip,
			"destroyable":             node.Destroyable,
			"auto_flag":               node.AutoFlag,
			"hardware_resource_type":  node.HardwareResourceType,
			"is_dynamic_spec":         node.IsDynamicSpec,
			"dynamic_pod_spec":        node.DynamicPodSpec,
			"support_modify_pay_mode": node.SupportModifyPayMode,
			"cdb_node_info":           []map[string]interface{}{cdbNodeInfoMap},
			"mc_multi_disks":          mcMultiDiskList,
			"tags":                    tagList,
		}
		ids = append(ids, (string)(*node.AppId))
		emrNodes = append(emrNodes, nodeMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("nodes", emrNodes)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), emrNodes); err != nil {
			return err
		}
	}
	return nil
}
