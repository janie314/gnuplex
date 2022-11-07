class MPVCmd
  def mpvsocket
    @mpvsocket ||= UNIXSocket.new "/tmp/mpvsocket"
  end

  def mpvcmd(obj)
    mpvsocket.write(JSON.generate(obj) + "\n")
    jsonstr = mpvsocket.readline
    res = JSON.parse(jsonstr)
    res["data"]
  rescue => err
    warn err.message
  end

  def play
    mpvcmd({command: ["set_property", "pause", false]})
  end

  def pause
    mpvcmd({command: ["set_property", "pause", true]})
  end

  def queue(mediafile)
    mpvcmd({command: ["loadfile", mediafile]})
  end

  def getmedia
    mpvcmd({command: ["get_property", "path"]})
  end

  def getvol
    mpvcmd({command: ["get_property", "volume"]})
  end

  def setvol(vol)
    mpvcmd({command: ["set_property", "volume", vol]})
  end

  def getpos
    mpvcmd({command: ["get_property", "time-pos"]})
  end

  def setpos(pos)
    mpvcmd({command: ["set_property", "time-pos", pos]})
  end
end
