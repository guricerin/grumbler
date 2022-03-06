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

  /user/:id/post-grumble (id : String) {
    sequence {
      Application.setPageWithAuthentication(Page::PostGrumble)
    }
  }

  /signout {
    Application.setPageWithAuthentication(Page::SignOut)
  }

  /unsubscribe {
    Application.setPageWithAuthentication(Page::Unsubscribe)
  }

  /user/:id/timeline (id : String) {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          if (user.id == id) {
            parallel {
              Application.setPage(Page::Timeline)
              Stores.Timeline.getTimeline(id)
            }
          } else {
            Application.setPage(Page::Error(403))
          }
      }
    }
  }

  /user/:id/grumbles (id : String) {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          parallel {
            Stores.PageUser.getGrumbles(id)
            Application.setPage(Page::UserGrumbles)
          }
      }
    }
  }

  /user/:id (id : String) {
    parallel {
      Stores.PageUser.getUserDetail(id)
      Application.setPageWithAuthentication(Page::UserDetail)
    }
  }

  * {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::Error(404))
    }
  }
}
