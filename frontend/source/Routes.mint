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

  /search {
    Application.setPageWithAuthentication(Page::Search)
  }

  /search?q=:q&k=:k (q : String, k : String) {
    sequence {
      Stores.Search.setSearchWord(q)
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
