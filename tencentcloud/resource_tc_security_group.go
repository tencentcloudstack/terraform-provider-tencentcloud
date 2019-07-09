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
				ValidateFunc: validateStringLengthInRange(2, 60),
			},

			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 100),
			},

			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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

	var projectID *string
	if projectIDInterface, exist := d.GetOk("project_id"); exist {
		projectIDStr := projectIDInterface.(string)
		projectID = &projectIDStr
	}

	id, err := vpcService.CreateSecurityGroup(ctx, name, desc, projectID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] SgId=%s", id)
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
		_ = d.Set("name", *securityGroup.SecurityGroupName)
		_ = d.Set("description", *securityGroup.SecurityGroupDesc)
		if securityGroup.ProjectId != nil {
			projectID, _ := strconv.Atoi(*securityGroup.ProjectId)
			_ = d.Set("project_id", projectID)
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
