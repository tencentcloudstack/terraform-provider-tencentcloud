/*
Provide a resource to create a placement group.

Example Usage

```hcl
resource "tencentcloud_placement_group" "foo" {
  name = "test"
  type = "HOST"
}
```

Import

Placement group can be imported using the id, e.g.

```
$ terraform import tencentcloud_placement_group.foo ps-ilan8vjf
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPlacementGroupCreate,
		Read:   resourceTencentCloudPlacementGroupRead,
		Update: resourceTencentCloudPlacementGroupUpdate,
		Delete: resourceTencentCloudPlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the placement group, 1-60 characters in length.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_PLACEMENT_GROUP_TYPE),
				Description:  "Type of the placement group, the available values include `HOST`,`SW` and `RACK`.",
			},

			// computed
			"cvm_quota_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of hosts in the placement group.",
			},
			"current_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of hosts in the placement group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the placement group.",
			},
		},
	}
}

func resourceTencentCloudPlacementGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_placement_group.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	placementName := d.Get("name").(string)
	placementType := d.Get("type").(string)
	var id string
	var errRet error
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		id, errRet = cvmService.CreatePlacementGroup(ctx, placementName, placementType)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_placement_group.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var placement *cvm.DisasterRecoverGroup
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		placement, errRet = cvmService.DescribePlacementGroupById(ctx, placementId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if placement == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", placement.Name)
	d.Set("type", placement.Type)
	d.Set("cvm_quota_total", placement.CvmQuotaTotal)
	d.Set("current_num", placement.CurrentNum)
	d.Set("create_time", placement.CreateTime)

	return nil
}

func resourceTencentCloudPlacementGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_placement_group.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	if d.HasChange("name") {
		placementName := d.Get("name").(string)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err := cvmService.ModifyPlacementGroup(ctx, placementId, placementName)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_placement_group.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := cvmService.DeletePlacementGroup(ctx, placementId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
