allowlist = ["testdb", "production"]

table "users" {
  pk = "id"

  column "full_name" {
    gen "faker" {
      provider = "first_name"
    }
  }

  column "email" {
    gen "faker" {
      provider = "email"
    }
  }

  column "status" {
    gen "constant" {
      value = "active"
    }
  }
}
