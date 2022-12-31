# coding: utf-8
require 'base64'
require 'fileutils'
require 'open-uri'
require 'json'
require 'aruba/cucumber'

def test_root
    expand_path('/')
end

class JSONVisitor
  def hash(h) end

  def array(a) end

  def traverse(v)
    case v
    when Array
      array(v)
      v.each do |e|
        traverse(e)
      end
    when Hash
      hash(v)
      v.each do |k, val|
        traverse(val)
      end
    end
  end
end

module Normalize
  def self.blacklisted?(key)
    [
        # ignore this element in metadata as it is changing for each request
        'requestDate'
    ].include? key
  end

  def self.sort_sets(h)
    SortSets.new.traverse h
  end

  def self.delete_keys(h)
    DeleteKeys.new.traverse h
  end

  def self.normalize(h)
    delete_keys(h)
    sort_sets(h)
  end


  Coords = Struct.new(:coordinateType, :coordinate1, :coordinate2, :version) do
    def <=>(other)
      [coordinateType, coordinate1, coordinate2, version] <=>
          [other.coordinateType, other.coordinate1, other.coordinate2, other.version]
    end

    def self.from_h(h)
      if h then
        # The values_at lets this not fail if the keys are out of order
        new(*h.values_at('coordinateType', 'coordinate1', 'coordinate2', 'version'))
      else
        empty
      end
    end

    def to_json
      to_h.to_json
    end

    def self.empty
      new('', '', '', '')
    end
  end

  class SortSets < JSONVisitor
    def hash(h)
      h.each do |k, v|
        if k == 'directs' then
          v.sort! do |a, b|
            ac = Coords.from_h(a['coords'])
            bc = Coords.from_h(b['coords'])
            ac <=> bc
          end
        elsif k == 'artifactRelationships' then
          v.sort! do |a, b|
            ac = a['child']
            bc = b['child']
            ac <=> bc
          end
        end
      end
    end
  end

  class DeleteKeys < JSONVisitor
    def hash(h)
      keys = h.keys
      keys.each do |k|
        h.delete k if Normalize.blacklisted? k
      end
    end
  end
end

Given('a server that authenticates the CLI') do 
  tries = 3
  begin
    retries ||= 0
    @server = Veracode::LocalServer.start_sinatra(8080) do
      init
    end
  rescue
    sleep 2
    retry if (retries += 1) < tries
    raise "unable to start server after #{tries} tries"
  end
end

# This is a dummy endpoint. Modify accordingly based on backend's API design.
def init
  post('/v1/init') {
    json configFileContents: '',
         id: 123,
         property: 'value'
  }
end

Given('an activated user') do
  dot_veracode = File.join @HOME, '.veracode'
  creds_ini = File.join dot_veracode, 'credentials'
  FileUtils.mkdir_p dot_veracode
  File.open(creds_ini, "wb") do |f|
    f.write("veracode_api_key_id=#{getKey}\n" + "veracode_api_key_secret=#{getSecret}")
  end

  File.chmod 0600, creds_ini
end

Given('an user with invalid credentials') do
  dot_veracode = File.join @HOME, '.veracode'
  creds_ini = File.join dot_veracode, 'credentials'
  FileUtils.mkdir_p dot_veracode

  File.open(creds_ini, "wb") do |f|
    f.write <<-CREDS_INI
veracode_api_key_id=foo
veracode_api_key_secret=bar
CREDS_INI
  end

  File.chmod 0600, creds_ini
end

Given('an API server') do 
  tries = 3
  begin
    retries ||= 0
    @server = Veracode::LocalServer.start_sinatra(8080) do
      scan
    end
  rescue 
    puts "Error during processing: #{$!}"
    sleep 2
    retry if (retries += 1) < tries
    raise "unable to start server after #{tries} tries"
  end
end

Then("the output should match json in {string}") do |filename|
  expected_json = JSON.parse(File.read(expected_outputs(filename)))
  actual_json = JSON.parse(last_command_started.output)
  expect(Normalize.normalize(actual_json)).to match(Normalize.normalize(expected_json))
end

def expected_outputs(filename)
  "features/resources/expectedoutputs/%s/%s" % [host_os_architecture, filename]
end

def host_os_architecture
  arch_map = { "aarch64" => "arm64", "x86_64" => "x86_64" }
  arch = `docker info --format "{{.Architecture}}"`.gsub("\n", "")
  arch_map.has_key?(arch) ? arch_map[arch] : arch
end

# This is a dummy endpoint. Modify accordingly based on backend's API design.
def scan
  get('/v1/scan') {
    json test: '',
         scanId: 123
  }
end

Given('existing cache') do
  if Dir["#{Dir.home}/.veracode"].empty?
    Dir.mkdir("#{Dir.home}/.veracode")
    Dir.mkdir("#{Dir.home}/.veracode/cache")
  end
end

Then("the cache directory should be cleared") do
  expect(Dir.exist?("#{Dir.home}/.veracode/cache")).to match(false)
end

Then("the output should contain the version and hash") do
  version = ENV.fetch('CI_COMMIT_TAG',"0.0.0")
  hash = ENV.fetch('CI_COMMIT_SHORT_SHA',`git rev-parse --short HEAD`)

  expected_version = "Veracode CLI v" + version + " -- " + hash
  actual_version = last_command_started.output
  expect(actual_version).to match(expected_version)
end

def getCredsFromLocalFile
  result = []
  File.open('features/support/local_credentials').each do |line|
    result << line.split("=")[1].gsub(/[[:space:]]/, '')
  end
  return result
end

def getKey
  ENV.fetch('APIKEY', getCredsFromLocalFile[0])
end

def getSecret
  ENV.fetch('APISECRET', getCredsFromLocalFile[1])
end

When(/^I type key and secret/) do
  steps %Q{ When I type "#{getKey}" }
  steps %Q{ When I type "#{getSecret}" }
end

And(/^check key and secret present in the credentials file/) do
  steps %Q{ And the file "~/.veracode/credentials" should contain "veracode_api_key_id=#{getKey}" }
  steps %Q{ And the file "~/.veracode/credentials" should contain "veracode_api_key_secret=#{getSecret}" }
end