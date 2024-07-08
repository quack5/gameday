terraform {
  cloud {
    organization = "<your-org>"

    workspaces {
      name = "<your-workspace-name>"
    }
  }
}