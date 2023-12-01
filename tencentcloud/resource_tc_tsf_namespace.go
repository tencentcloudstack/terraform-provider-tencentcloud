/*
Provides a resource to create a tsf namespace

Example Usage

```hcl
resource "tencentcloud_tsf_namespace" "namespace" {
  namespace_name = "namespace-name"
  # cluster_id = "cls-xxxx"
  namespace_desc = "namespace desc"
  # namespace_resource_type = ""
  namespace_type = "DEF"
  # namespace_id = ""
  is_ha_enable = "0"
  # program_id = ""
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfNamespaceCreate,
		Read:   resourceTencentCloudTsfNamespaceRead,
		Update: resourceTencentCloudTsfNamespaceUpdate,
		Delete: resourceTencentCloudTsfNamespaceDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"cluster_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster ID.",
			},

			"namespace_desc": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace description.",
			},

			"namespace_resource_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace resource type (default is DEF).",
			},

			"namespace_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether it is a global namespace (the default is DEF, which means a common namespace; GLOBAL means a global namespace).",
			},

			"namespace_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID.",
			},

			"is_ha_enable": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "whether to enable high availability.",
			},

			"program_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the dataset to be bound.",
			},

			"kube_inject_enable": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "KubeInjectEnable value.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"namespace_code": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace encoding.",
			},

			"is_default": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "default namespace.",
			},

			"namespace_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "namespace status.",
			},

			"delete_flag": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Delete ID.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "update time.",
			},
		},
	}
}

func resourceTencentCloudTsfNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tsf.NewCreateNamespaceRequest()
		response    = tsf.NewCreateNamespaceResponse()
		namespaceId string
	)
	if v, ok := d.GetOk("namespace_name"); ok {
		request.NamespaceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_desc"); ok {
		request.NamespaceDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_resource_type"); ok {
		request.NamespaceResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_type"); ok {
		request.NamespaceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_ha_enable"); ok {
		request.IsHaEnable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id"); ok {
		request.ProgramId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf namespace failed, reason:%+v", logId, err)
		return err
	}

	namespaceId = *response.Response.Result
	d.SetId(namespaceId)

	return resourceTencentCloudTsfNamespaceRead(d, meta)
}

func resourceTencentCloudTsfNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	namespaceId := d.Id()

	namespace, err := service.DescribeTsfNamespaceById(ctx, namespaceId)
	if err != nil {
		return err
	}

	if namespace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfNamespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if namespace.NamespaceName != nil {
		_ = d.Set("namespace_name", namespace.NamespaceName)
	}

	if namespace.ClusterId != nil {
		_ = d.Set("cluster_id", namespace.ClusterId)
	}

	if namespace.NamespaceDesc != nil {
		_ = d.Set("namespace_desc", namespace.NamespaceDesc)
	}

	if namespace.NamespaceResourceType != nil {
		_ = d.Set("namespace_resource_type", namespace.NamespaceResourceType)
	}

	if namespace.NamespaceType != nil {
		_ = d.Set("namespace_type", namespace.NamespaceType)
	}

	if namespace.NamespaceId != nil {
		_ = d.Set("namespace_id", namespace.NamespaceId)
	}

	if namespace.IsHaEnable != nil {
		_ = d.Set("is_ha_enable", namespace.IsHaEnable)
	}

	// if namespace.ProgramId != nil {
	// 	_ = d.Set("program_id", namespace.ProgramId)
	// }

	if namespace.KubeInjectEnable != nil {
		_ = d.Set("kube_inject_enable", namespace.KubeInjectEnable)
	}

	// if namespace.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", namespace.ProgramIdList)
	// }

	if namespace.NamespaceCode != nil {
		_ = d.Set("namespace_code", namespace.NamespaceCode)
	}

	if namespace.IsDefault != nil {
		_ = d.Set("is_default", namespace.IsDefault)
	}

	if namespace.NamespaceStatus != nil {
		_ = d.Set("namespace_status", namespace.NamespaceStatus)
	}

	if namespace.DeleteFlag != nil {
		_ = d.Set("delete_flag", namespace.DeleteFlag)
	}

	if namespace.CreateTime != nil {
		_ = d.Set("create_time", namespace.CreateTime)
	}

	if namespace.UpdateTime != nil {
		_ = d.Set("update_time", namespace.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTsfNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyNamespaceRequest()

	namespaceId := d.Id()

	request.NamespaceId = &namespaceId

	immutableArgs := []string{"namespace_name", "cluster_id", "namespace_resource_type", "namespace_type", "program_id", "kube_inject_enable", "program_id_list", "namespace_code", "is_default", "namespace_status", "delete_flag", "create_time", "update_time", "cluster_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("namespace_desc") {
		if v, ok := d.GetOk("namespace_desc"); ok {
			request.NamespaceDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_id") {
		if v, ok := d.GetOk("namespace_id"); ok {
			request.NamespaceId = helper.String(v.(string))
		}
	}

	if d.HasChange("is_ha_enable") {
		if v, ok := d.GetOk("is_ha_enable"); ok {
			request.IsHaEnable = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf namespace failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfNamespaceRead(d, meta)
}

func resourceTencentCloudTsfNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	namespaceId := d.Id()

	if err := service.DeleteTsfNamespaceById(ctx, namespaceId); err != nil {
		return err
	}

	return nil
}
