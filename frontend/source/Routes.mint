routes {
  / {
    Application.setPage(Page::Home)
  }

  /signup {
    Application.setPage(Page::SignUp)
  }

  /signin {
    Application.setPage(Page::SignIn)
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
