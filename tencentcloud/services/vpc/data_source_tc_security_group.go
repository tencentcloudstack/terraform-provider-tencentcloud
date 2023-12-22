package vpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func DataSourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.14.0. Please use 'tencentcloud_security_groups' instead.",
		Read:               dataSourceTencentCloudSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "ID of the security group to be queried. Conflict with `name`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  tccommon.ValidateStringLengthInRange(1, 60),
				ConflictsWith: []string{"security_group_id"},
				Description:   "Name of the security group to be queried. Conflict with `security_group_id`.",
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
		},
	}
}

func dataSourceTencentCloudSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_security_group.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := VpcService{client: client}

	var (
		sgId   *string
		sgName *string
	)

	if raw, ok := d.GetOk("security_group_id"); ok {
		sgId = common.StringPtr(raw.(string))
	}

	if raw, ok := d.GetOk("name"); ok {
		sgName = common.StringPtr(raw.(string))
	}

	sgs, err := vpcService.DescribeSecurityGroups(ctx, sgId, sgName, nil, map[string]string{})
	if err != nil {
		return err
	}

	if len(sgs) == 0 {
		return errors.New("security group not found with the query condition")
	}

	sg := sgs[0]
	in, out, _, err := vpcService.DescribeSecurityGroupPolices(ctx, *sg.SecurityGroupId)
	if err != nil {
		return err
	}

	d.SetId(*sg.SecurityGroupId)
	_ = d.Set("security_group_id", sg.SecurityGroupId)
	_ = d.Set("name", sg.SecurityGroupName)
	_ = d.Set("description", sg.SecurityGroupDesc)
	_ = d.Set("create_time", sg.CreatedTime)
	_ = d.Set("be_associate_count", len(in)+len(out))

	projectId, err := strconv.Atoi(*sg.ProjectId)
	if err != nil {
		return fmt.Errorf("project id is not valid number: %v", err)
	}

	_ = d.Set("project_id", projectId)

	return nil
}
