/*
Use this data source to query detailed information of security groups.

Example Usage

```hcl
data "tencentcloud_security_groups" "sglab" {
  security_group_id = "${tencentcloud_security_group.sglab.id}"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "project_id"},
				Description:   "ID of the security group to be queried. Conflict with `name` and `project_id`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateStringLengthInRange(1, 60),
				ConflictsWith: []string{"security_group_id"},
				Description:   "Name of the security group to be queried. Conflict with `security_group_id`.",
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"security_group_id"},
				Description:   "Project ID of the security group to be queried. Conflict with `security_group_id`.",
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"security_group_id"},
				Description:   "Tags of the security group to be queried. Conflict with `security_group_id`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of security group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the security group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the security group.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the security group.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of security group.",
						},
						"be_associate_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of security group binding resources.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID of the security group.",
						},
						"ingress": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Ingress rules set. For items like `[action]#[cidr_ip]#[port]#[protocol]`, it means a regular rule; for items like `sg-XXXX`, it means a nested security group.",
						},
						"egress": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Egress rules set. For items like `[action]#[cidr_ip]#[port]#[protocol]`, it means a regular rule; for items like `sg-XXXX`, it means a nested security group.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the security group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupsRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_security_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	var (
		sgId           *string
		sgName         *string
		inputProjectId *int
	)

	idBuilder := strings.Builder{}
	idBuilder.WriteString("securityGroups-")

	if raw, ok := d.GetOk("security_group_id"); ok {
		sgId = common.StringPtr(raw.(string))
		idBuilder.WriteString(*sgId)
		idBuilder.WriteRune('-')
	}

	if raw, ok := d.GetOk("name"); ok {
		sgName = common.StringPtr(raw.(string))
		idBuilder.WriteString(*sgName)
		idBuilder.WriteRune('-')
	}

	if raw, ok := d.GetOk("project_id"); ok {
		inputProjectId = common.IntPtr(raw.(int))
		idBuilder.WriteString(strconv.Itoa(*inputProjectId))
	}

	tags := helper.GetTags(d, "tags")

	sgs, err := vpcService.DescribeSecurityGroups(ctx, sgId, sgName, inputProjectId, tags)
	if err != nil {
		return err
	}

	if len(sgs) == 0 {
		_ = d.Set("security_groups", []map[string]interface{}{})
		d.SetId(idBuilder.String())
		return nil
	}

	sgMap := make(map[string]*vpc.SecurityGroup, len(sgs))
	for _, sg := range sgs {
		if sg.SecurityGroupId == nil {
			return errors.New("security group id is nil")
		}
		sgMap[*sg.SecurityGroupId] = sg
	}

	sgIds := make([]string, 0, len(sgs))
	for _, sg := range sgs {
		sgIds = append(sgIds, *sg.SecurityGroupId)
	}

	associateSet, err := vpcService.DescribeSecurityGroupsAssociate(ctx, sgIds)
	if err != nil {
		return err
	}

	sgInstances := make([]map[string]interface{}, 0, len(sgs))
	for _, associate := range associateSet {
		if associate.SecurityGroupId == nil {
			return errors.New("associate statistics security group id is nil")
		}

		if sg, ok := sgMap[*associate.SecurityGroupId]; ok {
			count := int(*associate.CVM + *associate.ENI + *associate.CDB + *associate.CLB)

			// normally projectId default value is 0
			if sg.ProjectId == nil {
				return errors.New("associate statistics project id is nil")
			}

			projectId, err := strconv.Atoi(*sg.ProjectId)
			if err != nil {
				return fmt.Errorf("securtiy group %s project id invalid: %v", *sg.SecurityGroupId, err)
			}

			respIngress, respEgress, exist, err := vpcService.DescribeSecurityGroupPolices(ctx, *sg.SecurityGroupId)
			if err != nil {
				return err
			}

			if !exist {
				// when read security group all rules, it doesn't exist, maybe delete on other places, ignore it
				continue
			}

			respTags, err := tagService.DescribeResourceTags(ctx, "cvm", "sg", region, *sg.SecurityGroupId)
			if err != nil {
				return err
			}

			ingress := make([]string, 0, len(respIngress))
			for _, in := range respIngress {
				ingress = append(ingress, in.String())
			}

			egress := make([]string, 0, len(respEgress))
			for _, eg := range respEgress {
				egress = append(egress, eg.String())
			}

			sgInstances = append(sgInstances, map[string]interface{}{
				"security_group_id":  *sg.SecurityGroupId,
				"name":               *sg.SecurityGroupName,
				"description":        *sg.SecurityGroupDesc,
				"create_time":        *sg.CreatedTime,
				"be_associate_count": count,
				"project_id":         projectId,
				"ingress":            ingress,
				"egress":             egress,
				"tags":               respTags,
			})
		}
	}

	if len(sgInstances) != len(sgs) {
		return errors.New("security group associate statistics is not enough")
	}

	_ = d.Set("security_groups", sgInstances)
	d.SetId(idBuilder.String())

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), sgInstances); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]", logId, output.(string), err)
			return err
		}
	}

	return nil
}
