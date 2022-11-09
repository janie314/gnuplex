require "json"
require "sqlite3"
require "socket"
require "thor"
require_relative "gnuplex-display-server/litedb"
require_relative "gnuplex-display-server/mpv_cmd"

class MyCLI < Thor
  desc "play", "play current media file"
  def play
    puts MPVCmd.new.play
  end

  desc "pause", "play current media file"
  def pause
    puts MPVCmd.new.pause
  end

  desc "queue", "queue a media file"
  option :mediafile
  def queue(mediafile)
    puts MPVCmd.new.queue mediafile
  end

  desc "getmedia", "get filepath currently playing"
  def getmedia
    puts MPVCmd.new.getmedia
  end

  desc "getvol", "get volume (0 to 100 scale)"
  def getvol
    puts MPVCmd.new.getvol
  end

  desc "setvol", "set volume (0 to 100 scale)"
  def setvol(vol)
    puts MPVCmd.new.setvol vol
  end

  desc "getpos", "get position (seconds)"
  def getpos
    puts MPVCmd.new.getpos
  end

  desc "setpos", "set position (seconds)"
  def setpos(pos)
    puts MPVCmd.new.setpos pos
  end

  desc "savepos", "[EXPERIMENTAL] save current mediafile's position"
  def savepos
    cmd = MPVCmd.new
    LiteDB.new.savepos(cmd.getmedia, cmd.getpos || 0)
  end

  desc "loadpos", "[EXPERIMENTAL] load current mediafile's saved position"
  def loadpos
    a = MPVCmd.new.getmedia
    puts "a", a
    b = LiteDB.new.loadpos a
    puts "b", b
    setpos b
  end
end
