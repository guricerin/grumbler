component Header {
  connect Application exposing { isNavMenuActive }
  property userStatus : UserStatus

  style head {
    box-shadow: 2px 2px 8px rgba(0,0,0,.06),0px .5px 1px rgba(0,0,0,.05);
  }

  fun render : Html {
    <div>
      <nav::head class="navbar is-fixed-top">
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
    if (isNavMenuActive) {
      "is-active"
    } else {
      "burger"
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
      UserStatus::Guest =>
        [
          <NavbarItem
            route="/signin"
            title="サインイン"/>,
          <NavbarItem
            route="/signup"
            title="新規登録"/>
        ]

      UserStatus::SignIn(user) =>
        [
          <NavbarItem
            route="/user/#{user.id}/post-grumble"
            title="ぼやく"/>,
          <NavbarItem
            route="/search"
            title="検索"/>,
          <NavbarItem
            route="/signout"
            title="サインアウト"/>,
          <NavbarItem
            route="/unsubscribe"
            title="退会"/>
        ]
    }
  }

  fun navbarUser : Html {
    case (userStatus) {
      UserStatus::Guest =>
        <NavbarUser
          route="/"
          title="Guest"
          icon="fas fa-user"/>

      UserStatus::SignIn(user) =>
        <NavbarUser
          route="/user/#{user.id}"
          title="#{user.name}@#{user.id}"
          icon="fas fa-user"/>
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
      onClick={Application.resetMenu}>

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
      onClick={Application.resetMenu}>

      <i::icon class="#{icon}"/>
      <span>"#{title}"</span>

    </a>
  }
}
