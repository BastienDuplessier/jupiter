defmodule Geolocation do
  @moduledoc """
  Documentation for Geolocation.
  """

  @doc """
    Return job offers inside a specified zone.
  """
  def compute(latitude, longitude, radius) do
    all_jobs = :ets.match_object(:jobs, {:'_', :'$1'})
    Enum.filter(all_jobs, fn({_, job}) ->
      {lat, lon} = {job["office_latitude"], job["office_longitude"]}
      dist = Distance.GreatCircle.distance({latitude, longitude}, {lat, lon})
      dist < radius * 1000
    end) |> Enum.map(fn({_, job}) -> job end)
  end
                                
  @doc """
    Fills ETS table with data from jobs.csv
  """                            
  def init_jobs() do
    :ets.new(:jobs, [:set, :protected, :named_table])
    result = File.stream!("data/jobs.csv") |> CSV.decode(headers: true)
    Enum.map(result, fn({:ok, val}) ->
      lat = parse_latlon(val["office_latitude"])
      lon = parse_latlon(val["office_longitude"])
      job = Map.merge(val, %{"office_latitude" => lat, "office_longitude" => lon})
      id = System.unique_integer [:monotonic, :positive]
      :ets.insert(:jobs, {id, job})
    end)
  end


  defp parse_latlon(str) do
    case Float.parse(str) do
      :error -> 0.0
      {res, _} -> res
    end
  end
end
