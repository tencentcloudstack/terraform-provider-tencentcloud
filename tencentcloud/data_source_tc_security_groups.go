/*
Use this data source to query detailed information of security group.

Example Usage

```hcl
data "tencentcloud_security_group" "sglab" {
    security_group_id = "sg-fh48e762"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "project_id"},
				Description:   "ID of the security group to be queried.",
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateStringLengthInRange(1, 60),
				ConflictsWith: []string{"security_group_id"},
				Description:   "Name of the security group to be queried.",
			},

			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"security_group_id"},
			},

			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupsRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_groups.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		sgId      *string
		sgName    *string
		projectId *int
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
		projectId = common.IntPtr(raw.(int))
		idBuilder.WriteString(strconv.Itoa(*projectId))
	}

	sgs, err := service.DescribeSecurityGroups(ctx, sgId, sgName, projectId)
	if err != nil {
		return err
	}

	sgInstances := make([]map[string]interface{}, 0, len(sgs))

	for _, sg := range sgs {
		associateSet, err := service.DescribeSecurityGroupsAssociate(ctx, []string{*sg.SecurityGroupId})
		if err != nil {
			return err
		}

		var associate *vpc.SecurityGroupAssociationStatistics
		for _, v := range associateSet {
			if *v.SecurityGroupId == *sg.SecurityGroupId {
				associate = v
				break
			}
		}

		var count int
		if associate != nil {
			count += int(*associate.CVM)
			count += int(*associate.ENI)
			count += int(*associate.CDB)
			count += int(*associate.CLB)
		}

		pjId, err := strconv.Atoi(*sg.ProjectId)
		if err != nil {
			return fmt.Errorf("securtiy group %s project id invalid: %v", *sg.SecurityGroupId, err)
		}

		sgInstances = append(sgInstances, map[string]interface{}{
			"id":                 *sg.SecurityGroupId,
			"name":               *sg.SecurityGroupName,
			"description":        *sg.SecurityGroupDesc,
			"create_time":        *sg.CreatedTime,
			"be_associate_count": count,
			"project_id":         pjId,
		})
	}

	_ = d.Set("security_groups", sgInstances)

	d.SetId(idBuilder.String())

	return nil
}
