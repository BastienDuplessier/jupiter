defmodule API.Controller.Jobs do
  alias Geolocation

  def index(conn) do
    {lat, _} = Float.parse(conn.params["latitude"])
    {lon, _} = Float.parse(conn.params["longitude"])
    {rad, _} = Float.parse(conn.params["radius"])

    response = Geolocation.compute(lat, lon, rad) |> Poison.encode!
    Plug.Conn.send_resp(conn, 200, response)
  end
end
