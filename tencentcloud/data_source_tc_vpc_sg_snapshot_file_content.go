package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcSgSnapshotFileContent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSgSnapshotFileContentRead,
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy IDs.",
			},

			"snapshot_file_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot file ID.",
			},

			"security_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Security group ID.",
			},

			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Security group ID.",
			},

			"backup_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup time.",
			},

			"operator": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Operator.",
			},

			"original_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Original data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index number of security group rules, which dynamically changes with the rules. This parameter can be obtained via the `DescribeSecurityGroupPolicies` API and used with the `Version` field in the returned value of the API.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol. Valid values: TCP, UDP, ICMP, ICMPv6, ALL.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port (`all`, a single port, or a port range).Note: If the `Protocol` value is set to `ALL`, the `Port` value also needs to be set to `all`.",
						},
						"service_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Protocol port ID or protocol port group ID. ServiceTemplate and Protocol+Port are mutually exclusive.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port ID, such as `ppm-f5n1f8da`.",
									},
									"service_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port group ID, such as `ppmg-f5n1f8da`.",
									},
								},
							},
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either `CidrBlock` or `Ipv6CidrBlock can be specified. Note that if `0.0.0.0/n` is entered, it is mapped to 0.0.0.0/0.",
						},
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block or IPv6 (mutually exclusive).",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group instance ID, such as `sg-ohuuioma`.",
						},
						"address_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address ID or IP address group ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the IP address, such as `ipm-2uw6ujo6`.",
									},
									"address_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the IP address group, such as `ipmg-2uw6ujo6`.",
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACCEPT or DROP.",
						},
						"policy_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group policy description.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last modification time of the security group.",
						},
					},
				},
			},

			"backup_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Backup data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The index number of security group rules, which dynamically changes with the rules. This parameter can be obtained via the `DescribeSecurityGroupPolicies` API and used with the `Version` field in the returned value of the API.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol. Valid values: TCP, UDP, ICMP, ICMPv6, ALL.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port (`all`, a single port, or a port range).Note: If the `Protocol` value is set to `ALL`, the `Port` value also needs to be set to `all`.",
						},
						"service_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Protocol port ID or protocol port group ID. ServiceTemplate and Protocol+Port are mutually exclusive.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port ID, such as `ppm-f5n1f8da`.",
									},
									"service_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port group ID, such as `ppmg-f5n1f8da`.",
									},
								},
							},
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either `CidrBlock` or `Ipv6CidrBlock can be specified. Note that if `0.0.0.0/n` is entered, it is mapped to 0.0.0.0/0.",
						},
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block or IPv6 (mutually exclusive).",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group instance ID, such as `sg-ohuuioma`.",
						},
						"address_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address ID or IP address group ID.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the IP address, such as `ipm-2uw6ujo6`.",
									},
									"address_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the IP address group, such as `ipmg-2uw6ujo6`.",
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACCEPT or DROP.",
						},
						"policy_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group policy description.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last modification time of the security group.",
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

