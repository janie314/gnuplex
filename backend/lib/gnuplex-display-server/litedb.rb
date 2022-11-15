require "fileutils"
require "sqlite3"
require_relative "mpv_cmd"

class LiteDB
  def initialize
    db.execute <<-SQL
      create table if not exists pos_cache (
        filepath string not null primary key,
        pos int
      );
    SQL
    db.execute <<-SQL
      create table if not exists history (
	      id	integer not null unique,
      	mediafile	text,
	      primary key("id" AUTOINCREMENT)
      );
    SQL
  end

  def db
    FileUtils.mkdir_p "tmp/"
    @db ||= SQLite3::Database.new "tmp/gnuplex.sqlite3"
  end

  def loadpos(filepath)
    res = db.execute <<-SQL, [filepath]
      select pos from pos_cache where filepath = ?;
    SQL
    res.flatten[0] || 0
  end

  def savepos(filepath, pos)
    db.execute <<-SQL, [filepath, pos]
      insert or replace into pos_cache (filepath, pos) values (?, ?);
    SQL
  end

  def addhist(filepath)
    db.execute <<-SQL, [filepath]
      insert into history (mediafile) values (?);
    SQL
  end

  def last25
    res = []
    query = <<-SQL
      select mediafile from history order by id desc limit 25;
    SQL
    db.execute(query) do |row|
      res.append row
    end
    JSON.generate(res.flatten)
  end
end
