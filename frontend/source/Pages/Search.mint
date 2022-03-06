component Pages.Search {
  state searchWord : String = ""

  connect Stores.Search exposing { apiStatus }

  fun setSearchWord (v : String) : Promise(Never, Void) {
    next { searchWord = v }
  }

  fun handleInput (
    onChange : Function(String, Promise(Never, Void)),
    event : Html.Event
  ) : a {
    onChange(Dom.getValue(event.target))
  }

  fun searchWithUserId : Promise(Never, Void) {
    sequence {
      Window.navigate("/search?q=#{searchWord}&k=user_id")
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

  style content {
    flex-direction: column;
  }

  style button {
    margin-top: 20px;
  }

  fun showStatus (status : Api.Status(SearchResultKind)) : Html {
    case (status) {
      Api.Status::Initial => Html.empty()
      Api.Status::Error(err) => <Errors errors={es}/>
      Api.Status::Ok(res) => showResult(res)
    }
  } where {
    es =
      Api.errorsOf("error", apiStatus)
  }

  fun showResult (res : SearchResultKind) : Html {
    case (res) {
      SearchResultKind::Initial => Html.empty()

      SearchResultKind::Users(users) =>
        <div>
          <p>"検索結果 : #{Array.size(users.users)}件"</p>
          <UserList users={users}/>
        </div>
    }
  }

  fun render : Html {
    <div::content class="column">
      <div class="box form-box">
        <{ error }>

        <form>
          <div class="field">
            <input
              class="input"
              type="text"
              placeholder="検索ワード"
              value={searchWord}
              onChange={handleInput(setSearchWord)}/>
          </div>
        </form>

        <button::button
          class="button is-primary"
          type="submit"
          onClick={searchWithUserId}>

          <{ "ID検索" }>

        </button>
      </div>

      <{ showStatus(apiStatus) }>
    </div>
  }
}
