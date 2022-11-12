require "fileutils"
require_relative "mpv_cmd"

class LiteDB
  def initialize
    db.execute <<-SQL
      create table if not exists pos_cache (
        filepath string not null primary key,
        pos int
      );
    SQL
  end

  def db
    FileUtils.mkdir_p "tmp/"
    @db ||= SQLite3::Database.new "tmp/gnuplex.sqlite3"
  end

  def loadpos(filepath)
    res = db.execute <<-SQL, [filepath]
      SELECT pos FROM pos_cache WHERE filepath = ?;
    SQL
    res.flatten[0] || 0
  end

  def savepos(filepath, pos)
    db.execute <<-SQL, [filepath, pos]
      INSERT OR REPLACE INTO pos_cache (filepath, pos) VALUES (?, ?);
    SQL
  end
end
