require "sinatra"
require_relative "gnuplex-display-server/mpv_cmd"
require "erb"

class GNUPlexDisplayServer
  def self.run
    spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
  end
end

GNUPlexDisplayServer.run

def mpvcmd
  @mpvcmd ||= MPVCmd.new
end

set :port, 40000
set :root, File.join(File.dirname(__FILE__), "..")

post "/api/play" do
  content_type :json
  mpvcmd.play
end

post "/api/pause" do
  content_type :json
  mpvcmd.pause
end

post "/api/queue" do
  content_type :json
  mpvcmd.queue params["mediafile"]
end

get "/api/vol" do
  content_type :json
  mpvcmd.getvol
end

post "/api/vol" do
  content_type :json
  mpvcmd.setvol params["vol"]
end

get "/api/pos" do
  content_type :json
  mpvcmd.getpos
end

post "/api/pos" do
  content_type :json
  mpvcmd.setpos params["pos"]
end
