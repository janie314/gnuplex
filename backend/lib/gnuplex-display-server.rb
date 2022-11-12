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

set :port, 50001
set :root, File.join(File.dirname(__FILE__), "..")

post "/play" do
  content_type :json
  mpvcmd.play
end

post "/pause" do
  content_type :json
  mpvcmd.pause
end

post "/queue" do
  content_type :json
  mpvcmd.queue params["mediafile"]
end

get "/vol" do
  content_type :json
  mpvcmd.getvol
end

post "/vol" do
  content_type :json
  mpvcmd.setvol params["vol"]
end

get "/pos" do
  content_type :json
  mpvcmd.getpos
end

post "/pos" do
  content_type :json
  mpvcmd.setvol params["pos"]
end
