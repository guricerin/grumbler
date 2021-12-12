routes {
  / {
    Application.setPage(Page::Home)
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
