allowlist = ["testdb", "production"]

table "users" {
  pk = "id"

  column "username" {
    gen "faker" {
      provider = "username"
    }
  }

  column "full_name" {
    gen "faker" {
      provider = "full_name"
    }
  }

  column "email" {
    gen "faker" {
      provider = "email"
    }
  }

  column "company" {
    gen "faker" {
      provider = "company_name"
    }
  }

  column "phone" {
    gen "faker" {
      provider = "phone_number"
    }
  }

  column "ip_address" {
    gen "faker" {
      provider = "ipv4_address"
    }
  }

  column "bio" {
    gen "faker" {
      provider = "short_text"
    }
  }

  column "status" {
    gen "constant" {
      value = "masked"
    }
  }
}
