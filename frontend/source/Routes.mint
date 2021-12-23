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
      Stores.Search.search(query, kind)
      Application.setPageWithAuthentication(Page::Search)
    }
  }

  /search {
    sequence {
      Stores.Search.resetApiStatus()
      Application.setPageWithAuthentication(Page::Search)
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
