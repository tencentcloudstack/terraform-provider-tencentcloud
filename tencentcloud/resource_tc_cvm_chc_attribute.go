/*
Provides a resource to create a cvm chc_attribute

Example Usage

```hcl
resource "tencentcloud_cvm_chc_attribute" "chc_attribute" {
	chc_id = "chc-xxxxx"
	instance_name = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmChcAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcAttributeCreate,
		Read:   resourceTencentCloudCvmChcAttributeRead,
		Delete: resourceTencentCloudCvmChcAttributeDelete,

		Schema: map[string]*schema.Schema{
			"chc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CHC host ID.",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CHC host name.",
			},

			"device_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Server type.",
			},

			"bmc_user": {
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"password"},
				Type:         schema.TypeString,
				Description:  "Valid characters: Letters, numbers, hyphens and underscores.",
			},

			"password": {
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				RequiredWith: []string{"bmc_user"},
				Type:         schema.TypeString,

				Description: "The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&amp;amp;*-+=_|{}).",
			},

			"bmc_security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "BMC network security group list.",
			},
		},
	}
}

func resourceTencentCloudCvmChcAttributeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cvm.NewModifyChcAttributeRequest()
		chcId   string
	)
	chcId = d.Get("chc_id").(string)
	request.ChcIds = []*string{&chcId}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("device_type"); ok {
		request.DeviceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bmc_user"); ok {
		request.BmcUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bmc_security_group_ids"); ok {
		bmcSecurityGroupIdsSet := v.(*schema.Set).List()
		for i := range bmcSecurityGroupIdsSet {
			bmcSecurityGroupIds := bmcSecurityGroupIdsSet[i].(string)
			request.BmcSecurityGroupIds = append(request.BmcSecurityGroupIds, &bmcSecurityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyChcAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(chcId)

	return resourceTencentCloudCvmChcAttributeRead(d, meta)
}

func resourceTencentCloudCvmChcAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	chcId := d.Id()

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}

	if len(chcHosts) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s host` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	chcHost := chcHosts[0]
	_ = d.Set("instance_name", chcHost.InstanceName)
	_ = d.Set("device_type", chcHost.DeviceType)

	bmcSecurityGroupIds := make([]string, 0)
	if len(chcHost.BmcSecurityGroupIds) > 0 {
		for _, sg := range chcHost.BmcSecurityGroupIds {
			bmcSecurityGroupId := *sg
			bmcSecurityGroupIds = append(bmcSecurityGroupIds, bmcSecurityGroupId)
		}
	}
	_ = d.Set("bmc_security_group_ids", bmcSecurityGroupIds)

	return nil
}

func resourceTencentCloudCvmChcAttributeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_chc_attribute.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
