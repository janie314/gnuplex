class DB
  # def self.f
  #    puts "FUCKYEAH"
  # end
  def db
    @db ||= SQLite3::Database.new "gnuplex.sqlite3"
  end
end
