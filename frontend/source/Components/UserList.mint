component UserList {
  property users : Users = Users.empty()

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
              "#{user.profile}"
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
      <p>"検索結果 : #{Array.size(users.users)}件"</p>

      <ul>
        <{ Array.map(userListItem, users.users) }>
      </ul>
    </div>
  }
}
