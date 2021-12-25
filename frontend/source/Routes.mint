routes {
  / {
    sequence {
      Application.dbgUser()
      Application.initializeWithPage(Page::Home)
    }
  }

  /signup {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::SignUp)

        UserStatus::SignIn => Window.navigate("/")
      }
    }
  }

  /signin {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::SignIn)

        UserStatus::SignIn => Window.navigate("/")
      }
    }
  }

  /search?q=:query&k=:kind (query : String, kind : String) {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(401))

        UserStatus::SignIn =>
          sequence {
            Stores.Search.search(query, kind)
            Application.setPage(Page::Search)
          }
      }
    }
  }

  /search {
    sequence {
      Stores.Search.resetApiStatus()
      Application.setPageWithAuthentication(Page::Search)
    }
  }

  /grumble {
    sequence {
      Application.setPageWithAuthentication(Page::Grumble)
    }
  }

  /signout {
    Application.setPageWithAuthentication(Page::SignOut)
  }

  /unsubscribe {
    Application.setPageWithAuthentication(Page::Unsubscribe)
  }

  /user/:id/timeline (id : String) {
    Application.setPageWithAuthorization(id, Page::Timeline)
  }

  * {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::Error(404))
    }
  }
}
