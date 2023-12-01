/*
Provides a resource to create a oceanus resource

Example Usage

```hcl
resource "tencentcloud_oceanus_resource" "example" {
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  resource_type          = 1
  remark                 = "remark."
  name                   = "tf_example"
  resource_config_remark = "config remark."
  folder_id              = "folder-7ctl246z"
  work_space_id          = "space-2idq8wbr"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOceanusResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusResourceCreate,
		Read:   resourceTencentCloudOceanusResourceRead,
		Update: resourceTencentCloudOceanusResourceUpdate,
		Delete: resourceTencentCloudOceanusResourceDelete,

		Schema: map[string]*schema.Schema{
			"resource_loc": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Resource location.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The available storage types for resource location are currently limited to 1:COS.",
						},
						"param": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Json to describe resource location.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource bucket.",
									},
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource path.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource region, if not set, use resource region, note: this field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},
			"resource_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Resource type, only support JAR now, value is 1.",
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource description.",
			},
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource name.",
			},
			"resource_config_remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource version description.",
			},
			"folder_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Folder id.",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace serialId.",
			},
			"resource_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},
			"version": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Resource Version.",
			},
		},
	}
}

func resourceTencentCloudOceanusResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_resource.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = oceanus.NewCreateResourceRequest()
		response   = oceanus.NewCreateResourceResponse()
		resourceId string
		version    string
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "resource_loc"); ok {
		resourceLoc := oceanus.ResourceLoc{}
		if v, ok := dMap["storage_type"]; ok {
			resourceLoc.StorageType = helper.IntInt64(v.(int))
		}

		if paramMap, ok := helper.InterfaceToMap(dMap, "param"); ok {
			resourceLocParam := oceanus.ResourceLocParam{}
			if v, ok := paramMap["bucket"]; ok {
				resourceLocParam.Bucket = helper.String(v.(string))
			}

			if v, ok := paramMap["path"]; ok {
				resourceLocParam.Path = helper.String(v.(string))
			}

			if v, ok := paramMap["region"]; ok {
				resourceLocParam.Region = helper.String(v.(string))
			}

			resourceLoc.Param = &resourceLocParam
		}

		request.ResourceLoc = &resourceLoc
	}

	if v, ok := d.GetOkExists("resource_type"); ok {
		request.ResourceType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_config_remark"); ok {
		request.ResourceConfigRemark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("folder_id"); ok {
		request.FolderId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOceanusClient().CreateResource(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus resource not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus resource failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.ResourceId
	versionInt := *response.Response.Version
	version = strconv.FormatInt(versionInt, 10)
	d.SetId(strings.Join([]string{resourceId, version}, FILED_SP))

	return resourceTencentCloudOceanusResourceRead(d, meta)
}

func resourceTencentCloudOceanusResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_resource.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]
	version := idSplit[1]

	resourceItem, err := service.DescribeOceanusResourceConfigById(ctx, resourceId, version)
	if err != nil {
		return err
	}

	if resourceItem == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusResource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("resource_id", resourceItem.ResourceId)
	_ = d.Set("version", resourceItem.Version)

	if resourceItem.ResourceLoc != nil {
		resourceLocMap := map[string]interface{}{}

		if resourceItem.ResourceLoc.StorageType != nil {
			resourceLocMap["storage_type"] = resourceItem.ResourceLoc.StorageType
		}

		if resourceItem.ResourceLoc.Param != nil {
			paramMap := map[string]interface{}{}

			if resourceItem.ResourceLoc.Param.Bucket != nil {
				paramMap["bucket"] = resourceItem.ResourceLoc.Param.Bucket
			}

			if resourceItem.ResourceLoc.Param.Path != nil {
				paramMap["path"] = resourceItem.ResourceLoc.Param.Path
			}

			if resourceItem.ResourceLoc.Param.Region != nil {
				paramMap["region"] = resourceItem.ResourceLoc.Param.Region
			}

			resourceLocMap["param"] = []interface{}{paramMap}
		}

		_ = d.Set("resource_loc", []interface{}{resourceLocMap})
	}

	if resourceItem.ResourceType != nil {
		_ = d.Set("resource_type", resourceItem.ResourceType)
	}

	return nil
}

func resourceTencentCloudOceanusResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_resource.update")()
	defer inconsistentCheck(d, meta)()

	immutableArgs := []string{"resource_loc", "resource_type", "remark", "name", "resource_config_remark", "folder_id", "work_space_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudOceanusResourceRead(d, meta)
}

func resourceTencentCloudOceanusResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_oceanus_resource.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]

	if err := service.DeleteOceanusResourceById(ctx, resourceId); err != nil {
		return err
	}

	return nil
}
