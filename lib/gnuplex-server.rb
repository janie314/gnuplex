class GNUPlexServer
  def self.run
    pid = spawn("mpv --idle=yes --input-ipc-server=/tmp/mpvsocket --fs")
    Process.wait pid
  end
end
