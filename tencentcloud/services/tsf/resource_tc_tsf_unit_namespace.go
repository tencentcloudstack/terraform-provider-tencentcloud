package tsf

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfUnitNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfUnitNamespaceCreate,
		Read:   resourceTencentCloudTsfUnitNamespaceRead,
		Delete: resourceTencentCloudTsfUnitNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "gateway instance Id.",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "namespace id.",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "namespace name.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time. Note: This field may return null, indicating that no valid value was found.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time. Note: This field may return null, indicating that no valid value was found.",
			},
		},
	}
}

func resourceTencentCloudTsfUnitNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_unit_namespace.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request           = tsf.NewCreateUnitNamespacesRequest()
		response          = tsf.NewCreateUnitNamespacesResponse()
		gatewayInstanceId string
		namespaceId       string
	)
	if v, ok := d.GetOk("gateway_instance_id"); ok {
		gatewayInstanceId = v.(string)
		request.GatewayInstanceId = helper.String(v.(string))
	}

	unitNamespace := tsf.UnitNamespace{}
	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		unitNamespace.NamespaceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("namespace_name"); ok {
		unitNamespace.NamespaceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("gateway_instance_id"); ok {
		unitNamespace.GatewayInstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("created_time"); ok {
		unitNamespace.CreatedTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("updated_time"); ok {
		unitNamespace.UpdatedTime = helper.String(v.(string))
	}
	request.UnitNamespaceList = append(request.UnitNamespaceList, &unitNamespace)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().CreateUnitNamespaces(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf unitNamespace failed, reason:%+v", logId, err)
		return err
	}

	if !*response.Response.Result {
		return fmt.Errorf("[CRITAL]%s create tsf unitNamespace failed", logId)
	}

	d.SetId(gatewayInstanceId + tccommon.FILED_SP + namespaceId)

	return resourceTencentCloudTsfUnitNamespaceRead(d, meta)
}

func resourceTencentCloudTsfUnitNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_unit_namespace.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayInstanceId := idSplit[0]
	namespaceId := idSplit[1]

	unitNamespace, err := service.DescribeTsfUnitNamespaceById(ctx, gatewayInstanceId, namespaceId)
	if err != nil {
		return err
	}

	if unitNamespace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfUnitNamespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_instance_id", gatewayInstanceId)
	_ = d.Set("namespace_id", namespaceId)

	if unitNamespace.NamespaceName != nil {
		_ = d.Set("namespace_name", unitNamespace.NamespaceName)
	}

	if unitNamespace.CreatedTime != nil {
		_ = d.Set("created_time", unitNamespace.CreatedTime)
	}

	if unitNamespace.UpdatedTime != nil {
		_ = d.Set("updated_time", unitNamespace.UpdatedTime)
	}

	return nil
}

func resourceTencentCloudTsfUnitNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_unit_namespace.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayInstanceId := idSplit[0]
	namespaceId := idSplit[1]

	unitNamespace, err := service.DescribeTsfUnitNamespaceById(ctx, gatewayInstanceId, namespaceId)
	if err != nil {
		return err
	}

	if err := service.DeleteTsfUnitNamespaceById(ctx, gatewayInstanceId, *unitNamespace.Id); err != nil {
		return err
	}

	return nil
}
