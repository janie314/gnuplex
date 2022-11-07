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
    @db ||= SQLite3::Database.new "tmp/gnuplex.sqlite3"
  end

  def savepos(filepath, pos)
    db.execute <<-SQL, [filepath, pos]
      insert into pos_cache values (?, ?);
    SQL
  end
end
