routes {
  / {
    Application.setPage(Page::Home)
  }

  /signup {
    Application.setPage(Page::SignUp)
  }

  /user/:id {
    sequence {
      Application.setPage(Page::Initial)
    }
  }

  * {
    Application.setPage(Page::NotFound)
  }
}
