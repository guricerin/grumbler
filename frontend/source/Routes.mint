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
    Application.setPageWithAuthentication(id, Page::SignOut)
  }

  /user/:id/unsubscribe (id : String) {
    Application.setPageWithAuthentication(id, Page::Unsubscribe)
  }

  /user/:id/timeline (id : String) {
    Application.setPageWithAuthentication(id, Page::Timeline)
  }

  * {
    sequence {
      Application.dbgUser()
      Application.setPage(Page::Error(404))
    }
  }
}
