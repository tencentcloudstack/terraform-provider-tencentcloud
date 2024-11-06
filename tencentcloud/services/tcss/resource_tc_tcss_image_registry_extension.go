package tcss

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcssv20201101 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcss/v20201101"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcssImageRegistryReadOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := tccommon.GetLogId(ctx)
	service := TcssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	respData, err := service.DescribeTcssImageRegistryById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tcss_image_registry` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Url != nil {
		_ = d.Set("url", respData.Url)
	}

	if respData.RegistryType != nil {
		_ = d.Set("registry_type", respData.RegistryType)
	}

	if respData.NetType != nil {
		_ = d.Set("net_type", respData.NetType)
	}

	if respData.RegistryVersion != nil {
		_ = d.Set("registry_version", respData.RegistryVersion)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.ConnDetectDetail != nil {
		tmpList := make([]map[string]interface{}, 0, len(respData.ConnDetectDetail))
		for _, item := range respData.ConnDetectDetail {
			dMap := map[string]interface{}{}
			if item.Quuid != nil {
				dMap["quuid"] = item.Quuid
			}

			if item.Uuid != nil {
				dMap["uuid"] = item.Uuid
			}

			tmpList = append(tmpList, dMap)
		}

		_ = d.Set("conn_detect_config", tmpList)
	}

	if respData.SyncStatus != nil {
		_ = d.Set("sync_status", respData.SyncStatus)
	}

	return nil
}

func resourceTencentCloudTcssImageRegistryUpdatePreRequest0(ctx context.Context, req *tcssv20201101.UpdateAssetImageRegistryRegistryDetailRequest) *resource.RetryError {
	d := tccommon.ResourceDataFromContext(ctx)
	if v, ok := d.GetOk("name"); ok {
		req.Name = helper.String(v.(string))
	}

	return nil
}
