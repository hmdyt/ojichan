terraform {
  required_providers {
    sakuracloud = {
      source  = "sacloud/sakuracloud"
      version = "2.26.0"
    }
  }
}

provider "sakuracloud" {}

resource "sakuracloud_apprun_application" "ojichan" {
  name            = "ojichan"
  timeout_seconds = 60
  port            = 8080
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "ojichan"
    max_cpu    = "1"
    max_memory = "512Mi"
    deploy_source {
      container_registry {
        image    = "hmdyt.sakuracr.jp/ojichan:latest"
        server   = "hmdyt.sakuracr.jp"
        username = var.SAKURACLOUD_APPRUN_USER
        password = var.SAKURACLOUD_APPRUN_PASSWORD
      }
    }
    env {
      key   = "OJICHAN_DISCORD_TOKEN"
      value = var.OJICHAN_DISCORD_TOKEN
    }
    env {
      key   = "OJICHAN_DISCORD_CHANNEL_ID"
      value = var.OJICHAN_DISCORD_CHANNEL_ID
    }
    env {
      key   = "OJICHAN_EMOJI_NUM"
      value = var.OJICHAN_EMOJI_NUM
    }
    env {
      key   = "OJICHAN_UNCTUATION_LEVEL"
      value = var.OJICHAN_UNCTUATION_LEVEL
    }
    env {
      key   = "OJICHAN_MYSELF_URL"
      value = "https://app-64a2d214-23c2-4010-8041-a986b4e4e27f.ingress.apprun.sakura.ne.jp"
    }
    probe {
      http_get {
        path = "/"
        port = 8080
      }
    }
  }
  traffics {
    version_index = 0
    percent       = 100
  }
}

output "OJICHAN_DISCORD_CHANNEL_ID" {
  value = var.OJICHAN_DISCORD_CHANNEL_ID
}

output "OJICHAN_EMOJI_NUM" {
  value = var.OJICHAN_EMOJI_NUM
}

output "OJICHAN_UNCTUATION_LEVEL" {
  value = var.OJICHAN_UNCTUATION_LEVEL
}

data "sakuracloud_apprun_application" "ojichan" {
  name = sakuracloud_apprun_application.ojichan.name
}

output "public_link" {
  value = data.sakuracloud_apprun_application.ojichan.public_url
}
