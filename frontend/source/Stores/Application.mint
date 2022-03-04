enum Page {
  Initial
  Home
  SignUp
  SignIn
  SignOut
  Unsubscribe
  Search
  User
  Timeline
  PostGrumble
  UserGrumbles
  Error(Number)
}

enum UserStatus {
  Guest
  SignIn(User)
}

store Application {
  state isNavMenuActive : Bool = false
  state page : Page = Page::Initial
  state userStatus : UserStatus = UserStatus::Guest

  fun dbgUser : String {
    case (userStatus) {
      UserStatus::Guest => Debug.log("Guest")
      UserStatus::SignIn(user) => Debug.log("signin -  id: #{user.id}, name: #{user.name}")
    }
  }

  fun toggleMenu : Promise(Never, Void) {
    next { isNavMenuActive = !isNavMenuActive }
  }

  fun resetMenu : Promise(Never, Void) {
    next { isNavMenuActive = false }
  }

  fun signinCheck : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/signin-check")
        |> Api.send(User.decodes)

      case (status) {
        Api.Status::Ok(user) => next { userStatus = UserStatus::SignIn(user) }
        => next { userStatus = UserStatus::Guest }
      }
    }
  }

  fun setPage (page : Page) : Promise(Never, Void) {
    next { page = page }
  }

  fun initializeWithPage (page : Page) : Promise(Never, Void) {
    sequence {
      signinCheck()
      setPage(page)
    }
  }

  fun setPageWithAuthentication (page : Page) : Promise(Never, Void) {
    sequence {
      signinCheck()

      case (userStatus) {
        UserStatus::Guest => setPage(Page::Error(401))

        UserStatus::SignIn(user) =>
          setPage(page)
      }
    }
  }

  fun setPageWithAuthorization (userId : String, page : Page) : Promise(Never, Void) {
    sequence {
      signinCheck()

      case (userStatus) {
        UserStatus::Guest => setPage(Page::Error(403))

        UserStatus::SignIn(user) =>
          if (user.id == userId) {
            setPage(page)
          } else {
            setPage(Page::Error(403))
          }
      }
    }
  }

  fun signin (user : User) : Promise(Never, Void) {
    sequence {
      next { userStatus = UserStatus::SignIn(user) }

      Window.navigate("/user/#{user.id}/timeline")
    }
  }

  fun signout : Promise(Never, Void) {
    sequence {
      next { userStatus = UserStatus::Guest }

      Window.navigate("/")
    }
  }
}
