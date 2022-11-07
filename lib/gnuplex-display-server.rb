class GNUPlexDisplayServer
  def self.run
    mpv_proc = spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
    Process.wait mpv_proc
  end
end
