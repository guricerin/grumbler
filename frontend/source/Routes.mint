routes {
  / {
    Application.setPage(Page::Home)
  }

  /signup {
    Application.setPage(Page::SignUp)
  }

  /signin {
    Application.setPage(Page::SignIn)
  }

  /user/:id/signout (id : String) {
    sequence {
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
    Application.setPage(Page::Error(404))
  }
}
