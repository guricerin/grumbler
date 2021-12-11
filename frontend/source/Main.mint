component Main {
  connect Application exposing { page }

  style app {
    justify-content: center;
    flex-direction: column;
    align-items: center;
    display: flex;

    background-color: #282C34;
    height: 100vh;
    width: 100vw;

    font-family: Open Sans;
    font-weight: bold;
  }

  fun render : Html {
     /* index() */
     case (page) {
       Page::Initial => Html.empty()

       Page::Home => index()

       Page::NotFound => notFound()
     }
  }

  fun index : Html {
    <div::app>
      <Logo/>

      <Info mainPath="source/Main.mint"/>

      <Link href="https://www.mint-lang.com/">
        "Learn Mint"
      </Link>
      <Link href={ apiurl }>
        "api"
      </Link>

      <Footer/>
    </div>
  } where {
    apiurl = "#{ @ENDPOINT }/api"
  }

  fun notFound : Html {
    <div>
      <p>"what's the fuck."</p>
    </div>
  }
}
