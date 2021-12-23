enum Api.Status(a) {
  Initial
  Error(Map(String, Array(String)))
  Ok(a)
}

module Api {
  fun errorStatus (key : String, val : String) : Api.Status(a) {
    Api.Status::Error(err)
  } where {
    err =
      Map.empty()
      |> Map.set(key, [val])
  }

  fun errorsOf (key : String, status : Api.Status(a)) : Array(String) {
    case (status) {
      Api.Status::Error(errors) =>
        errors
        |> Map.get(key)
        |> Maybe.withDefault([])

      => []
    }
  }

  fun decodeErrors (res : Http.Response) : Api.Status(a) {
    try {
      body =
        Json.parse(res.body)
        |> Maybe.toResult("decodeErrors() json parse error.")

      Debug.log(body)

      errors =
        decode body as Map(String, Array(String))

      Api.Status::Error(errors)
    } catch Object.Error => err {
      errorStatus("error", Object.Error.toString(err))
    } catch String => err {
      errorStatus("error", err)
    }
  }

  fun send (
    decoder : Function(Object, Result(Object.Error, a)),
    req : Http.Request
  ) : Promise(Never, Api.Status(a)) {
    sequence {
      res =
        req
        |> Http.header("Content-Type", "application/json")
        |> Http.withCredentials(true)
        |> Http.send()

      case (res.status) {
        200 =>
          try {
            obj =
              Json.parse(res.body)
              |> Maybe.toResult("response json parse error.")

            Debug.log(obj)

            data =
              decoder(obj)

            Api.Status::Ok(data)
          } catch {
            errorStatus("error", "someting went wrong.")
          }

        =>
          decodeErrors(res)
      }
    } catch Http.ErrorResponse => err {
      errorStatus("error", "status code: #{Number.toString(err.status)}, url: #{err.url}")
    }
  }

  fun getUser (userId : String) : Promise(Never, Api.Status(User)) {
    sequence {
      getUserReq =
        { id = userId }

      reqBody =
        encode getUserReq

      status =
        Http.get("#{@ENDPOINT}/user/#{userId}")
        |> Http.jsonBody(reqBody)
        |> Api.send(User.decodes)

      status
    }
  }
}
