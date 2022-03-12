component Components.Header {
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
          <Components.NavbarItem
            route="/signin"
            title="サインイン"/>,
          <Components.NavbarItem
            route="/signup"
            title="新規登録"/>
        ]

      UserStatus::SignIn(user) =>
        [
          <Components.NavbarItem
            route="/timeline"
            title="タイムライン"/>,
          <Components.NavbarItem
            route="/post-grumble"
            title="ぼやく"/>,
          <Components.NavbarItem
            route="/search"
            title="検索"/>,
          <Components.NavbarItem
            route="/signout"
            title="サインアウト"/>,
          <Components.NavbarItem
            route="/unsubscribe"
            title="退会"/>
        ]
    }
  }

  fun navbarUser : Html {
    case (userStatus) {
      UserStatus::Guest =>
        <Components.NavbarUser
          route="/"
          title="Guest"
          icon="fas fa-user"/>

      UserStatus::SignIn(user) =>
        <Components.NavbarUser
          route="/user/#{user.id}"
          title="#{user.name}@#{user.id}"
          icon="fas fa-user"/>
    }
  }
}

component Components.NavbarItem {
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

component Components.NavbarUser {
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
