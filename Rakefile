require "peppermint/rake"
require "fileutils"

desc "start development server"
task :dev do
  Dir.chdir(__dir__) do
    sh "bun i --cwd frontend"
    FileUtils.mkdir_p "tmp"
    caddy = Process.spawn "caddy run"
    frontend = Process.spawn "bun run --cwd frontend dev"
    backend = Process.spawn "go run -C backend . -verbose"
    Signal.trap("TERM") {
      [caddy, frontend, backend].each { |p| Process.kill "HUP", p }
      exit
    }
    [caddy, frontend, backend].each { |p| Process.waitpid p }
  end
end

desc "start development server (compiled build)"
task dev_compiled: [:build] do
  Dir.chdir(__dir__) do
    FileUtils.mkdir_p "tmp"
    caddy = Process.spawn "caddy run"
    backend = Process.spawn "./backend/bin/gnuplex -verbose -static_files ./backend/static"
    Signal.trap("TERM") {
      [caddy, backend].each { |p| Process.kill "HUP", p }
      exit
    }
    [caddy, backend].each { |p| Process.waitpid p }
  end
end

desc "build frontend"
task :frontend_build do
  Dir.chdir(__dir__) do
    sh "bun i --cwd frontend"
    sh "bun run --cwd frontend build"
  end
end

desc "build backend go code"
task :go_build do
  Dir.chdir(__dir__) do
    target = ENV.fetch("TARGET", "bin/gnuplex")
    sh "go", "build", "-C", "backend", "-o", target, "-ldflags", "-X main.SourceHash=" + source_hash, "."
  end
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
  Dir.chdir(__dir__) do
    exit 1 unless File.exist? File.join(__dir__, "backend/bin/gnuplex")
    build_hash = `./backend/bin/gnuplex -source_hash`.strip
    exit source_hash == build_hash
  end
end

desc "print go source code hash"
task :go_source_hash do
  puts source_hash
end

def source_hash
  Dir.chdir(__dir__) do
    `find backend -type f -name '*.go' | sort | while read i; do echo "$i"; cat "$i"; done | sha256sum | cut -d ' ' -f 1`.strip
  end
end
