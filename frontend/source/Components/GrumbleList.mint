component GrumbleList {
  property grumbles : Grumbles = Grumbles.empty()

  style wrap {
    overflow-wrap: break-word;
  }

  style date {
    margin: 7px;
  }

  fun grumbleBox (grumble : Grumble) : Html {
    <div class="box">
      <article class="media">
        <div class="media-content">
          <div class="content">
            <p>
              <a href="/user/#{grumble.userId}">
                // <strong>

                // <{ user.name }>

                // </strong>
                <small>"@#{grumble.userId}"</small>
              </a>

              <small::date>"#{grumble.createdAt}"</small>

              <br/>
              "#{grumble.content}"
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

  fun sample : Html {
    <div>
      <div class="box">
        <article class="media">
          <div class="media-left">
            <figure class="image is-64x64">
              <img
                src="https://bulma.io/images/placeholders/128x128.png"
                alt="Image"/>
            </figure>
          </div>

          <div class="media-content">
            <div class="content">
              <p>
                <strong>"John Smith"</strong>
                <small>"@johnsmith"</small>
                <small>"31m"</small>
                <br/>
                "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean efficitur sit amet massa fringilla egestas. Nullam condimentum luctus turpis.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
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
    </div>
  }

  fun render : Html {
    <div>
      <{ Array.map(grumbleListItem, grumbles.grumbles) }>
      <{ sample() }>
    </div>
  }
}
