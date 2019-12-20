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
{{if ne .import ""}}
## Import

{{.import}}
{{end}}
`
	idxTPL = `
<% wrap_layout :inner do %>
    <% content_for :sidebar do %>
        <div class="docs-sidebar hidden-print affix-top" role="complementary">
            <ul class="nav docs-sidenav">
                <li>
                    <a href="/docs/providers/index.html">All Providers</a>
                </li>
                <li>
                    <a href="/docs/providers/{{.cloud_mark}}/index.html">{{.cloud_title}} Provider</a>
                </li>
                {{range .Products}}
                <li>
                    <a href="#">{{.Name}}</a>
                    <ul class="nav">
                        <li>
                            <a href="#">Data Sources</a>
                            <ul class="nav nav-auto-expand">{{range $Resource := .DataSources}}
                                <li>
                                    <a href="/docs/providers/{{$.cloud_mark}}/d/{{replace $Resource $.cloudPrefix ""}}.html">{{$Resource}}</a>
                                </li>{{end}}
                            </ul>
                        </li>
                        <li>
                            <a href="#">Resources</a>
                            <ul class="nav nav-auto-expand">{{range $Resource := .Resources}}
                                <li>
                                    <a href="/docs/providers/{{$.cloud_mark}}/r/{{replace $Resource $.cloudPrefix ""}}.html">{{$Resource}}</a>
                                </li>{{end}}
                            </ul>
                        </li>
                    </ul>
                </li>{{end}}
            </ul>
        </div>
    <% end %>
    <%= yield %>
<% end %>
`
)
