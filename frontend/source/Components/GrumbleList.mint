component GrumbleList {
  connect Application exposing { userStatus }
  property grumbles : Grumbles = Grumbles.empty()

  fun doBookmark (grumble : Grumble, event : Html.Event) : Promise(Never, Void) {
    case (userStatus) {
      UserStatus::Guest => next { }

      UserStatus::SignIn(user) =>
        sequence {
          bookmarkReq =
            {
              grumblePk = grumble.pk,
              byUserId = user.id
            }

          status =
            Http.post("#{@ENDPOINT}/auth/bookmark")
            |> Http.jsonBody(encode bookmarkReq)
            |> Api.send(BookmarkRes.decodes)

          case (status) {
            Api.Status::Initial => next { }
            Api.Status::Ok(res) => Window.navigate("/settings")
            Api.Status::Error(err) => Window.navigate("/")
          }
        }
    }
  }

  /* --> テキストを親要素内で折り返し */
  style wrap {
    overflow-wrap: break-word;
  }

  style child {
    width: 100%;
  }

  /* <-- テキストを親要素内で折り返し */
  style text {
    white-space: pre-wrap;
  }

  style date {
    margin-left: 7px;
  }

  fun grumbleBox (grumble : Grumble) : Html {
    <div::wrap class="box">
      <article class="media">
        <div::child class="media-content">
          <div class="content">
            <p::text>
              <a href="/user/#{grumble.userId}">
                <strong>"#{grumble.userName}"</strong>
                <small>"@#{grumble.userId}"</small>
              </a>

              <small::date>"#{grumble.createdAt}"</small>

              <br/>
              <p>"#{grumble.content}"</p>
            </p>
          </div>

          <nav class="level is-mobile">
            <a
              class="level-item"
              aria-label="reply">

              <span class="icon is-small">
                <i
                  class="fas fa-reply"
                  aria-hidden="true"/>
              </span>

            </a>

            <a
              class="level-item"
              aria-label="retweet">

              <span class="icon is-small">
                <i
                  class="fas fa-retweet"
                  aria-hidden="true"/>
              </span>

            </a>

            <a
              class="level-item"
              aria-label="like"
              onClick={doBookmark(grumble)}>

              <span class="icon is-small">
                <i
                  class="fas fa-bookmark"
                  aria-hidden="true"/>
              </span>

            </a>
          </nav>
        </div>
      </article>
    </div>
  }

  fun grumbleListItem (grumble : Grumble) : Html {
    <div>
      <{ grumbleBox(grumble) }>
    </div>
  }

  fun render : Html {
    <div>
      <{ Array.map(grumbleListItem, grumbles.grumbles) }>
    </div>
  }
}
