component UserList {
  property users : Users = Users.empty()

  fun userBox (user : User) : Html {
    <div class="box">
      <article class="media">
        <div class="media-content">
          <div class="content">
            <p>
              <strong>
                <{ user.name }>
              </strong>

              <small>"@#{user.id}"</small>
              <br/>
              "#{user.profile}"
            </p>
          </div>
        </div>
      </article>
    </div>
  }

  fun userList (user : User) : Html {
    <li>
      <{ userBox(user) }>
    </li>
  }

  fun render : Html {
    <div>
      <p>"検索結果 : #{Array.size(users.users)}件"</p>

      <ul>
        <{ Array.map(userList, users.users) }>
      </ul>
    </div>
  }
}
