component Pages.Reply {
  connect Application exposing { userStatus }
  connect Stores.GrumbleDetail exposing { apiStatus }

  fun render : Html {
    case (apiStatus) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>

      Api.Status::Ok(grumbleDetail) => core(grumbleDetail.root)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun core (grumble : Grumble) : Html {
    case (userStatus) {
      /* unreachable! */
      UserStatus::Guest => Html.empty()

      UserStatus::SignIn(user) =>
        <div>
          <Components.GrumbleBox
            signinUser={user}
            grumble={grumble}/>

          <hr/>
        </div>
    }
  }
}
