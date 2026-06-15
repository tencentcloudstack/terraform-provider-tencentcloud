//go:build tools
// +build tools

package main

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/katbyte/terrafmt"
	_ "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	_ "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"
)
