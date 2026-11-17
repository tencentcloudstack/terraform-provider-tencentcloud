package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDbInstanceSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance ID, which can be obtained via the DescribeDBInstances API. DBInstanceId, ReadOnlyGroupId at least pass one; if you want to query the security group associated with the instance, only pass the DBInstanceId field.",
			},

			"read_only_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Read-only group ID, which can be obtained via the DescribeReadOnlyGroups API. DBInstanceId, ReadOnlyGroupId at least pass one; if you want to query the security group associated with the read-only group, only pass the ReadOnlyGroupId.",
			},

			"security_group_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security group information list of the instance.",
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
							Description: "Creation time.",
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
						"security_group_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group remark.",
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
										Description: "Policy, ACCEPT or DROP.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source or destination IP or IP range, e.g. 172.16.0.0/12.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol, supports UDP, TCP, etc.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
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
										Description: "Policy, ACCEPT or DROP.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source or destination IP or IP range, e.g. 172.16.0.0/12.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol, supports UDP, TCP, etc.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlDbInstanceSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_db_instance_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	dbInstanceId := ""
	if v, ok := d.GetOk("db_instance_id"); ok {
		dbInstanceId = v.(string)
	}

	readOnlyGroupId := ""
	if v, ok := d.GetOk("read_only_group_id"); ok {
		readOnlyGroupId = v.(string)
	}

	var securityGroups []*postgresql.SecurityGroup
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDbInstanceSecurityGroups(ctx, dbInstanceId, readOnlyGroupId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("describe postgresql_db_instance_security_groups failed, response is nil or empty."))
		}

		securityGroups = result
		return nil
	})
	if err != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return err
	}

	ids := make([]string, 0, len(securityGroups))
	tmpList := make([]map[string]interface{}, 0, len(securityGroups))
	for _, securityGroup := range securityGroups {
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

		if securityGroup.SecurityGroupDescription != nil {
			securityGroupMap["security_group_description"] = securityGroup.SecurityGroupDescription
		}

		if securityGroup.Inbound != nil {
			inboundList := make([]map[string]interface{}, 0, len(securityGroup.Inbound))
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

				if inbound.Description != nil {
					inboundMap["description"] = inbound.Description
				}

				inboundList = append(inboundList, inboundMap)
			}
			securityGroupMap["inbound"] = inboundList
		}

		if securityGroup.Outbound != nil {
			outboundList := make([]map[string]interface{}, 0, len(securityGroup.Outbound))
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

				if outbound.Description != nil {
					outboundMap["description"] = outbound.Description
				}

				outboundList = append(outboundList, outboundMap)
			}
			securityGroupMap["outbound"] = outboundList
		}

		if securityGroup.SecurityGroupId != nil {
			ids = append(ids, *securityGroup.SecurityGroupId)
		}
		tmpList = append(tmpList, securityGroupMap)
	}

	_ = d.Set("security_group_set", tmpList)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
