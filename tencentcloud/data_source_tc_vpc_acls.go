/*
Use this data source to query VPC Network ACL information.

Example Usage

```hcl
data "tencentcloud_vpc_instances" "foo" {
}

data "tencentcloud_vpc_acls" "foo" {
  vpc_id            = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
}

data "tencentcloud_vpc_acls" "foo" {
  name            	= "test_acl"
}

```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcACLRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "ID of the VPC instance.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 60),
				Description:  "Name of the network ACL.",
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "`ID` of the network ACL instance.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"acl_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The information list of the VPC. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`ID` of the VPC instance.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`ID` of the network ACL instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the network ACL.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subnets associated with the network ACL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the VPC instance.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet instance `ID`.",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet name.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The `IPv4` `CIDR` of the subnet.",
									},
									"tags": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "Tags of the subnet.",
									},
								},
							},
						},
						"ingress": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inbound rules of the network ACL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of ip protocol.",
									},
									"port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Range of the port.",
									},
									"policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule policy of.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IP address network or segment.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
									},
								},
							},
						},
						"egress": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Outbound rules of the network ACL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of ip protocol.",
									},
									"port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Range of the port.",
									},
									"policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule policy of.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IP address network or segment.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
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

func dataSourceTencentCloudVpcACLRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_acls.read")()
	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

		vpcID string
		name  string
		id    string
	)

	if temp, ok := d.GetOk("vpc_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			vpcID = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}
	if temp, ok := d.GetOk("id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			id = tempStr
		}
	}

	networkAcls, err := service.DescribeNetWorkAcls(ctx, id, vpcID, name)
	if err != nil {
		return err
	}

	aclList := make([]map[string]interface{}, 0, len(networkAcls))
	ids := make([]string, 0, len(networkAcls))

	for _, info := range networkAcls {
		subnetInfo := info.SubnetSet
		subnets := make([]map[string]interface{}, 0, len(subnetInfo))
		for i := range subnetInfo {
			v := subnetInfo[i]
			subnet := make(map[string]interface{}, 5)
			subnet["vpc_id"] = v.VpcId
			subnet["subnet_id"] = v.SubnetId
			subnet["subnet_name"] = v.SubnetName
			subnet["cidr_block"] = v.CidrBlock

			tag := make(map[string]interface{}, len(v.TagSet))
			for t := range v.TagSet {
				tagValue := v.TagSet[t]
				tag[*tagValue.Key] = tagValue.Value
			}
			subnet["tags"] = tag

			subnets = append(subnets, subnet)
		}

		ingressInfo := info.IngressEntries
		ingress := make([]map[string]interface{}, 0, len(ingressInfo))
		for i := range ingressInfo {
			v := ingressInfo[i]
			egressMap := make(map[string]interface{}, 5)
			egressMap["protocol"] = v.Protocol
			egressMap["port"] = v.Port
			egressMap["cidr_block"] = v.CidrBlock
			egressMap["policy"] = v.Action
			egressMap["description"] = v.Description

			ingress = append(ingress, egressMap)
		}

		egressInfo := info.EgressEntries
		egress := make([]map[string]interface{}, 0, len(egressInfo))
		for i := range egressInfo {
			v := egressInfo[i]
			egressMap := make(map[string]interface{}, 5)
			egressMap["protocol"] = v.Protocol
			egressMap["port"] = v.Port
			egressMap["cidr_block"] = v.CidrBlock
			egressMap["policy"] = v.Action
			egressMap["description"] = v.Description

			egress = append(egress, egressMap)
		}

		aclResult := map[string]interface{}{
			"vpc_id":      info.VpcId,
			"id":          info.NetworkAclId,
			"name":        info.NetworkAclName,
			"create_time": info.CreatedTime,
			"subnets":     subnets,
			"ingress":     ingress,
			"egress":      egress,
		}
		aclList = append(aclList, aclResult)
		ids = append(ids, *info.NetworkAclId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("acl_list", aclList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set acl list fail, reason:%v \n ", logId, err)
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), aclList); err != nil {
			return err
		}
	}
	return nil
}
