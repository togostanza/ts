#!/usr/bin/env ruby

require "erb"
require "json"

templates = {}

Dir["templates/*.{html,rq}"].each do |path|
  key = File.basename(path)
  templates[key] = File.read(path)
end

templates_json = JSON.dump(templates)

erb_template = File.read("index.html.erb")

puts "---"
puts ERB.new(erb_template).result(binding)
