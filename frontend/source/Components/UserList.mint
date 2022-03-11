component Components.UserList {
  property users : Users = Users.empty()

  style text {
    white-space: pre-wrap;
  }

  fun userBox (user : User) : Html {
    <div class="box">
      <article class="media">
        <div class="media-content">
          <div class="content">
            <p>
              <a href="/user/#{user.id}">
                <strong>
                  <{ user.name }>
                </strong>

                <small>"@#{user.id}"</small>
              </a>

              <br/>
              <div::text>"#{user.profile}"</div>
            </p>
          </div>
        </div>
      </article>
    </div>
  }

  fun userListItem (user : User) : Html {
    <li>
      <{ userBox(user) }>
    </li>
  }

  fun render : Html {
    <div>
      <ul>
        <{ Array.map(userListItem, users.users) }>
      </ul>
    </div>
  }
}
