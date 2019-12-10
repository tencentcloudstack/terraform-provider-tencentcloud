/*
Provides a resource to create security group.

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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the security group.",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

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

	if tags := getTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("cvm", "sg", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	id := d.Id()

	securityGroup, err := vpcService.DescribeSecurityGroup(ctx, id)
	if err != nil {
		return err
	}

	if securityGroup == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", *securityGroup.SecurityGroupName)
	d.Set("description", *securityGroup.SecurityGroupDesc)

	projectId, err := strconv.Atoi(*securityGroup.ProjectId)
	if err != nil {
		return fmt.Errorf("securtiy group %s project id invalid: %v", *securityGroup.SecurityGroupId, err)
	}
	d.Set("project_id", projectId)

	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "sg", region, id)
	if err != nil {
		return err
	}
	d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := m.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	id := d.Id()

	d.Partial(true)

	var (
		newName    *string
		newDesc    *string
		attrUpdate []string
	)

	if d.HasChange("name") {
		newName = stringToPointer(d.Get("name").(string))
		attrUpdate = append(attrUpdate, "name")
	}

	if d.HasChange("description") {
		newDesc = stringToPointer(d.Get("description").(string))
		attrUpdate = append(attrUpdate, "description")
	}

	if len(attrUpdate) > 0 {
		if err := vpcService.ModifySecurityGroup(ctx, id, newName, newDesc); err != nil {
			return err
		}

		for _, attr := range attrUpdate {
			d.SetPartial(attr)
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("cvm", "sg", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudSecurityGroupRead(d, m)
}

func resourceTencentCloudSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_security_group.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	vpcService := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	// wait until all instances unbind this security group
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		associateSet, err := vpcService.DescribeSecurityGroupsAssociate(ctx, []string{id})
		if err != nil {
			return resource.RetryableError(err)
		}

		if len(associateSet) == 0 {
			return nil
		}

		statistics := associateSet[0]
		if *statistics.CVM > 0 {
			return resource.RetryableError(fmt.Errorf("security group %s still bind %d CVM instances", id, *statistics.CVM))
		}

		if *statistics.CLB > 0 {
			return resource.RetryableError(fmt.Errorf("security group %s still bind %d CLB instances", id, *statistics.CLB))
		}

		if *statistics.CDB > 0 {
			return resource.RetryableError(fmt.Errorf("security group %s still bind %d CDB instances", id, *statistics.CDB))
		}

		if *statistics.ENI > 0 {
			return resource.RetryableError(fmt.Errorf("security group %s still bind %d ENI instances", id, *statistics.ENI))
		}

		if *statistics.SG > 0 {
			return resource.RetryableError(fmt.Errorf("security group %s still bind %d SG instances", id, *statistics.SG))
		}

		return nil
	}); err != nil {
		return err
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		e := vpcService.DeleteSecurityGroup(ctx, id)
		if e != nil {
			return resource.RetryableError(fmt.Errorf("security group delete failed: %v", e))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s security group delete failed: %v", logId, err)
		return err
	}

	return nil
}
