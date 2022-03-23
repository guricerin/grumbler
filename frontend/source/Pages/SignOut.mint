component Pages.SignOut {
  connect Application exposing { userStatus, signout }

  state apiStatus : Api.Status(SignOutRes) = Api.Status::Initial

  fun setApiStatus (v : Api.Status(SignOutRes)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  get error : Html {
    case (apiStatus) {
      Api.Status::Error => <Errors errors={es}/>
      => Html.empty()
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun cancel (user : User, event : Html.Event) : Promise(Never, Void) {
    sequence {
      Window.navigate("/timeline")
    }
  }

  fun doSignOut (user : User, event : Html.Event) : Promise(Never, Void) {
    sequence {
      status =
        Http.post("#{@ENDPOINT}/auth/signout")
        |> Api.send(SignOutRes.decodes)

      case (status) {
        Api.Status::Ok(res) => signout()
        => setApiStatus(status)
      }
    }
  }

  style content {
    flex-direction: column;
  }

  style button {
    margin: 5px;
  }

  fun core (user : User) : Html {
    <div::content class="column">
      <div class="box form-box">
        <{ error }>
        <p>"サインアウトしますか？"</p>
        <br/>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={doSignOut(user)}>

          <{ "はい" }>

        </button>

        <button::button
          class="button"
          type="submit"
          onClick={cancel(user)}>

          <{ "いいえ" }>

        </button>
      </div>
    </div>
  }

  fun render : Html {
    case (userStatus) {
      /* unreachable */
      UserStatus::Guest => Html.empty()
      UserStatus::SignIn(user) => core(user)
    }
  }
}
