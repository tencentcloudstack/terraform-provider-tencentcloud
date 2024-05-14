package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCvmImportImageOsReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return map[string]interface{}{
		"import_image_os_list_supported": d.Get("import_image_os_list_supported"),
		"import_image_os_version_set":    d.Get("import_image_os_version_set"),
	}
}
