workflow "New workflow" {
  on = "push"
  resolves = ["new-action"]
}

action "Setup Go for use with actions" {
  uses = "actions/setup-go@419ae75c254126fa6ae3e3ef573ce224a919b8fe"
}

action "new-action" {
  uses = "owner/repo/path@ref"
  needs = ["Setup Go for use with actions"]
}
