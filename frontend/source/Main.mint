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

      <Footer/>
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

          Page::Timeline => <Pages.Timeline/>

          Page::NotFound => notFound()
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

  style notFound {
    height: 100vh;

    /* width: 100vw; */
    align-items: center;
    justify-content: center;
    display: flex;
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

  fun notFound : Html {
    <div::notFound>
      <p>"404 page not found."</p>
    </div>
  }
}
