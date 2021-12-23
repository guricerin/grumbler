enum SearchKind {
  UserId
}

enum SearchResultKind {
  Initial
  Users(Users)
}

store Stores.Search {
  state apiStatus : Api.Status(SearchResultKind) = Api.Status::Initial

  fun setApiStatus (v : Api.Status(SearchResultKind)) : Promise(Never, Void) {
    next { apiStatus = v }
  }

  fun resetApiStatus : Promise(Never, Void) {
    next { apiStatus = Api.Status::Initial }
  }

  fun search (query : String, kind : String) : Promise(Never, Void) {
    case (kind) {
      "user_id" => searchByUserId(query)

      =>
        Application.setPage(Page::Error(400))
    }
  }

  fun searchByUserId (userId : String) : Promise(Never, Void) {
    sequence {
      status =
        Http.get("#{@ENDPOINT}/auth/search?q=#{userId}&k=user_id")
        |> Api.send(Users.decodes)

      case (status) {
        Api.Status::Ok(users) =>
          sequence {
            result =
              SearchResultKind::Users(users)

            next { apiStatus = Api.Status::Ok(result) }
          }

        Api.Status::Initial => next { apiStatus = Api.Status::Initial }
        Api.Status::Error(err) => next { apiStatus = Api.Status::Error(err) }
      }
    }
  }
}
