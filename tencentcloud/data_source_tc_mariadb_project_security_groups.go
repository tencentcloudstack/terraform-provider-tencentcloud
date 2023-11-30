package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbProjectSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbProjectSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engine name. Valid value: `mariadb`.",
			},
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
			"groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Security group details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of yyyy-mm-dd hh:mm:ss.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group ID.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group name.",
						},
						"security_group_remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group remarks.",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inbound rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy, which can be `ACCEPT` or `DROP`.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source IP or source IP range, such as 192.168.0.0/16.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol. UDP and TCP are supported.",
									},
								},
							},
						},
						"outbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Outbound rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy, which can be `ACCEPT` or `DROP`.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source IP or source IP range, such as 192.168.0.0/16.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol. UDP and TCP are supported.",
									},
								},
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
		},
	}
}

func dataSourceTencentCloudMariadbProjectSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_project_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		groups  []*mariadb.SecurityGroup
		Product string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
		Product = v.(string)
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbProjectSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		groups = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(groups))

	if groups != nil {
		for _, securityGroup := range groups {
			securityGroupMap := map[string]interface{}{}

			if securityGroup.ProjectId != nil {
				securityGroupMap["project_id"] = securityGroup.ProjectId
			}

			if securityGroup.CreateTime != nil {
				securityGroupMap["create_time"] = securityGroup.CreateTime
			}

			if securityGroup.SecurityGroupId != nil {
				securityGroupMap["security_group_id"] = securityGroup.SecurityGroupId
			}

			if securityGroup.SecurityGroupName != nil {
				securityGroupMap["security_group_name"] = securityGroup.SecurityGroupName
			}

			if securityGroup.SecurityGroupRemark != nil {
				securityGroupMap["security_group_remark"] = securityGroup.SecurityGroupRemark
			}

			if securityGroup.Inbound != nil {
				inboundList := []interface{}{}
				for _, inbound := range securityGroup.Inbound {
					inboundMap := map[string]interface{}{}

					if inbound.Action != nil {
						inboundMap["action"] = inbound.Action
					}

					if inbound.CidrIp != nil {
						inboundMap["cidr_ip"] = inbound.CidrIp
					}

					if inbound.PortRange != nil {
						inboundMap["port_range"] = inbound.PortRange
					}

					if inbound.IpProtocol != nil {
						inboundMap["ip_protocol"] = inbound.IpProtocol
					}

					inboundList = append(inboundList, inboundMap)
				}

				securityGroupMap["inbound"] = inboundList
			}

			if securityGroup.Outbound != nil {
				outboundList := []interface{}{}
				for _, outbound := range securityGroup.Outbound {
					outboundMap := map[string]interface{}{}

					if outbound.Action != nil {
						outboundMap["action"] = outbound.Action
					}

					if outbound.CidrIp != nil {
						outboundMap["cidr_ip"] = outbound.CidrIp
					}

					if outbound.PortRange != nil {
						outboundMap["port_range"] = outbound.PortRange
					}

					if outbound.IpProtocol != nil {
						outboundMap["ip_protocol"] = outbound.IpProtocol
					}

					outboundList = append(outboundList, outboundMap)
				}

				securityGroupMap["outbound"] = outboundList
			}
			tmpList = append(tmpList, securityGroupMap)
		}

		_ = d.Set("groups", tmpList)
	}

	d.SetId(Product)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
