package cvm

import (
	"context"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func dataSourceTencentCloudCvmImportImageOsReadOutputContent(ctx context.Context) interface{} {
	d := tccommon.ResourceDataFromContext(ctx)

	return map[string]interface{}{
		"import_image_os_list_supported": d.Get("import_image_os_list_supported"),
		"import_image_os_version_set":    d.Get("import_image_os_version_set"),
	}
}

func dataSourceTencentCloudCvmImportImageOsReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeImportImageOsResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)

	d.SetId(helper.BuildToken())

	return nil
}
