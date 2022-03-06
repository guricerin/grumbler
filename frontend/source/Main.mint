component Main {
  connect Application exposing { page, userStatus }

  style app {
    font-family: Open Sans;
  }

  fun render : Html {
    <div::app>
      <Header userStatus={userStatus}/>

      <Content
        page={page}
        userStatus={userStatus}/>
    </div>
  }
}

component Content {
  property page : Page
  property userStatus : UserStatus

  style section {
    margin-top: 50px;
  }

  fun render : Html {
    <div class="container sf-site-all">
      <section::section class="section sf-site-content">
        case (page) {
          Page::Initial => Html.empty()
          Page::Home => <Pages.Home/>
          Page::SignUp => <Pages.SignUp/>
          Page::SignIn => <Pages.SignIn/>
          Page::Search => <Pages.Search/>
          Page::SignOut => <Pages.SignOut/>
          Page::Unsubscribe => <Pages.Unsubscribe/>
          Page::UserDetail => <Pages.UserDetail/>
          Page::Timeline => <Pages.Timeline/>
          Page::PostGrumble => <Pages.PostGrumble/>
          Page::UserGrumbles => <Pages.UserGrumbles/>
          Page::Error(statusCode) => <Pages.Error statusCode={statusCode}/>
        }
      </section>
    </div>
  }

  style app {
    justify-content: center;
    flex-direction: column;
    align-items: center;
    display: flex;

    background-color: #eeeeee;

    height: 100vh;

    /* width: 100vw; */
    font-family: Open Sans;
    font-weight: bold;
  }

  fun index : Html {
    <div::app>
      <Logo/>

      <Info mainPath="source/Main.mint"/>

      <Link href="https://www.mint-lang.com/">
        "Learn Mint"
      </Link>

      <Link href={apiurl}>
        "api"
      </Link>
    </div>
  } where {
    apiurl =
      "#{@ENDPOINT}/api"
  }
}
