component GrumbleList {
  connect Application exposing { userStatus }
  property grumbles : Grumbles = Grumbles.empty()

  fun doBookmark (grumble : Grumble, apiUrl : String, event : Html.Event) : Promise(Never, Void) {
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
            Http.post(apiUrl)
            |> Http.jsonBody(encode bookmarkReq)
            |> Api.send(BookmarkRes.decodes)

          case (status) {
            Api.Status::Initial => next { }
            Api.Status::Ok(res) => `location.reload()`
            Api.Status::Error(err) => Window.navigate("/")
          }
        }
    }
  }

  style iconNumber {
    margin-left: 5px;
  }

  fun bookmarkIcon (grumble : Grumble) : Html {
    if (grumble.isBookmarkedBySigninUser) {
      <div class="level-item">
        <a
          aria-label="bookmark"
          onClick={doBookmark(grumble, "#{@ENDPOINT}/auth/delete-bookmark")}>

          <span class="icon is-small">
            <i
              class="fas fa-bookmark"
              aria-hidden="true"/>
          </span>

          <span::iconNumber>
            <{ Number.toString(grumble.bookmarkedCount) }>
          </span>

        </a>
      </div>
    } else {
      <div class="level-item">
        <a
          aria-label="like"
          onClick={doBookmark(grumble, "#{@ENDPOINT}/auth/bookmark")}>

          <span class="icon is-small">
            <i
              class="far fa-bookmark"
              aria-hidden="true"/>
          </span>

          <span::iconNumber>
            <{ Number.toString(grumble.bookmarkedCount) }>
          </span>

        </a>
      </div>
    }
  }

  fun icons (grumble : Grumble) : Html {
    <nav class="level is-mobile">
      <div class="level-item">
        <a aria-label="reply">
          <span class="icon s-small">
            <i
              class="fas fa-reply"
              aria-hidden="true"/>
          </span>
        </a>
      </div>

      <div class="level-item">
        <a aria-label="retweet">
          <span class="icon is-small">
            <i
              class="fas fa-retweet"
              aria-hidden="true"/>
          </span>
        </a>
      </div>

      <{ bookmarkIcon(grumble) }>
    </nav>
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

          <{ icons(grumble) }>
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
