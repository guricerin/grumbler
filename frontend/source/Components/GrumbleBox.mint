component Components.GrumbleBox {
  property signinUser : User
  property grumble : Grumble

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

  style both-ends {
    display: flex;
    justify-content: space-between;
  }

  fun render : Html {
    <div>
      <{ anchor() }>

      <div::wrap class="box">
        <article class="media">
          <div::child class="media-content">
            <div class="content">
              <p::text>
                <div::both-ends>
                  <div>
                    <a href="/user/#{grumble.userId}">
                      <strong>"#{grumble.userName}"</strong>
                      <small>"@#{grumble.userId}"</small>
                    </a>

                    <small::date>"#{grumble.createdAt}"</small>
                  </div>

                  <{ signinUserMenu() }>
                </div>

                <{ replyTo() }>

                <a::content href="/user/#{grumble.userId}/grumble/#{grumble.pk}/##{grumble.pk}">
                  <div>"#{grumble.content}"</div>
                </a>
              </p>
            </div>

            <{ icons() }>
          </div>
        </article>
      </div>
    </div>
  }

  fun signinUserMenu : Html {
    if (signinUser.id == grumble.userId) {
      <a href="/delete-grumble/#{grumble.pk}">
        <i
          class="fas fa-trash"
          aria-hidden="true"/>
      </a>
    } else {
      Html.empty()
    }
  }

  fun anchor : Html {
    <a name="#{grumble.pk}"/>
  }

  fun replyTo : Html {
    if (Grumble.isReply(grumble)) {
      <small>
        "返信先: "

        <a href="/user/#{grumble.reply.dstUserId}">
          "@#{grumble.reply.dstUserId}"
        </a>
      </small>
    } else {
      Html.empty()
    }
  }

  fun navigateToReplyPage (event : Html.Event) : Promise(Never, Void) {
    Window.navigate("/reply/#{grumble.pk}")
  }

  style iconNumber {
    margin-left: 5px;
  }

  fun replyIcon : Html {
    <div class="level-item">
      <a
        aria-label="reply"
        title="返信"
        onClick={navigateToReplyPage}>

        <span class="icon s-small">
          <i
            class="far fa-comment"
            aria-hidden="true"/>
        </span>

        <span::iconNumber>
          <{ Number.toString(grumble.reply.repliedCount) }>
        </span>

      </a>
    </div>
  }

  fun doRegrumble : Promise(Never, Void) {
    sequence {
      Window.confirm("リグランブルしますか？")

      regrumbleReq =
        { grumblePk = grumble.pk }

      status =
        Http.post("#{@ENDPOINT}/auth/regrumble")
        |> Http.jsonBody(encode regrumbleReq)
        |> Api.send(RegrumbleRes.decodes)

      case (status) {
        Api.Status::Initial => next { }
        Api.Status::Ok(res) => `location.reload()`
        Api.Status::Error(err) => Window.navigate("/")
      }
    } catch String => error {
      Promise.never()
    }
  }

  fun unRegrumble : Promise(Never, Void) {
    sequence {
      Window.confirm("リグランブルを取り消しますか？")

      regrumbleReq =
        { grumblePk = grumble.pk }

      status =
        Http.post("#{@ENDPOINT}/auth/delete-regrumble")
        |> Http.jsonBody(encode regrumbleReq)
        |> Api.send(RegrumbleRes.decodes)

      case (status) {
        Api.Status::Initial => next { }
        Api.Status::Ok(res) => `location.reload()`
        Api.Status::Error(err) => Window.navigate("/")
      }
    } catch String => error {
      Promise.never()
    }
  }

  style notRegrumbledBySigninUser {
    color: gray;
  }

  fun regrumbleIcon : Html {
    if (grumble.regrumble.isRegrumbledBySigninUser) {
      <div class="level-item">
        <a
          aria-label="regrumble"
          title="リグランブルを取り消す"
          onClick={unRegrumble}>

          <span class="icon is-small">
            <i
              class="fa fa-retweet"
              aria-hidden="true"/>
          </span>

          <span::iconNumber>
            <{ Number.toString(grumble.regrumble.regrumbledCount) }>
          </span>

        </a>
      </div>
    } else {
      <div class="level-item">
        <a
          aria-label="regrumble"
          title="リグランブル"
          onClick={doRegrumble}>

          <span class="icon is-small">
            <i::notRegrumbledBySigninUser
              class="fa fa-retweet"
              aria-hidden="true"/>
          </span>

          <span::iconNumber>
            <{ Number.toString(grumble.regrumble.regrumbledCount) }>
          </span>

        </a>
      </div>
    }
  }

  fun doBookmark (apiUrl : String, event : Html.Event) : Promise(Never, Void) {
    sequence {
      bookmarkReq =
        { grumblePk = grumble.pk }

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

  fun bookmarkIcon : Html {
    if (grumble.isBookmarkedBySigninUser) {
      <div class="level-item">
        <a
          aria-label="bookmark"
          title="ブックマークを取り消す"
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
          aria-label="bookmark"
          title="ブックマーク"
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
      <{ replyIcon() }>
      <{ regrumbleIcon() }>
      <{ bookmarkIcon() }>
    </nav>
  }
}
