require "sinatra"
require_relative "gnuplex-display-server/mpv_cmd"
require_relative "gnuplex-display-server/litedb"
require "erb"

class GNUPlexDisplayServer
  def self.run
    loop do
      pid = spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs --save-position-on-quit", pgroup: nil)
      Process.wait pid
      sleep 3
    end
  end
end

def mpvcmd
  @mpvcmd ||= MPVCmd.new
end

def db
  @db ||= LiteDB.new
end

set :bind, "0.0.0.0"
set :port, 40000
set :root, File.join(File.dirname(__FILE__), "..")

get "/" do
  redirect "/index.html"
end

post "/api/play" do
  content_type :json
  mpvcmd.play
end

post "/api/pause" do
  content_type :json
  mpvcmd.pause
end

get "/api/media" do
  content_type :json
  mpvcmd.getmedia
end

post "/api/media" do
  content_type :json
  db.addhist params["mediafile"]
  mpvcmd.setmedia params["mediafile"]
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

get "/api/last25" do
  content_type :json
  db.last25
end

get "/api/medialist" do
  content_type :json
  db.medialist
end

post "/api/medialist" do
  content_type :json
  file = File.open(File.join(File.dirname(__FILE__), "medialist"))
  medialist = file.readlines.map(&:chomp)
  db.refresh_medialist medialist
  200
end

Thread.new {
  GNUPlexDisplayServer.run
}
