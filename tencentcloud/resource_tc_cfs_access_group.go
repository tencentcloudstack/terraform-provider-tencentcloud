/*
Provides a resource to create a CFS access group.

Example Usage

```hcl
resource "tencentcloud_cfs_access_group" "foo" {
  name        = "test_access_group"
  description = "test"
}
```

Import

CFS access group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_access_group.foo pgroup-7nx89k7l
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCfsAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsAccessGroupCreate,
		Read:   resourceTencentCloudCfsAccessGroupRead,
		Update: resourceTencentCloudCfsAccessGroupUpdate,
		Delete: resourceTencentCloudCfsAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Name of the access group, and max length is 64.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 255),
				Description:  "Description of the access group, and max length is 255.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the access group.",
			},
		},
	}
}

func resourceTencentCloudCfsAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_group.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	name := d.Get("name").(string)
	description := ""
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	accessGroupId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		id, errRet := cfsService.CreateAccessGroup(ctx, name, description)
		if errRet != nil {
			return retryError(errRet)
		}
		accessGroupId = id
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(accessGroupId)

	return resourceTencentCloudCfsAccessGroupRead(d, meta)
}

func resourceTencentCloudCfsAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	id := d.Id()
	var accessGroup *cfs.PGroupInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		accessGroups, errRet := cfsService.DescribeAccessGroup(ctx, id, "")
		if errRet != nil {
			return retryError(errRet)
		}
		if len(accessGroups) > 0 {
			accessGroup = accessGroups[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if accessGroup == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", accessGroup.Name)
	_ = d.Set("description", accessGroup.DescInfo)
	_ = d.Set("create_time", accessGroup.CDate)

	return nil
}

func resourceTencentCloudCfsAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_group.update")()
	logId := getLogId(contextNil)

	request := cfs.NewUpdateCfsPGroupRequest()
	if d.HasChange("name") {
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		request.DescInfo = helper.String(d.Get("description").(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCfsClient().UpdateCfsPGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudCfsAccessGroupRead(d, meta)
}

func resourceTencentCloudCfsAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfs_access_group.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	cfsService := CfsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cfsService.DeleteAccessGroup(ctx, id)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
