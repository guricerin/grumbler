enum Page {
  Initial
  Home
  NotFound
}

enum UserStatus {
  LoggedOut
  LoggedIn
}

store Application {
  state page : Page = Page::Initial
  state userStatus : UserStatus = UserStatus::LoggedOut

  fun initializeWithPage (page : Page) : Promise(Never, Void) {
    sequence {
      setPage(page)
    }
  }

  fun setPage (page : Page) : Promise(Never, Void) {
    next { page = page }
  }
}
