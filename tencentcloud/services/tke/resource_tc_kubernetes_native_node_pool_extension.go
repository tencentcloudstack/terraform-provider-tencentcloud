package tke

import (
	"context"
	"encoding/base64"

	v20220501 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesNativeNodePoolReadPostHandleResponse0(ctx context.Context, resp *v20220501.NodePool) error {
	d := tccommon.ResourceDataFromContext(ctx)
	respData := resp

	if respData.Native != nil {
		nativeMap := d.Get("native").([]interface{})
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

			nativeMap[0].(map[string]interface{})["lifecycle"] = []interface{}{lifecycleMap}
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
