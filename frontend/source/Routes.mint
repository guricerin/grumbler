routes {
  / {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::Home)
    }
  }

  /signup {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::SignUp)
    }
  }

  /signin {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::SignIn)
    }
  }

  /user/:id/signout (id : String) {
    sequence {
      Application.dbgUser()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          if (user.id == id) {
            Application.setPage(Page::SignOut)
          } else {
            Application.setPage(Page::Error(403))
          }
      }
    }
  }

  /user/:id/timeline (id : String) {
    sequence {
      Application.dbgUser()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          if (user.id == id) {
            Application.setPage(Page::Timeline)
          } else {
            Application.setPage(Page::Error(403))
          }
      }
    }
  }

  * {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::Error(404))
    }
  }
}
