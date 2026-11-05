output "tool_id" {
  value       = tencentcloud_ags_sandbox_tool.this.id
  description = "Use this ID with the runtime SDK; it is not a connection token."
}

output "tool_name" {
  value = tencentcloud_ags_sandbox_tool.this.tool_name
}

output "status" {
  value = tencentcloud_ags_sandbox_tool.this.status
}

output "image_digest" {
  value = tencentcloud_ags_sandbox_tool.this.image_digest
}
