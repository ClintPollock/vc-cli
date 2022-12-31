# based on <github.com/jnicklas/capybara/blob/ab62b27/lib/capybara/server.rb>
require 'fileutils'
require 'net/http'
require 'rack/handler/webrick'
require 'json'
require 'sinatra/base'
require 'rack/vcr'
require_relative 'api_helpers'

module Veracode
  class LocalServer
    class Identify < Struct.new(:app)
      def call(env)
        if env["PATH_INFO"] == "/__identify__"
          [200, {}, [app.object_id.to_s]]
        else
          app.call(env)
        end
      end
    end

    class JsonParamsParser < Struct.new(:app)
      def call(env)
        if env['rack.input'] and not input_parsed?(env) and type_match?(env)
          env['rack.request.form_input'] = env['rack.input']
          data = env['rack.input'].read
          env['rack.request.form_hash'] = data.empty? ? {} : JSON.parse(data)
        end
        app.call(env)
      end

      def input_parsed? env
        env['rack.request.form_input'].eql? env['rack.input']
      end

      def type_match? env
        type = env['CONTENT_TYPE'] and
            type.split(/\s*[;,]\s*/, 2).first.downcase =~ /[\/+]json$/
      end
    end

    class App < Sinatra::Base
      def invoke
        res = super
        content_type :json unless response.content_type
        response.body = '{}' if blank_response?(response.body) ||
            (response.body.respond_to?(:[]) && blank_response?(response.body[0]))
        res
      end

      def blank_response?(obj)
        obj.nil? || (obj.respond_to?(:empty?) && obj.empty?)
      end
    end

    def self.start_sinatra(cassette_name = nil, port, &block)
      klass = Class.new(App)
      klass.use JsonParamsParser

      if cassette_name then
        klass.use Rack::VCR, replay: true,
                  cassette: cassette_name, record: :none
      end
      klass.set :environment, :test
      klass.disable :protection
      klass.error(404, 401) { content_type :json; nil }
      klass.class_eval(&block)
      klass.before do
        if request.request_method == 'POST' then
          request.body.rewind
          @request_payload = JSON.parse request.body.read
        end
      end

      klass.helpers do
        include Veracode::ApiHelpers
      end

      new(klass.new, port).start
    end

    attr_reader :app, :host, :port, :server_thread
    attr_accessor :server

    def initialize(app, port,  host = 'localhost')
      @app = app
      @host = host
      @port = port
      @server = nil
      @server_thread = nil
    end

    def responsive?
      return false if @server_thread && @server_thread.join(0)

      res = Net::HTTP.start(host, port) { |http| http.get('/__identify__') }

      res.is_a?(Net::HTTPSuccess) and res.body == app.object_id.to_s
    rescue Errno::ECONNREFUSED, Errno::EBADF
      return false
    end

    def start
      @server_thread = start_handler(Identify.new(app)) do |server, host|
        self.server = server
      end

      Timeout.timeout(60) { @server_thread.join(0.01) until responsive? }
    rescue Timeout::Error
      raise "Rack application timed out during boot"
    else
      self
    end

    def start_handler(app)
      server = nil
      thread = Rack::Handler::WEBrick.run(app, server_options) { |s| server = s }
      addr = server.listeners[0].addr
      yield server, addr[3]
      return thread
    end

    def server_options
      {:Port => @port,
       :BindAddress => 'localhost',
       :ShutdownSocketWithoutClose => true,
       :ServerType => Thread,
       :AccessLog => [],
       :Logger => WEBrick::Log::new(nil, 0)
      }
    end

    def stop
      server.shutdown
      @server_thread.join
    end
  end
end