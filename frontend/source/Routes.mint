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

  /post-grumble {
    sequence {
      Application.setPageWithAuthentication(Page::PostGrumble)
    }
  }

  /reply/:grumblePk (grumblePk : String) {
    sequence {
      Stores.GrumbleDetail.getGrumbleDetail(grumblePk)
      Application.setPageWithAuthentication(Page::Reply)
    }
  }

  /signout {
    Application.setPageWithAuthentication(Page::SignOut)
  }

  /unsubscribe {
    Application.setPageWithAuthentication(Page::Unsubscribe)
  }

  /timeline {
    sequence {
      Application.signinCheck()

      case (Application.userStatus) {
        UserStatus::Guest => Application.setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          sequence {
            Stores.Timeline.getTimeline(user.id)
            Application.setPage(Page::Timeline)
          }
      }
    }
  }

  /settings {
    Application.setPageWithAuthentication(Page::UserSettings)
  }

  /user/:id/follows (id : String) {
    sequence {
      Stores.PageUser.getUserDetail(id)
      Stores.PageUser.setShowKind(UserDetailShowKind::Follows)
      Application.setPageWithAuthentication(Page::UserDetail)
    }
  }

  /user/:id/followers (id : String) {
    sequence {
      Stores.PageUser.getUserDetail(id)
      Stores.PageUser.setShowKind(UserDetailShowKind::Followers)
      Application.setPageWithAuthentication(Page::UserDetail)
    }
  }

  /user/:id/bookmarks (id : String) {
    sequence {
      Stores.PageUser.getUserDetail(id)
      Stores.PageUser.setShowKind(UserDetailShowKind::Bookmarks)
      Application.setPageWithAuthentication(Page::UserDetail)
    }
  }

  /user/:userId/grumble/:grumblePk/:anchor (userId : String, grumblePk : String, anchor : String) {
    sequence {
      Stores.GrumbleDetail.getGrumbleDetail(grumblePk)
      Application.setPageWithAuthentication(Page::GrumbleDetail)
      Window.triggerHashJump()
    }
  }

  /user/:userId/grumble/:grumblePk (userId : String, grumblePk : String) {
    sequence {
      Stores.GrumbleDetail.getGrumbleDetail(grumblePk)
      Application.setPageWithAuthentication(Page::GrumbleDetail)
    }
  }

  /user/:id (id : String) {
    sequence {
      Stores.PageUser.getUserDetail(id)
      Stores.PageUser.setShowKind(UserDetailShowKind::Grumbles)
      Application.setPageWithAuthentication(Page::UserDetail)
    }
  }

  * {
    Application.setPage(Page::Error(404))
  }
}
