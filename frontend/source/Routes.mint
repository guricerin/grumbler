routes {
  / {
    Application.setPage(Page::Home)
  }

  /signup {
    Application.setPage(Page::SignUp)
  }

  /user/:id {
    sequence {
      Application.setPage(Page::NotFound)
    }
  }

  * {
    Application.setPage(Page::NotFound)
  }
}
