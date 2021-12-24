component Pages.Home {
  style content {
    margin-bottom: 50px;
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
            <li>
              "ログインは"

              <a href="/sign-in">
                "こちら。"
              </a>
            </li>

            <li>
              "アカウント新規登録は"

              <a href="/sign-up">
                "こちら。"
              </a>
            </li>

            <li>"ログイン後、一週間以内に再アクセスしなかった場合、セッションが切れてログアウト状態になります。"</li>
          </ul>
        </div>
      </div>
    </div>
  }
}
