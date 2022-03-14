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
          <Components.GrumbleList grumbles={ancestors}/>

          <Components.GrumbleBox
            grumble={grumbleDetail.target}
            signinUser={user}/>

          <hr/>
          <p>"このぼやきに対する返信一覧"</p>
          <br/>
          <Components.GrumbleList grumbles={replies}/>
        </div>
    }
  } where {
    ancestors =
      Grumbles(grumbleDetail.ancestors)

    replies =
      Grumbles(grumbleDetail.replies)
  }
}
