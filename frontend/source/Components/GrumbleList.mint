component GrumbleList {
  property grumbles : Grumbles = Grumbles.empty()

  style wrap {
    overflow-wrap: break-word;
  }

  style text {
    white-space: pre-wrap;
  }

  style date {
    margin-left: 7px;
  }

  fun grumbleBox (grumble : Grumble) : Html {
    <div::wrap class="box">
      <article class="media">
        <div class="media-content">
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
            <div class="level-left">
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
                aria-label="like">

                <span class="icon is-small">
                  <i
                    class="fas fa-heart"
                    aria-hidden="true"/>
                </span>

              </a>
            </div>
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
