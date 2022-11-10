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

get "/play" do
  MPVCmd.new.play
end

get "/pause" do
  MPVCmd.new.pause
end

get "/index" do
  erb "OKGREAT"
end
