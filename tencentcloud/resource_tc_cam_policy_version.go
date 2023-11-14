/*
Provides a resource to create a cam policy_version

Example Usage

```hcl
resource "tencentcloud_cam_policy_version" "policy_version" {
  policy_id =
  policy_document = ""
  set_as_default =
}
```

Import

cam policy_version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_policy_version.policy_version policy_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCamPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamPolicyVersionCreate,
		Read:   resourceTencentCloudCamPolicyVersionRead,
		Update: resourceTencentCloudCamPolicyVersionUpdate,
		Delete: resourceTencentCloudCamPolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},

			"policy_document": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Policy Text Information.",
			},

			"set_as_default": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Set as the version of the current policy.",
			},
		},
	}
}

func resourceTencentCloudCamPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cam.NewCreatePolicyVersionRequest()
		response  = cam.NewCreatePolicyVersionResponse()
		policyId  int
		versionId int
	)
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int)
		request.PolicyId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy_document"); ok {
		request.PolicyDocument = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("set_as_default"); ok {
		request.SetAsDefault = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreatePolicyVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam PolicyVersion failed, reason:%+v", logId, err)
		return err
	}

	policyId = *response.Response.PolicyId
	d.SetId(strings.Join([]string{helper.Int64ToStr(policyId), helper.Int64ToStr(int64(versionId))}, FILED_SP))

	return resourceTencentCloudCamPolicyVersionRead(d, meta)
}

func resourceTencentCloudCamPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	versionId := idSplit[1]

	PolicyVersion, err := service.DescribeCamPolicyVersionById(ctx, policyId, versionId)
	if err != nil {
		return err
	}

	if PolicyVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamPolicyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if PolicyVersion.PolicyId != nil {
		_ = d.Set("policy_id", PolicyVersion.PolicyId)
	}

	if PolicyVersion.PolicyDocument != nil {
		_ = d.Set("policy_document", PolicyVersion.PolicyDocument)
	}

	if PolicyVersion.SetAsDefault != nil {
		_ = d.Set("set_as_default", PolicyVersion.SetAsDefault)
	}

	return nil
}

func resourceTencentCloudCamPolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"policy_id", "policy_document", "set_as_default"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudCamPolicyVersionRead(d, meta)
}

func resourceTencentCloudCamPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_version.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CamService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	versionId := idSplit[1]

	if err := service.DeleteCamPolicyVersionById(ctx, policyId, versionId); err != nil {
		return err
	}

	return nil
}
