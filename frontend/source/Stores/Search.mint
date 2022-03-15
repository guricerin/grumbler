enum SearchKind {
  UserId
}

enum SearchResultKind {
  Initial
  Users(Users)
  Grumbles(Grumbles)
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
      "user_id" => searchUser(query, kind)
      "user_name" => searchUser(query, kind)
      "grumble" => searchGrumble(query)

      =>
        Application.setPage(Page::Error(400))
    }
  }

  fun searchUser (keyword : String, kind : String) : Promise(Never, Void) {
    sequence {
      req =
        {
          keyword = keyword,
          kind = kind
        }

      status =
        Http.post("#{@ENDPOINT}/auth/search")
        |> Http.jsonBody(encode req)
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

  fun searchGrumble (keyword : String) : Promise(Never, Void) {
    sequence {
      req =
        {
          keyword = keyword,
          kind = "grumble"
        }

      status =
        Http.post("#{@ENDPOINT}/auth/search")
        |> Http.jsonBody(encode req)
        |> Api.send(Grumbles.decodes)

      case (status) {
        Api.Status::Ok(grumbles) =>
          sequence {
            result =
              SearchResultKind::Grumbles(grumbles)

            next { apiStatus = Api.Status::Ok(result) }
          }

        Api.Status::Initial => next { apiStatus = Api.Status::Initial }
        Api.Status::Error(err) => next { apiStatus = Api.Status::Error(err) }
      }
    }
  }
}
