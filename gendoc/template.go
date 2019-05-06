package main

const (
	docTPL = `---
layout: "{{.cloud_mark}}"
page_title: "{{.cloud_title}}: {{.name}}"
sidebar_current: "docs-{{.cloud_mark}}-{{.dtype}}-{{.resource}}"
description: |-
  {{.description_short}}
---

# {{.name}}

{{.description}}

## Example Usage

{{.example}}

## Argument Reference

The following arguments are supported:

{{.arguments}}
{{if ne .attributes ""}}
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

{{.attributes}}
{{end}}
`
)
