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

desc "build gnuplex"
task :build do
  sh "bun run --cwd frontend build"
  sh "go build -C backend ."
end
