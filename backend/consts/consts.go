package consts

import "fmt"

const (
	GNUPlexVersion      string = "0.100 Mellon (PQA1)"
	DBVersion           int    = 1
	DevStaticFilespath  string = "./public"
	ProdStaticFilespath string = "/var/gnuplex/home"
	IssuesURL           string = "https://github.com/janie314/gnuplex/issues"
)

func InitStatements() []string {
	return []string{
		"create table if not exists pos_cache (filepath string not null primary key, pos int);",
		"create table if not exists history (id integer not null unique, mediafile	text, primary key(id AUTOINCREMENT));",
		"create table if not exists medialist (filepath text not null,  primary key(filepath)) ;",
		"create table if not exists mediadirs (filepath text not null, primary key(filepath)) ;",
		"create table if not exists file_exts (ext text not null, exclude int, primary key(ext)) ;",
		"create table if not exists version_info (key string not null primary key, value string);",
		fmt.Sprintf("insert or ignore into version_info values ('db_schema_version', %d);", DBVersion)}
}
