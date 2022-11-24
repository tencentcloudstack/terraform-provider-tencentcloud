/*
Provides a resource to create a CAM policy.

Example Usage

```hcl
resource "tencentcloud_cam_policy_by_name" "foo" {
  name        = "cam-policy-test"
  document    = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": [
        "name/sts:AssumeRole"
      ],
      "effect": "allow",
      "resource": [
        "*"
      ]
    }
  ]
}
EOF
  description = "test"
}
```

Import

CAM policy can be imported using the name, e.g.

```
$ terraform import tencentcloud_cam_policy_by_name.foo name
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamPolicyByName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamPolicyByNameCreate,
		Read:   resourceTencentCloudCamPolicyByNameRead,
		Update: resourceTencentCloudCamPolicyByNameUpdate,
		Delete: resourceTencentCloudCamPolicyByNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of CAM policy.",
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					var oldJson interface{}
					err := json.Unmarshal([]byte(olds), &oldJson)
					if err != nil {
						return olds == news
					}
					var newJson interface{}
					err = json.Unmarshal([]byte(news), &newJson)
					if err != nil {
						return olds == news
					}
					flag := reflect.DeepEqual(oldJson, newJson)
					return flag
				},
				Description: "Document of the CAM policy. The syntax refers to [CAM POLICY](https://intl.cloud.tencent.com/document/product/598/10604). There are some notes when using this para in terraform: 1. The elements in JSON claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when it appears, it must be replaced with the uin it stands for.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CAM policy.",
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of the policy strategy. Valid values: `1`, `2`.  `1` means customer strategy and `2` means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM policy.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the CAM policy.",
			},
		},
	}
}

func resourceTencentCloudCamPolicyByNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_by_name.create")()

	logId := getLogId(contextNil)

	name := d.Get("name").(string)
	document := d.Get("document").(string)

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	documentErr := camService.PolicyDocumentForceCheck(document)
	if documentErr != nil {
		return documentErr
	}
	request := cam.NewCreatePolicyRequest()
	request.PolicyName = &name
	request.PolicyDocument = &document
	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	var response *cam.CreatePolicyResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreatePolicy(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "PolicyNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.PolicyId == nil {
		return fmt.Errorf("CAM policy id is nil")
	}
	d.SetId(name)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		parmas := make(map[string]interface{})
		parmas["name"] = name
		instances, e := camService.DescribePoliciesByFilter(ctx, parmas)
		if e != nil {
			return retryError(e)
		}
		if len(instances) == 0 {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamPolicyByNameRead(d, meta)
}

func resourceTencentCloudCamPolicyByNameRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_by_name.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	policyName := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var policies []*cam.StrategyInfo
	params := make(map[string]interface{})
	params["name"] = policyName
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var innerErr error
		policies, innerErr = camService.DescribePoliciesByFilter(ctx, params)
		if innerErr != nil {
			return retryError(innerErr, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(policies) == 0 {
		return nil
	}
	var instance *cam.GetPolicyResponse
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		policyId := strconv.Itoa(int(*policies[0].PolicyId))
		result, e := camService.DescribePolicyById(ctx, policyId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.PolicyName == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", *instance.Response.PolicyName)
	//document with special change rule, the `\\/` must be replaced with `/`
	_ = d.Set("document", strings.Replace(*instance.Response.PolicyDocument, "\\/", "/", -1))
	_ = d.Set("create_time", *instance.Response.AddTime)
	_ = d.Set("update_time", *instance.Response.UpdateTime)
	_ = d.Set("type", int(*instance.Response.Type))
	if instance.Response.Description != nil {
		_ = d.Set("description", *instance.Response.Description)
	}
	return nil
}

func resourceTencentCloudCamPolicyByNameUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_by_name.update")()

	logId := getLogId(contextNil)

	policyName := d.Id()
	request := cam.NewUpdatePolicyRequest()
	request.PolicyName = &policyName
	changeFlag := false

	if d.HasChange("description") {
		request.Description = helper.String(d.Get("description").(string))
		changeFlag = true

	}
	if d.HasChange("name") {
		request.PolicyName = helper.String(d.Get("name").(string))
		changeFlag = true
	}

	if d.HasChange("document") {
		document := d.Get("document").(string)
		camService := CamService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		documentErr := camService.PolicyDocumentForceCheck(document)
		if documentErr != nil {
			return documentErr
		}
		request.PolicyDocument = &document
		changeFlag = true

	}
	if changeFlag {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdatePolicy(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM policy description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCamPolicyByNameRead(d, meta)
}

func resourceTencentCloudCamPolicyByNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy_by_name.delete")()

	logId := getLogId(contextNil)
	var policies []*cam.StrategyInfo
	policyName := d.Id()
	params := make(map[string]interface{})
	params["name"] = policyName
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var innerErr error
		policies, innerErr = camService.DescribePoliciesByFilter(ctx, params)
		if innerErr != nil {
			return retryError(innerErr, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(policies) == 0 {
		return nil
	}

	policyId := policies[0].PolicyId
	request := cam.NewDeletePolicyRequest()
	request.PolicyId = []*uint64{policyId}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DeletePolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
