require 'sinatra/base'
# Our api helpers and stubbed responses

module Veracode
  module ApiHelpers
    def json(value)
      content_type :json
      JSON.generate value
    end
  end
end
