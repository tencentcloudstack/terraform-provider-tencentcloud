package dcdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "security group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group id.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group name.",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "inbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cidr ip.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "policy action.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "internet protocol.",
									},
								},
							},
						},
						"outbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "outbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cidr ip.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "policy action.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "internet protocol.",
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

func dataSourceTencentCloudDcdbSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var groups []*dcdb.SecurityGroup
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groups = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb groups failed, reason:%+v", logId, err)
		return err
	}

	groupList := []interface{}{}
	ids := make([]string, 0, len(groups))
	if groups != nil {
		for _, group := range groups {
			groupMap := map[string]interface{}{}
			if group.ProjectId != nil {
				groupMap["project_id"] = group.ProjectId
			}
			if group.CreateTime != nil {
				groupMap["create_time"] = group.CreateTime
			}
			if group.SecurityGroupId != nil {
				groupMap["security_group_id"] = group.SecurityGroupId
			}
			if group.SecurityGroupName != nil {
				groupMap["security_group_name"] = group.SecurityGroupName
			}
			if group.Inbound != nil {
				inboundList := []interface{}{}
				for _, inbound := range group.Inbound {
					inboundMap := map[string]interface{}{}
					if inbound.CidrIp != nil {
						inboundMap["cidr_ip"] = inbound.CidrIp
					}
					if inbound.Action != nil {
						inboundMap["action"] = inbound.Action
					}
					if inbound.PortRange != nil {
						inboundMap["port_range"] = inbound.PortRange
					}
					if inbound.IpProtocol != nil {
						inboundMap["ip_protocol"] = inbound.IpProtocol
					}

					inboundList = append(inboundList, inboundMap)
				}
				groupMap["inbound"] = inboundList
			}
			if group.Outbound != nil {
				outboundList := []interface{}{}
				for _, outbound := range group.Outbound {
					outboundMap := map[string]interface{}{}
					if outbound.CidrIp != nil {
						outboundMap["cidr_ip"] = outbound.CidrIp
					}
					if outbound.Action != nil {
						outboundMap["action"] = outbound.Action
					}
					if outbound.PortRange != nil {
						outboundMap["port_range"] = outbound.PortRange
					}
					if outbound.IpProtocol != nil {
						outboundMap["ip_protocol"] = outbound.IpProtocol
					}

					outboundList = append(outboundList, outboundMap)
				}
				groupMap["outbound"] = outboundList
			}
			ids = append(ids, *group.SecurityGroupId)
			groupList = append(groupList, groupMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", groupList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupList); e != nil {
			return e
		}
	}

	return nil
}
