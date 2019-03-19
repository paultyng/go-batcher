workflow "golang" {
  on = "push"
  resolves = ["docker://golang:1.12"]
}

action "docker://golang:1.12" {
  uses = "docker://golang:1.12"
  runs = "go"
  args = "test -v -cover -race ./..."
}
