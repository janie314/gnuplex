require "sinatra"

class GNUPlexDisplayServer
  def self.run
    spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
  end
end

GNUPlexDisplayServer.run

set :port, 50000

get "/test" do
  "OKAY"
end
