component Header {
  connect Application exposing { navMenuStatus }
  property userStatus : UserStatus

  style head {
    box-shadow: 2px 2px 8px rgba(0,0,0,.06),0px .5px 1px rgba(0,0,0,.05);
  }

  fun render : Html {
    <div>
      <nav::head class="navbar is-fixed-top ">
        <div class="navbar-brand">
          <a
            class="navbar-item"
            href="/"
            onClick={Application.resetMenu}>

            <img
              src={@asset(../../assets/logo.svg)}
              width="64"
              height="64"
              alt="grumbler"/>

            <h1>"Grumbler"</h1>

          </a>

          <{ navbarBurger() }>
        </div>

        <{ navbarMenu() }>
      </nav>
    </div>
  }

  get getNavMenuStatus : String {
    case (navMenuStatus) {
      NavMenuStatus::Active => "is-active"
      NavMenuStatus::Reset => "burger"
    }
  }

  fun navbarBurger : Html {
    <div
      class="navbar-burger burger #{getNavMenuStatus}"
      data-target="navMenu"
      onClick={Application.toggleMenu}>

      <span/>
      <span/>
      <span/>

    </div>
  }

  fun navbarMenu : Html {
    <div
      id="navMenu"
      class="navbar-menu #{getNavMenuStatus}">

      <div class="navbar-start">
        <{ navbarItems() }>
      </div>

      <div class="navbar-end">
        <{ navbarUser() }>
      </div>

    </div>
  }

  fun navbarItems : Array(Html) {
    case (userStatus) {
      UserStatus::LoggedOut =>
        [
          <NavbarItem
            route="/sign-in"
            title="Sign In"/>,
          <NavbarItem
            route="/sign-up"
            title="Sign Up"/>
        ]

      UserStatus::LoggedIn =>
        [
          <NavbarItem
            route="/sign-out"
            title="Sign Out"/>
        ]
    }
  }

  fun navbarUser : Html {
    case (userStatus) {
      UserStatus::LoggedOut =>
        <NavbarUser
          title="Guest"
          icon="fas fa-user"/>

      UserStatus::LoggedIn =>
        <NavbarUser
          route="/"
          title="Login User"/>
    }
  }
}

component NavbarItem {
  property route : String = ""
  property title : String = ""
  fun render : Html {
    <a
      class="navbar-item"
      href={route}
      onClick={Application.toggleMenu}>
      <span>"#{title}"</span>

    </a>
  }
}

component NavbarUser {
  property route : String = ""
  property title : String = ""
  property icon : String = ""

  style icon {
    margin-right: 10px;
  }
  fun render : Html {
    <a
      class="navbar-item"
      href={route}
      onClick={Application.toggleMenu}>
      <i::icon class="#{icon}"></i>
      <span>"#{title}"</span>

    </a>
  }
}
