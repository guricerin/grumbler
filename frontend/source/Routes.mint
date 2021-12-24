routes {
  / {
    sequence {
      Application.dbgUser()
      Application.initializeWithPage(Page::Home)
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

  /user/:id/grumble (id : String) {
    sequence {
      Application.setPageWithAuthorization(id, Page::Grumble)
    }
  }

  /user/:id/signout (id : String) {
    Application.setPageWithAuthorization(id, Page::SignOut)
  }

  /user/:id/unsubscribe (id : String) {
    Application.setPageWithAuthorization(id, Page::Unsubscribe)
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
