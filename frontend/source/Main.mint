component Main {
  connect Application exposing { page, userStatus }

  style app {
    font-family: Open Sans;
  }

  fun render : Html {
    <div::app>
      <Components.Header userStatus={userStatus}/>

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
          Page::DeleteGrumble => <Pages.DeleteGrumble/>
          Page::UserSettings => <Pages.UserSettings/>
          Page::GrumbleDetail => <Pages.GrumbleDetail/>
          Page::Reply => <Pages.Reply/>
          Page::Error(statusCode) => <Pages.Error statusCode={statusCode}/>
        }
      </section>
    </div>
  }
}
