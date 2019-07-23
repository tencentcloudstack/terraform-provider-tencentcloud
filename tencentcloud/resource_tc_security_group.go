/*
Provide a resource to create security group.

Example Usage

```hcl
resource "tencentcloud_security_group" "sglab" {
  name        = "mysg"
  description = "favourite sg"
  project_id  = 0
}
```

Import

Security group can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group.sglab sg-ey3wmiz1
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func resourceTencentCloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupCreate,
		Read:   resourceTencentCloudSecurityGroupRead,
		Update: resourceTencentCloudSecurityGroupUpdate,
		Delete: resourceTencentCloudSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the security group to be queried.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "Description of the security group.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Project ID of the security group.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	name := d.Get("name").(string)
	desc := d.Get("description").(string)

	var projectId *int
	if projectIdInterface, exist := d.GetOk("project_id"); exist {
		projectId = common.IntPtr(projectIdInterface.(int))
	}

	id, err := vpcService.CreateSecurityGroup(ctx, name, desc, projectId)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	securityGroup, has, err := vpcService.DescribeSecurityGroup(ctx, id)
	if err != nil {
		return err
	}

	switch has {
	default:
		err := fmt.Errorf("one security_group_id read get %d security_group info", has)
		log.Printf("[CRITAL]%s %v", logId, err)

		return err

	case 0:
		d.SetId("")
		return nil

	case 1:
		d.Set("name", *securityGroup.SecurityGroupName)
		d.Set("description", *securityGroup.SecurityGroupDesc)
		if securityGroup.ProjectId != nil {
			projectId, err := strconv.Atoi(*securityGroup.ProjectId)
			if err != nil {
				return fmt.Errorf("securtiy group %s project id invalid: %v", *securityGroup.SecurityGroupId, err)
			}
			d.Set("project_id", projectId)
		}
	}

	return nil
}

func resourceTencentCloudSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	attributeUpdate := d.HasChange("name") || d.HasChange("description")
	var (
		newName *string
		newDesc *string
	)

	if d.HasChange("name") {
		newName = common.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		newDesc = common.StringPtr(d.Get("description").(string))
	}

	if !attributeUpdate {
		return nil
	}

	if err := vpcService.ModifySecurityGroup(ctx, id, newName, newDesc); err != nil {
		return err
	}

	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_security_group.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return vpcService.DeleteSecurityGroup(ctx, id)
}
