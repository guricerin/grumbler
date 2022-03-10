component GrumbleBox {
  property signinUser : User
  property grumble : Grumble

  fun doBookmark (apiUrl : String, event : Html.Event) : Promise(Never, Void) {
    sequence {
      bookmarkReq =
        {
          grumblePk = grumble.pk,
          byUserId = signinUser.id
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

  style iconNumber {
    margin-left: 5px;
  }

  fun bookmarkIcon : Html {
    if (grumble.isBookmarkedBySigninUser) {
      <div class="level-item">
        <a
          aria-label="bookmark"
          onClick={doBookmark("#{@ENDPOINT}/auth/delete-bookmark")}>

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
          onClick={doBookmark("#{@ENDPOINT}/auth/bookmark")}>

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

  fun icons : Html {
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

      <{ bookmarkIcon() }>
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

  style content {
    color: black;
  }

  fun render : Html {
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

              <a::content href="/user/#{grumble.userId}/grumble/#{grumble.pk}">
                <div>"#{grumble.content}"</div>
              </a>
            </p>
          </div>

          <{ icons() }>
        </div>
      </article>
    </div>
  }
}
