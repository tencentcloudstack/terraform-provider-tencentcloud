package tke

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	v20220501 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesNativeNodePoolReadPostHandleResponse0(ctx context.Context, resp *v20220501.NodePool) error {
	d := tccommon.ResourceDataFromContext(ctx)
	respData := resp

	if respData.Native != nil {
		nativeMap := d.Get("native").([]interface{})[0].(map[string]interface{})
		lifecycleMap := map[string]interface{}{}
		if respData.Native.Lifecycle != nil {
			if respData.Native.Lifecycle.PreInit != nil {
				lifecycleMap["pre_init"] = base64.StdEncoding.EncodeToString([]byte(*respData.Native.Lifecycle.PreInit))
				//lifecycleMap["pre_init"] = respData.Native.Lifecycle.PreInit
			}

			if respData.Native.Lifecycle.PostInit != nil {
				lifecycleMap["post_init"] = base64.StdEncoding.EncodeToString([]byte(*respData.Native.Lifecycle.PostInit))
				//lifecycleMap["post_init"] = respData.Native.Lifecycle.PostInit
			}
			nativeMap["lifecycle"] = []interface{}{lifecycleMap}
			_ = d.Set("native", []interface{}{nativeMap})
		}
	}

	annotationsList := make([]map[string]interface{}, 0, len(respData.Annotations))
	if respData.Annotations != nil {
		for _, annotations := range respData.Annotations {
			annotationsMap := map[string]interface{}{}

			if annotations.Name != nil && tkeNativeNodePoolAnnotationsMap[*annotations.Name] != "" {
				continue
			}

			if annotations.Name != nil {
				annotationsMap["name"] = annotations.Name
			}

			if annotations.Value != nil {
				annotationsMap["value"] = annotations.Value
			}

			annotationsList = append(annotationsList, annotationsMap)
		}

		_ = d.Set("annotations", annotationsList)
	}

	return nil
}

func resourceTencentCloudKubernetesNativeNodePoolDeletePostHandleResponse0(ctx context.Context, resp *v20220501.DeleteNodePoolResponse) error {
	// wait for delete ok
	logId := tccommon.GetLogId(tccommon.ContextNil)
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	var (
		request = v20220501.NewDeleteNodePoolRequest()
	)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		respData, errRet := service.DescribeKubernetesNativeNodePoolById(ctx, clusterId, nodePoolId)
		if errRet != nil {
			errCode := errRet.(*sdkErrors.TencentCloudSDKError).Code
			if strings.Contains(errCode, "InternalError") {
				return nil
			}
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if respData != nil && *respData.LifeState == "Deleting" {
			log.Printf("[DEBUG]%s api[%s] native node pool %s still alive and status is %s", logId, request.GetAction(), nodePoolId, *respData.LifeState)
			return resource.RetryableError(fmt.Errorf("native node pool %s still alive and status is %s", nodePoolId, *respData.LifeState))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
