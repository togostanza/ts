#!/usr/bin/env ruby

require "erb"
require "json"

templates = {}

Dir["templates/*.{html,rq}"].each do |path|
  key = File.basename(path)
  templates[key] = File.read(path)
end

templates_json = JSON.dump(templates)
index_js = File.read("index.js")

erb_template = File.read("index.html.erb")

puts "---"
puts ERB.new(erb_template).result(binding)
