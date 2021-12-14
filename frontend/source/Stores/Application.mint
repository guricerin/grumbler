enum Page {
  Initial
  Home
  SignUp
  Timeline
  NotFound
}

enum UserStatus {
  Guest
  SignIn(User)
}

store Application {
  state isNavMenuActive : Bool = false
  state page : Page = Page::Initial
  state userStatus : UserStatus = UserStatus::Guest

  fun toggleMenu : Promise(Never, Void) {
    next { isNavMenuActive = !isNavMenuActive }
  }

  fun resetMenu : Promise(Never, Void) {
    next { isNavMenuActive = false }
  }

  fun initializeWithPage (page : Page) : Promise(Never, Void) {
    sequence {
      setPage(page)
    }
  }

  fun setPage (page : Page) : Promise(Never, Void) {
    next { page = page }
  }

  fun signin (user : User) : Promise(Never, Void) {
    sequence {
      next { userStatus = UserStatus::SignIn(user) }

      Window.navigate("/user/#{user.id}/timeline")
    }
  }
}
