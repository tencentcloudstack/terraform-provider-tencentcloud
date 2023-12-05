package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCamRoleByName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRoleByNameCreate,
		Read:   resourceTencentCloudCamRoleByNameRead,
		Update: resourceTencentCloudCamRoleByNameUpdate,
		Delete: resourceTencentCloudCamRoleByNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of CAM role.",
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
				Description: "Document of the CAM role. The syntax refers to [CAM POLICY](https://intl.cloud.tencent.com/document/product/598/10604). There are some notes when using this para in terraform: 1. The elements in json claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when appears, it must be replaced with the uin it stands for.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CAM role.",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether the CAM role can login or not.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM role.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the CAM role.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},
		},
	}
}

func resourceTencentCloudCamRoleByNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_by_name.create")()

	logId := getLogId(contextNil)

	name := d.Get("name").(string)
	document := d.Get("document").(string)

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cam.NewCreateRoleRequest()
	request.RoleName = &name
	request.PolicyDocument = &document
	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		loginBool := v.(bool)
		loginInt := uint64(1)
		if !loginBool {
			loginInt = uint64(0)
		}
		request.ConsoleLogin = &loginInt
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreateRole(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "RoleNameInUse") {
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(name)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instances []*cam.RoleInfo
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		params := make(map[string]interface{})
		params["name"] = name
		var innerErr error
		instances, innerErr = camService.DescribeRolesByFilter(ctx, params)
		if innerErr != nil {
			return retryError(innerErr)
		}
		if len(instances) == 0 {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		var instance *cam.RoleInfo
		if len(instances) != 0 {
			instance = instances[0]
		}
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		resourceName := BuildTagResourceName("cam", "role", "", *instance.RoleId)
		log.Printf("resourceName: %v", resourceName)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamRoleByNameRead(d, meta)
}

func resourceTencentCloudCamRoleByNameRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_by_name.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleName := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cam.RoleInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		params := make(map[string]interface{})
		params["name"] = roleName
		instances, e := camService.DescribeRolesByFilter(ctx, params)

		if e != nil {
			return retryError(e)
		}
		if len(instances) != 0 {
			instance = instances[0]
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", instance.RoleName)
	_ = d.Set("document", instance.PolicyDocument)
	_ = d.Set("create_time", instance.AddTime)
	_ = d.Set("update_time", instance.UpdateTime)
	if instance.Description != nil {
		_ = d.Set("description", instance.Description)
	}

	if instance.ConsoleLogin != nil {
		if int(*instance.ConsoleLogin) == 1 {
			_ = d.Set("console_login", true)
		} else {
			_ = d.Set("console_login", false)
		}
	} else {
		_ = d.Set("console_login", false)
	}

	//tags
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	tags, err := tagService.DescribeResourceTags(ctx, "cam", "role", "", *instance.RoleId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCamRoleByNameUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_by_name.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	d.Partial(true)

	roleName := d.Id()

	description := ""
	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}
		mDescRequest := cam.NewUpdateRoleDescriptionRequest()
		mDescRequest.Description = &description
		mDescRequest.RoleName = &roleName
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateRoleDescription(mDescRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, mDescRequest.GetAction(), mDescRequest.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mDescRequest.GetAction(), mDescRequest.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM role description failed, reason:%s\n", logId, err.Error())
			return err
		}

	}
	document := ""
	if d.HasChange("document") {

		document = d.Get("document").(string)
		camService := CamService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		documentErr := camService.PolicyDocumentForceCheck(document)
		if documentErr != nil {
			return documentErr
		}
		mDocRequest := cam.NewUpdateAssumeRolePolicyRequest()
		mDocRequest.PolicyDocument = &document
		mDocRequest.RoleName = &roleName
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdateAssumeRolePolicy(mDocRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, mDocRequest.GetAction(), mDocRequest.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mDocRequest.GetAction(), mDocRequest.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM role document failed, reason:%s\n", logId, err.Error())
			return err
		}

	}

	d.Partial(false)

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		var instance *cam.RoleInfo
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			params := make(map[string]interface{})
			params["name"] = roleName
			camService := CamService{
				client: meta.(*TencentCloudClient).apiV3Conn,
			}
			instances, e := camService.DescribeRolesByFilter(ctx, params)

			if e != nil {
				return retryError(e)
			}
			if len(instances) != 0 {
				instance = instances[0]
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
			return err
		}
		if instance == nil {
			return fmt.Errorf("Instance can not find by name!")
		}
		resourceName := BuildTagResourceName("cam", "role", "", *instance.RoleId)
		err = tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	return resourceTencentCloudCamRoleByNameRead(d, meta)
}

func resourceTencentCloudCamRoleByNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_role_by_name.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	roleName := d.Id()
	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteRoleByName(ctx, roleName)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
