/*
Use this data source to query detailed information of security group.

## Example Usage

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

	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the security group to be queried.",
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Name of the security group to be queried.",
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
				Description:  "Description of the security group.",
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
	}
}

func dataSourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	sgId := d.Get("security_group_id").(string)

	sg, has, err := service.DescribeSecurityGroup(ctx, sgId)
	if err != nil {
		return err
	}

	if has == 0 {
		return fmt.Errorf("security group %s not found", *sg.SecurityGroupId)
	}

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

	if associate == nil {
		return fmt.Errorf("security group %s associate statistic not found", *sg.SecurityGroupId)
	}

	var count int
	count += int(*associate.CVM)
	count += int(*associate.ENI)
	count += int(*associate.CDB)
	count += int(*associate.CLB)

	d.SetId(*sg.SecurityGroupId)
	_ = d.Set("security_group_id", *sg.SecurityGroupId)
	_ = d.Set("name", *sg.SecurityGroupName)
	_ = d.Set("description", *sg.SecurityGroupDesc)
	_ = d.Set("create_time", *sg.CreatedTime)
	_ = d.Set("be_associate_count", count)

	if sg.ProjectId != nil {
		projectId, err := strconv.Atoi(*sg.ProjectId)
		if err != nil {
			return fmt.Errorf("project id %s is invalid", *sg.ProjectId)
		}
		_ = d.Set("project_id", projectId)
	}

	return nil
}
