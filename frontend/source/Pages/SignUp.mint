component Pages.SignUp {
  connect Stores.SignUp exposing {
    setUserId,
    setUserName,
    setPassword,
    userId,
    userName,
    password,
    setApiStatus
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

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <form>
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

        <button
          class="button is-primary"
          type="submit"
          onClick={submit}>

          <{ "登録" }>

        </button>
      </div>
    </div>
  }
}
