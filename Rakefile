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
task :frontend_build do
  sh "bun i --cwd frontend"
  sh "bun run --cwd frontend build"
end

desc "build backend go code"
task :go_build do
  target = ENV.fetch("TARGET", "bin/gnuplex")
  sh "go", "build", "-C", "backend", "-o", target, "-ldflags", "-X main.SourceHash=" + source_hash, "."
end

desc "build backend (CI use only)"
task :go_build_ci do
  Dir.chdir(File.join(__dir__, "backend")) do
    sh "go", "build", "-o", "/tmp/gnuplex", "-ldflags", "-X main.SourceHash=" + source_hash, "."
  end
end

desc "build gnuplex"
task build: [:frontend_build, :go_build]

desc "is the go build current?"
task :go_build_current do
  exit 1 unless File.exist? File.join(__dir__, "backend/bin/gnuplex")
  build_hash = `./backend/bin/gnuplex -source_hash`.strip
  exit source_hash == build_hash
end

desc "print go source code hash"
task :go_source_hash do
  puts source_hash
end

def source_hash
  `find backend -type f -name '*.go' -or -name 'go*' | xargs sha512sum | cut -d ' ' -f 1 | sort | sha512sum | cut -d ' ' -f 1`.strip
end
