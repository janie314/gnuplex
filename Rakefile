require "peppermint/rake"
require "fileutils"

desc "start development server"
task :dev do
  FileUtils.mkdir_p "tmp"
  caddy = Process.spawn "caddy run"
  frontend = Process.spawn "bun run --cwd frontend dev"
  backend = Process.spawn "go run -C backend ."
  Signal.trap("TERM") {
    [caddy, frontend, backend].each { |p| Process.kill "HUP", p }
    exit
  }
  [caddy, frontend, backend].each { |p| Process.waitpid p }
end

desc "build frontend"
task :build_frontend do
  sh "bun i --cwd frontend"
  sh "bun run --cwd frontend build"
end

desc "build backend"
task :build_backend do
  source_hash = `find backend -type f -not -path 'backend/bin/gnuplex' | xargs sha512sum | sha512sum | cut -d ' ' -f 1`.strip
  target = ENV.fetch("TARGET", "bin/gnuplex")
  sh "go", "build", "-C", "backend", "-o", target, "-ldflags", "-X main.SourceHash=" + source_hash, "."
end

desc "build gnuplex"
task build: [:build_frontend, :build_backend]

desc "is the go build current?"
task :go_build_current do
  exit 1 unless File.exist? File.join(__dir__, "backend/bin/gnuplex")
  puts "b"
  build_hash = `./backend/bin/gnuplex -source_hash`.strip
  puts source_hash
  puts build_hash
  puts "b"
  exit source_hash == build_hash
end

def source_hash
  `find backend -type f -not -path 'backend/bin/gnuplex' | xargs sha512sum | sha512sum | cut -d ' ' -f 1`.strip
end
