routes {
  / {
    Application.setPage(Page::Home)
  }

  /signup {
    Application.setPage(Page::SignUp)
  }

  /user/:id/timeline {
    sequence {
      Application.setPage(Page::Timeline)
    }
  }

  * {
    Application.setPage(Page::NotFound)
  }
}
