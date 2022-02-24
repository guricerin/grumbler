component Pages.SignIn {
  state userId : String = ""
  state password : String = ""
  state apiStatus : Api.Status(User) = Api.Status::Initial

  fun setUserId (v : String) : Promise(Never, Void) {
    next { userId = v }
  }

  fun setPassword (v : String) : Promise(Never, Void) {
    next { password = v }
  }

  fun setApiStatus (v : Api.Status(User)) : Promise(Never, Void) {
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

  style content {
    flex-direction: column;
  }

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  fun submit : Promise(Never, Void) {
    sequence {
      signinUser =
        {
          id = userId,
          password = password
        }

      reqBody =
        encode signinUser

      status =
        Http.post("#{@ENDPOINT}/signin")
        |> Http.jsonBody(reqBody)
        |> Api.send(User.decodes)

      case (status) {
        Api.Status::Ok(user) => Application.signin(user)
        => setApiStatus(status)
      }
    }
  }

  style button {
    margin-top: 20px;
  }

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <{ error }>

        <form>
          <div class="field">
            <label class="label">
              "ユーザID"
            </label>

            <input
              class="input"
              type="text"
              placeholder="ユーザID"
              value={userId}
              onChange={handleInput(setUserId)}/>
          </div>

          <div class="field">
            <label class="label">
              "パスワード"
            </label>

            <input
              class="input"
              type="password"
              placeholder="パスワード"
              value={password}
              onChange={handleInput(setPassword)}/>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit}>

          <{ "サインイン" }>

        </button>
      </div>
    </div>
  }
}
