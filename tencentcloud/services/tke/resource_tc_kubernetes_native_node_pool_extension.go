package tke

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesNativeNodePoolReadPostHandleResponse0(ctx context.Context, resp *v20220501.NodePool) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	respData := resp
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
}
