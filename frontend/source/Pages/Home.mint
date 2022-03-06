component Pages.Home {
  connect Application exposing { userStatus }

  style content {
    margin-bottom: 50px;
  }

  fun showGuide : Html {
    case (userStatus) {
      UserStatus::Guest =>
        <div>
          <li>
            "サインインは"

            <a href="/signin">
              "こちら。"
            </a>
          </li>

          <li>
            "アカウント新規登録は"

            <a href="/signup">
              "こちら。"
            </a>
          </li>
        </div>

      UserStatus::SignIn(user) =>
        <div>
          <li>
            "サインアウトは"

            <a href="/signout">
              "こちら。"
            </a>
          </li>

          <li>
            "退会は"

            <a href="/unsubscribe">
              "こちら。"
            </a>
          </li>
        </div>
    }
  }

  fun render : Html {
    <div>
      <div::content>
        <h1 class="title">
          "Welcome"
        </h1>

        <hr/>
        <p>"GrumblerはTwitterの劣化パクりWebアプリです。"</p>
        <p>"'Grumble' は「ぼやく」という意味らしいです。"</p>
      </div>

      <div::content>
        <h1 class="title">
          "Features"
        </h1>

        <hr/>

        <div class="content">
          <ul>
            <li>"無料です。"</li>
            <li>"DM、鍵アカウント、引用リツイートなどオリジナルと比べて色々機能が足りません。"</li>
            <li>"画像アップロード機能もありません。サーバ費用はできるだけ抑えたいので。"</li>
          </ul>
        </div>
      </div>

      <div::content>
        <h1 class="title">
          "Guide"
        </h1>

        <hr/>

        <div class="content">
          <ul>
            <{ showGuide() }>
            <li>"サインイン後、一週間以内に再アクセスしなかった場合、セッションが切れてサインアウト状態になります。"</li>
            <li>"本サービスは予告なく終了する場合があります。"</li>
          </ul>
        </div>
      </div>
    </div>
  }
}
