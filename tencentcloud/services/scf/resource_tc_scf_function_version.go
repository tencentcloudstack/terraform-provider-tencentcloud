package scf

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfFunctionVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionVersionCreate,
		Read:   resourceTencentCloudScfFunctionVersionRead,
		Delete: resourceTencentCloudScfFunctionVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Name of the released function.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Function description.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Function namespace.",
			},

			"function_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Version of the released function.",
			},
		},
	}
}

func resourceTencentCloudScfFunctionVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_version.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = scf.NewPublishVersionRequest()
		response        = scf.NewPublishVersionResponse()
		functionName    string
		namespace       string
		functionVersion string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().PublishVersion(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf FunctionVersion failed, reason:%+v", logId, err)
		return err
	}

	functionVersion = *response.Response.FunctionVersion
	d.SetId(functionName + tccommon.FILED_SP + namespace + tccommon.FILED_SP + functionVersion)

	// wait ready
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	err = waitScfFunctionReady(ctx, functionName, namespace, client.UseScfClient())
	if err != nil {
		return err
	}

	return resourceTencentCloudScfFunctionVersionRead(d, meta)
}

func resourceTencentCloudScfFunctionVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	functionVersion := idSplit[2]

	version, err := service.DescribeScfFunctionVersionById(ctx, functionName, namespace, functionVersion)
	if err != nil {
		return err
	}

	if version == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if version.Response.FunctionName != nil {
		_ = d.Set("function_name", version.Response.FunctionName)
	}

	if version.Response.Description != nil {
		_ = d.Set("description", version.Response.Description)
	}

	if version.Response.Namespace != nil {
		_ = d.Set("namespace", version.Response.Namespace)
	}

	if version.Response.FunctionVersion != nil {
		_ = d.Set("function_version", version.Response.FunctionVersion)
	}

	return nil
}

func resourceTencentCloudScfFunctionVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_version.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	functionVersion := idSplit[2]

	if err := service.DeleteScfFunctionVersionById(ctx, functionName, namespace, functionVersion); err != nil {
		return err
	}

	return nil
}
