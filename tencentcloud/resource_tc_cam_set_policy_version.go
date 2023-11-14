/*
Provides a resource to create a cam set_policy_version

Example Usage

```hcl
resource "tencentcloud_cam_set_policy_version" "set_policy_version" {
  policy_id =
  version_id =
}
```

Import

cam set_policy_version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_set_policy_version.set_policy_version set_policy_version_id
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

func resourceTencentCloudCamSetPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamSetPolicyVersionCreate,
		Read:   resourceTencentCloudCamSetPolicyVersionRead,
		Update: resourceTencentCloudCamSetPolicyVersionUpdate,
		Delete: resourceTencentCloudCamSetPolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID.",
			},

			"version_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The policy version number, which can be obtained from ListPolicyVersions.",
			},
		},
	}
}

func resourceTencentCloudCamSetPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version.create")()
	defer inconsistentCheck(d, meta)()

	var policyId int64
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int64)
	}

	var versionId uint64
	if v, ok := d.GetOkExists("version_id"); ok {
		versionId = v.(uint64)
	}

	d.SetId(strings.Join([]string{helper.Int64ToStr(policyId), helper.Int64ToStr(int64(versionId))}, FILED_SP))

	return resourceTencentCloudCamSetPolicyVersionUpdate(d, meta)
}

func resourceTencentCloudCamSetPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version.read")()
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

	SetPolicyVersion, err := service.DescribeCamSetPolicyVersionById(ctx, policyId, versionId)
	if err != nil {
		return err
	}

	if SetPolicyVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamSetPolicyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if SetPolicyVersion.PolicyId != nil {
		_ = d.Set("policy_id", SetPolicyVersion.PolicyId)
	}

	if SetPolicyVersion.VersionId != nil {
		_ = d.Set("version_id", SetPolicyVersion.VersionId)
	}

	return nil
}

func resourceTencentCloudCamSetPolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewSetDefaultPolicyVersionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	versionId := idSplit[1]

	request.PolicyId = &policyId
	request.VersionId = &versionId

	immutableArgs := []string{"policy_id", "version_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().SetDefaultPolicyVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cam SetPolicyVersion failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCamSetPolicyVersionRead(d, meta)
}

func resourceTencentCloudCamSetPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
