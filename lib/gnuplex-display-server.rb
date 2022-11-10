require "sinatra"
require_relative "gnuplex-display-server/mpv_cmd"
require "erb"

class GNUPlexDisplayServer
  def self.run
    spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
  end
end

GNUPlexDisplayServer.run

set :port, 50000
set :root, File.join(File.dirname(__FILE__), "..")

get "/play" do
  MPVCmd.new.play
end

get "/pause" do
  MPVCmd.new.pause
end

get "/index" do
  erb :index, locals: {x: 4444}
end
