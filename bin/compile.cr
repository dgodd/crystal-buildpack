#!/usr/bin/env crystal
require "yaml"
require "file_utils"

build_dir, cache_dir = ARGV

Dir.cd(build_dir) do
  has_lib = Dir.exists?("lib")
  if !has_lib && Dir.exists?("#{cache_dir}/lib")
    puts "Copy cached lib dir"
    FileUtils.cp_r("#{cache_dir}/lib", "lib")
  end

  shard = YAML.parse(File.read("./shard.yml")) || raise "Could not read shard.yml"
  shard_name = shard["name"] || raise "Could not find shard name in shard.yml"
  crystal_bin = "#{cache_dir}/crystal/bin/crystal"

  puts "Installing Dependencies"
  Process.run(crystal_bin, ["deps", "--production"], output: true, error: true)

  puts "Compiling src/#{shard_name}.cr (auto-detected from shard.yml)"
  Process.run(crystal_bin, ["build", "src/#{shard_name}.cr", "--release", "-o", "app"], output: true, error: true)

  ## Copy to cache directory
  puts "Copy lib dir to cache"
  FileUtils.rm_rf("#{cache_dir}/lib")
  FileUtils.cp_r("lib", "#{cache_dir}/lib")
end

File.write("/tmp/crystal-buildpack-release-step.yml", %q{---
addons:
default_process_types:
  web: ./app --port $PORT
})
