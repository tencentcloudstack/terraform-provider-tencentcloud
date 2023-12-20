package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmChcHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmChcHostsRead,
		Schema: map[string]*schema.Schema{
			"chc_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CHC host ID. Up to 100 instances per request is allowed. ChcIds and Filters cannot be specified at the same time.",
			},

			"filters": {
				Optional: true,
				Type:     schema.TypeList,
				Description: "- `zone` Filter by the availability zone, such as ap-guangzhou-1. Valid values: See [Regions and Availability Zones](https://www.tencentcloud.com/document/product/213/6091?from_cn_redirect=1).\n" +
					"- `instance-name` Filter by the instance name.\n" +
					"- `instance-state` Filter by the instance status. For status details, see [InstanceStatus](https://www.tencentcloud.com/document/api/213/15753?from_cn_redirect=1#InstanceStatus).\n" +
					"- `device-type` Filter by the device type.\n" +
					"- `vpc-id` Filter by the unique VPC ID.\n" +
					"- `subnet-id` Filter by the unique VPC subnet ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values.",
						},
					},
				},
			},

			"chc_host_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of returned instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC host ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server serial number.",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC host status&lt;br/&gt;&lt;ul&gt;&lt;li&gt;REGISTERED: The CHC host is registered, but the out-of-band network and deployment network are not configured.&lt;/li&gt;&lt;li&gt;VPC_READY: The out-of-band network and deployment network are configured.&lt;/li&gt;&lt;li&gt;PREPARED: It&#39;s ready and can be associated with a CVM.&lt;/li&gt;&lt;li&gt;ONLINE: It&#39;s already associated with a CVM.&lt;/li&gt;&lt;/ul&gt;.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device typeNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"placement": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zone.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the availability zone where the instance resides. You can call the [DescribeZones](https://www.tencentcloud.com/document/product/213/35071) API and obtain the ID in the returned Zone field.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID of the project to which the instance belongs. This parameter can be obtained from the projectId returned by DescribeProject. If this is left empty, the default project is used.",
									},
									"host_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "ID list of CDHs from which the instance can be created. If you have purchased CDHs and specify this parameter, the instances you purchase will be randomly deployed on the CDHs.",
									},
									"host_ips": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "IPs of the hosts to create CVMs.",
									},
									"host_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the CDH to which the instance belongs, only used as an output parameter.",
									},
								},
							},
						},
						"bmc_virtual_private_cloud": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Out-of-band networkNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.",
									},
									"as_vpc_gateway": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.",
									},
									"private_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
									},
									"ipv6_address_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of IPv6 addresses randomly generated for the ENI.",
									},
								},
							},
						},
						"bmc_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Out-of-band network IPNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"bmc_security_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Out-of-band network security group IDNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"deploy_virtual_private_cloud": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Deployment networkNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.",
									},
									"as_vpc_gateway": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.",
									},
									"private_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
									},
									"ipv6_address_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of IPv6 addresses randomly generated for the ENI.",
									},
								},
							},
						},
						"deploy_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment network IPNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"deploy_security_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Deployment network security group IDNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"cvm_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the associated CVMNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server creation time.",
						},
						"hardware_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance hardware description, including CPU cores, memory capacity and disk capacity.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU cores of the CHC hostNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory capacity of the CHC host (unit: GB)Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"disk": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk capacity of the CHC hostNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"bmc_mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address assigned under the out-of-band networkNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"deploy_mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address assigned under the deployment networkNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"tenant_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Management typeHOSTING: HostingTENANT: LeasingNote: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudCvmChcHostsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_chc_hosts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("chc_ids"); ok {
		chcIdsSet := v.(*schema.Set).List()
		paramMap["ChcIds"] = helper.InterfacesStringsPoint(chcIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cvm.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := cvm.Filter{}
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

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var chcHostSet []*cvm.ChcHost

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmChcHostsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		chcHostSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(chcHostSet))
	tmpList := make([]map[string]interface{}, 0, len(chcHostSet))

	if chcHostSet != nil {
		for _, chcHost := range chcHostSet {
			chcHostMap := map[string]interface{}{}

			if chcHost.ChcId != nil {
				chcHostMap["chc_id"] = chcHost.ChcId
			}

			if chcHost.InstanceName != nil {
				chcHostMap["instance_name"] = chcHost.InstanceName
			}

			if chcHost.SerialNumber != nil {
				chcHostMap["serial_number"] = chcHost.SerialNumber
			}

			if chcHost.InstanceState != nil {
				chcHostMap["instance_state"] = chcHost.InstanceState
			}

			if chcHost.DeviceType != nil {
				chcHostMap["device_type"] = chcHost.DeviceType
			}

			if chcHost.Placement != nil {
				placementMap := map[string]interface{}{}

				if chcHost.Placement.Zone != nil {
					placementMap["zone"] = chcHost.Placement.Zone
				}

				if chcHost.Placement.ProjectId != nil {
					placementMap["project_id"] = chcHost.Placement.ProjectId
				}

				if chcHost.Placement.HostIds != nil {
					placementMap["host_ids"] = chcHost.Placement.HostIds
				}

				if chcHost.Placement.HostIps != nil {
					placementMap["host_ips"] = chcHost.Placement.HostIps
				}

				if chcHost.Placement.HostId != nil {
					placementMap["host_id"] = chcHost.Placement.HostId
				}

				chcHostMap["placement"] = []interface{}{placementMap}
			}

			if chcHost.BmcVirtualPrivateCloud != nil {
				bmcVirtualPrivateCloudMap := map[string]interface{}{}

				if chcHost.BmcVirtualPrivateCloud.VpcId != nil {
					bmcVirtualPrivateCloudMap["vpc_id"] = chcHost.BmcVirtualPrivateCloud.VpcId
				}

				if chcHost.BmcVirtualPrivateCloud.SubnetId != nil {
					bmcVirtualPrivateCloudMap["subnet_id"] = chcHost.BmcVirtualPrivateCloud.SubnetId
				}

				if chcHost.BmcVirtualPrivateCloud.AsVpcGateway != nil {
					bmcVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.BmcVirtualPrivateCloud.AsVpcGateway
				}

				if chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses != nil {
					bmcVirtualPrivateCloudMap["private_ip_addresses"] = chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses
				}

				if chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount != nil {
					bmcVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount
				}

				chcHostMap["bmc_virtual_private_cloud"] = []interface{}{bmcVirtualPrivateCloudMap}
			}

			if chcHost.BmcIp != nil {
				chcHostMap["bmc_ip"] = chcHost.BmcIp
			}

			if chcHost.BmcSecurityGroupIds != nil {
				chcHostMap["bmc_security_group_ids"] = chcHost.BmcSecurityGroupIds
			}

			if chcHost.DeployVirtualPrivateCloud != nil {
				deployVirtualPrivateCloudMap := map[string]interface{}{}

				if chcHost.DeployVirtualPrivateCloud.VpcId != nil {
					deployVirtualPrivateCloudMap["vpc_id"] = chcHost.DeployVirtualPrivateCloud.VpcId
				}

				if chcHost.DeployVirtualPrivateCloud.SubnetId != nil {
					deployVirtualPrivateCloudMap["subnet_id"] = chcHost.DeployVirtualPrivateCloud.SubnetId
				}

				if chcHost.DeployVirtualPrivateCloud.AsVpcGateway != nil {
					deployVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.DeployVirtualPrivateCloud.AsVpcGateway
				}

				if chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses != nil {
					deployVirtualPrivateCloudMap["private_ip_addresses"] = chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses
				}

				if chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount != nil {
					deployVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount
				}

				chcHostMap["deploy_virtual_private_cloud"] = []interface{}{deployVirtualPrivateCloudMap}
			}

			if chcHost.DeployIp != nil {
				chcHostMap["deploy_ip"] = chcHost.DeployIp
			}

			if chcHost.DeploySecurityGroupIds != nil {
				chcHostMap["deploy_security_group_ids"] = chcHost.DeploySecurityGroupIds
			}

			if chcHost.CvmInstanceId != nil {
				chcHostMap["cvm_instance_id"] = chcHost.CvmInstanceId
			}

			if chcHost.CreatedTime != nil {
				chcHostMap["created_time"] = chcHost.CreatedTime
			}

			if chcHost.HardwareDescription != nil {
				chcHostMap["hardware_description"] = chcHost.HardwareDescription
			}

			if chcHost.CPU != nil {
				chcHostMap["cpu"] = chcHost.CPU
			}

			if chcHost.Memory != nil {
				chcHostMap["memory"] = chcHost.Memory
			}

			if chcHost.Disk != nil {
				chcHostMap["disk"] = chcHost.Disk
			}

			if chcHost.BmcMAC != nil {
				chcHostMap["bmc_mac"] = chcHost.BmcMAC
			}

			if chcHost.DeployMAC != nil {
				chcHostMap["deploy_mac"] = chcHost.DeployMAC
			}

			if chcHost.TenantType != nil {
				chcHostMap["tenant_type"] = chcHost.TenantType
			}

			ids = append(ids, *chcHost.ChcId)
			tmpList = append(tmpList, chcHostMap)
		}

		_ = d.Set("chc_host_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
