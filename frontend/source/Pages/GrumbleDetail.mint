component Pages.GrumbleDetail {
  connect Application exposing { userStatus }
  connect Stores.GrumbleDetail exposing { apiStatus }

  fun render : Html {
    case (apiStatus) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>

      Api.Status::Ok(grumbleDetail) => core(grumbleDetail)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun core (grumbleDetail : GrumbleDetail) : Html {
    case (userStatus) {
      /* unreachable! */
      UserStatus::Guest => Html.empty()

      UserStatus::SignIn(user) =>
        <div>
          <Components.GrumbleBox
            signinUser={user}
            grumble={grumbleDetail.root}/>

          <hr/>
          <Components.GrumbleList grumbles={gs}/>
        </div>
    }
  } where {
    gs =
      Grumbles(grumbleDetail.replies)
  }
}
