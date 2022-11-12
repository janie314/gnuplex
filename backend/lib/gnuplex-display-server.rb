require "sinatra"
require_relative "gnuplex-display-server/mpv_cmd"
require "erb"

class GNUPlexDisplayServer
  def self.run
    spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
  end
end

GNUPlexDisplayServer.run

set :port, 50001
set :root, File.join(File.dirname(__FILE__), "..")

post "/play" do
  content_type :json
  MPVCmd.new.play
end

post "/pause" do
  content_type :json
  MPVCmd.new.pause
end

post "/queue" do
  content_type :json
  MPVCmd.new.queue params["mediafile"]
end
