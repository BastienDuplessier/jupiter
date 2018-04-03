defmodule API.Router do
  use Plug.Router

  plug Plug.Logger
  plug Plug.Parsers, parsers: [:json], json_decoder: Poison
  plug :match
  plug :dispatch

  get "/jobs",   do: API.Controller.Jobs.index(conn)

  match _ do
    send_resp(conn, 404, "404 - Not Found")
  end
end
