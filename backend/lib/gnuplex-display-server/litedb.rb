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
    db.execute <<~SQL
           create table if not exists "medialist" (
          	"filepath"	text not null,
      	    primary key("filepath")
      ) 
    SQL
  end

  def db
    FileUtils.mkdir_p File.join(File.dirname(__FILE__), "..", "..", "..", "tmp")
    filepath = File.join(File.dirname(__FILE__), "..", "..", "..", "tmp", "gnuplex.sqlite3")
    @db ||= SQLite3::Database.new filepath
  end

  def first_or_default(query, args, default)
    res = db.execute query, args
    res.flatten[0] || default
  end

  def read_query(query, args)
    res = []
    db.execute(query) do |row|
      res.append row
    end
    res.flatten
  end

  def loadpos(filepath)
    first_or_default <<-SQL, [filepath]
      select pos from pos_cache where filepath = ?;
    SQL
  end

  def savepos(filepath, pos)
    db.execute <<-SQL, [filepath, pos]
      insert or replace into pos_cache (filepath, pos) values (?, ?);
    SQL
  end

  def medialist
    read_query <<-SQL, []
      select distinct filepath from medialist order by filepath;
    SQL
  end

  def refresh_medialist(medialist)
    medialist.each do |filepath|
      db.execute <<-SQL, [filepath]
        insert or replace into medialist (filepath) values (?);
      SQL
    end
  end

  def addhist(filepath)
    db.execute <<-SQL, [filepath]
      insert into history (mediafile) values (?);
    SQL
  end

  def last25
    read_query <<-SQL, []
      select distinct mediafile from history order by id desc limit 25;
    SQL
  end
end