func dataSourceTencentCloudVpcSgSnapshotFileContentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_sg_snapshot_file_content.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		paramMap["SnapshotPolicyId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_file_id"); ok {
		paramMap["SnapshotFileId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		paramMap["SecurityGroupId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var content *vpc.DescribeSgSnapshotFileContentResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSgSnapshotFileContent(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		content = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0)

	if content.InstanceId != nil {
		_ = d.Set("instance_id", content.InstanceId)
	}

	if content.BackupTime != nil {
		_ = d.Set("backup_time", content.BackupTime)
	}

	if content.Operator != nil {
		_ = d.Set("operator", content.Operator)
	}

	if content.OriginalData != nil {
		for _, securityGroupPolicy := range content.OriginalData {
			securityGroupPolicyMap := map[string]interface{}{}

			if securityGroupPolicy.PolicyIndex != nil {
				securityGroupPolicyMap["policy_index"] = securityGroupPolicy.PolicyIndex
			}

			if securityGroupPolicy.Protocol != nil {
				securityGroupPolicyMap["protocol"] = securityGroupPolicy.Protocol
			}

			if securityGroupPolicy.Port != nil {
				securityGroupPolicyMap["port"] = securityGroupPolicy.Port
			}

			if securityGroupPolicy.ServiceTemplate != nil {
				serviceTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.ServiceTemplate.ServiceId != nil {
					serviceTemplateMap["service_id"] = securityGroupPolicy.ServiceTemplate.ServiceId
				}

				if securityGroupPolicy.ServiceTemplate.ServiceGroupId != nil {
					serviceTemplateMap["service_group_id"] = securityGroupPolicy.ServiceTemplate.ServiceGroupId
				}

				securityGroupPolicyMap["service_template"] = []interface{}{serviceTemplateMap}
			}

			if securityGroupPolicy.CidrBlock != nil {
				securityGroupPolicyMap["cidr_block"] = securityGroupPolicy.CidrBlock
			}

			if securityGroupPolicy.Ipv6CidrBlock != nil {
				securityGroupPolicyMap["ipv6_cidr_block"] = securityGroupPolicy.Ipv6CidrBlock
			}

			if securityGroupPolicy.SecurityGroupId != nil {
				securityGroupPolicyMap["security_group_id"] = securityGroupPolicy.SecurityGroupId
			}

			if securityGroupPolicy.AddressTemplate != nil {
				addressTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.AddressTemplate.AddressId != nil {
					addressTemplateMap["address_id"] = securityGroupPolicy.AddressTemplate.AddressId
				}

				if securityGroupPolicy.AddressTemplate.AddressGroupId != nil {
					addressTemplateMap["address_group_id"] = securityGroupPolicy.AddressTemplate.AddressGroupId
				}

				securityGroupPolicyMap["address_template"] = []interface{}{addressTemplateMap}
			}

			if securityGroupPolicy.Action != nil {
				securityGroupPolicyMap["action"] = securityGroupPolicy.Action
			}

			if securityGroupPolicy.PolicyDescription != nil {
				securityGroupPolicyMap["policy_description"] = securityGroupPolicy.PolicyDescription
			}

			if securityGroupPolicy.ModifyTime != nil {
				securityGroupPolicyMap["modify_time"] = securityGroupPolicy.ModifyTime
			}

			ids = append(ids, *securityGroupPolicy.SecurityGroupId)
			tmpList = append(tmpList, securityGroupPolicyMap)
		}

		_ = d.Set("original_data", tmpList)
	}

	if content.BackupData != nil {
		for _, securityGroupPolicy := range content.BackupData {
			securityGroupPolicyMap := map[string]interface{}{}

			if securityGroupPolicy.PolicyIndex != nil {
				securityGroupPolicyMap["policy_index"] = securityGroupPolicy.PolicyIndex
			}

			if securityGroupPolicy.Protocol != nil {
				securityGroupPolicyMap["protocol"] = securityGroupPolicy.Protocol
			}

			if securityGroupPolicy.Port != nil {
				securityGroupPolicyMap["port"] = securityGroupPolicy.Port
			}

			if securityGroupPolicy.ServiceTemplate != nil {
				serviceTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.ServiceTemplate.ServiceId != nil {
					serviceTemplateMap["service_id"] = securityGroupPolicy.ServiceTemplate.ServiceId
				}

				if securityGroupPolicy.ServiceTemplate.ServiceGroupId != nil {
					serviceTemplateMap["service_group_id"] = securityGroupPolicy.ServiceTemplate.ServiceGroupId
				}

				securityGroupPolicyMap["service_template"] = []interface{}{serviceTemplateMap}
			}

			if securityGroupPolicy.CidrBlock != nil {
				securityGroupPolicyMap["cidr_block"] = securityGroupPolicy.CidrBlock
			}

			if securityGroupPolicy.Ipv6CidrBlock != nil {
				securityGroupPolicyMap["ipv6_cidr_block"] = securityGroupPolicy.Ipv6CidrBlock
			}

			if securityGroupPolicy.SecurityGroupId != nil {
				securityGroupPolicyMap["security_group_id"] = securityGroupPolicy.SecurityGroupId
			}

			if securityGroupPolicy.AddressTemplate != nil {
				addressTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.AddressTemplate.AddressId != nil {
					addressTemplateMap["address_id"] = securityGroupPolicy.AddressTemplate.AddressId
				}

				if securityGroupPolicy.AddressTemplate.AddressGroupId != nil {
					addressTemplateMap["address_group_id"] = securityGroupPolicy.AddressTemplate.AddressGroupId
				}

				securityGroupPolicyMap["address_template"] = []interface{}{addressTemplateMap}
			}

			if securityGroupPolicy.Action != nil {
				securityGroupPolicyMap["action"] = securityGroupPolicy.Action
			}

			if securityGroupPolicy.PolicyDescription != nil {
				securityGroupPolicyMap["policy_description"] = securityGroupPolicy.PolicyDescription
			}

			if securityGroupPolicy.ModifyTime != nil {
				securityGroupPolicyMap["modify_time"] = securityGroupPolicy.ModifyTime
			}

			ids = append(ids, *securityGroupPolicy.SecurityGroupId)
			tmpList = append(tmpList, securityGroupPolicyMap)
		}

		_ = d.Set("backup_data", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
