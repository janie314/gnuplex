require "json"
require "sqlite3"
require "socket"
require "thor"
require_relative "mpv_cmds"

class MyCLI < Thor
  desc "play", "play current media file"
  def play
    MPVCmds.new.play
  end

  desc "pause", "play current media file"
  def pause
    MPVCmds.new.pause
  end

  desc "queue", "queue a media file"
  option :mediafile
  def queue(mediafile)
    MPVCmds.new.queue mediafile
  end

  desc "getmedia", "get filepath currently playing"
  def getmedia
    MPVCmds.new.getmedia
  end

  desc "getvol", "get volume (0 to 100 scale)"
  def getvol
    MPVCmds.new.getvol
  end

  desc "setvol", "set volume (0 to 100 scale)"
  def setvol(vol)
    MPVCmds.new.setvol vol
  end

  desc "getpos", "get position (seconds)"
  def getpos
    MPVCmds.new.getpos
  end

  desc "setpos", "set volume (0 to 100 scale)"
  def setpos(pos)
    MPVCmds.new.setpos pos
  end
end
