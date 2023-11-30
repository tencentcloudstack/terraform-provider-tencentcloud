package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamSetPolicyVersionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamSetPolicyVersionConfigCreate,
		Read:   resourceTencentCloudCamSetPolicyVersionConfigRead,
		Update: resourceTencentCloudCamSetPolicyVersionConfigUpdate,
		Delete: resourceTencentCloudCamSetPolicyVersionConfigDelete,
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

func resourceTencentCloudCamSetPolicyVersionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version_config.create")()
	defer inconsistentCheck(d, meta)()

	return resourceTencentCloudCamSetPolicyVersionConfigUpdate(d, meta)
}

func resourceTencentCloudCamSetPolicyVersionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version_config.read")()
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

	SetPolicyVersionConfig, err := service.DescribeCamSetPolicyVersionById(ctx, policyId, versionId)
	if err != nil {
		return err
	}

	if SetPolicyVersionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamSetPolicyVersionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("policy_id", helper.StrToInt(policyId))

	if SetPolicyVersionConfig.VersionId != nil {
		_ = d.Set("version_id", SetPolicyVersionConfig.VersionId)
	}

	return nil
}

func resourceTencentCloudCamSetPolicyVersionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cam.NewSetDefaultPolicyVersionRequest()

	var policyId string
	var versionId string

	policyId = helper.IntToStr(d.Get("policy_id").(int))
	versionId = helper.IntToStr(d.Get("version_id").(int))
	request.PolicyId = helper.StrToUint64Point(policyId)
	request.VersionId = helper.StrToUint64Point(versionId)

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
		log.Printf("[CRITAL]%s update cam SetPolicyVersionConfig failed, reason:%+v", logId, err)
		return err
	}
	d.SetId(policyId + FILED_SP + versionId)

	return resourceTencentCloudCamSetPolicyVersionConfigRead(d, meta)
}

func resourceTencentCloudCamSetPolicyVersionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_set_policy_version_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
