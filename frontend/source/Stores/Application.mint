enum NavMenuStatus {
  Active
  Reset
}

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
  state navMenuStatus : NavMenuStatus = NavMenuStatus::Reset
  state page : Page = Page::Initial
  state userStatus : UserStatus = UserStatus::LoggedOut

  fun toggleMenu : Promise(Never, Void) {
    case (navMenuStatus) {
      NavMenuStatus::Active => next { navMenuStatus = NavMenuStatus::Reset }
      NavMenuStatus::Reset => next { navMenuStatus = NavMenuStatus::Active }
    }
  }

  fun resetMenu : Promise(Never, Void) {
    next { navMenuStatus = NavMenuStatus::Reset }
  }

  fun initializeWithPage (page : Page) : Promise(Never, Void) {
    sequence {
      setPage(page)
    }
  }

  fun setPage (page : Page) : Promise(Never, Void) {
    next { page = page }
  }
}
