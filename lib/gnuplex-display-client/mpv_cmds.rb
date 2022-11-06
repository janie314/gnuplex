class MPVCmds
  def mpvsocket
    @mpvsocket ||= UNIXSocket.new "/tmp/mpvsocket"
  end

  def sendcmd(obj)
    mpvsocket.write(JSON.generate(obj) + "\n")
    puts mpvsocket.readline
    mpvsocket.close
  end

  def play
    sendcmd({command: ["set_property", "pause", false]})
  end

  def pause
    sendcmd({command: ["set_property", "pause", true]})
  end

  def queue(mediafile)
    sendcmd({command: ["loadfile", mediafile]})
  end

  def getmedia
    sendcmd({command: ["get_property", "path"]})
  end

  def getvol
    sendcmd({command: ["get_property", "volume"]})
  end

  def setvol(vol)
    sendcmd({command: ["set_property", "volume", vol]})
  end

  def getpos
    sendcmd({command: ["get_property", "time-pos"]})
  end

  def setpos(pos)
    sendcmd({command: ["set_property", "time-pos", pos]})
  end
end
