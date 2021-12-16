component Pages.SignUp {
  state userId : String = ""
  state userName : String = ""
  state password : String = ""
  state apiStatus : Api.Status(User) = Api.Status::Initial

  fun setUserId (v : String) : Promise(Never, Void) {
    next { userId = v }
  }

  fun setUserName (v : String) : Promise(Never, Void) {
    next { userName = v }
  }

  fun setPassword (v : String) : Promise(Never, Void) {
    next { password = v }
  }

  fun setApiStatus (v : Api.Status(User)) : Promise(Never, Void) {
    next { apiStatus = v }
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
      signupUser =
        {
          id = userId,
          name = userName,
          password = password
        }

      reqBody =
        encode signupUser

      status =
        Http.post("#{@ENDPOINT}/signup")
        |> Http.jsonBody(reqBody)
        |> Api.send(User.decodes)

      case (status) {
        Api.Status::Ok(user) => Application.signin(user)
        => setApiStatus(status)
      }
    }
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

  style button {
    margin-top: 20px;
  }

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <form>
        <{error}>
          <div class="field">
            <label class="label">
              "User ID"
            </label>

            <input
              class="input"
              type="text"
              placeholder="User ID"
              value={userId}
              onChange={handleInput(setUserId)}/>
          </div>

          <div class="field">
            <label class="label">
              "User Name"
            </label>

            <input
              class="input"
              type="text"
              placeholder="User Name"
              value={userName}
              onChange={handleInput(setUserName)}/>
          </div>

          <div class="field">
            <label class="label">
              "Password"
            </label>

            <input
              class="input"
              type="password"
              placeholder="Password"
              value={password}
              onChange={handleInput(setPassword)}/>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={submit}>

          <{ "新規登録" }>

        </button>
      </div>
    </div>
  }
}
