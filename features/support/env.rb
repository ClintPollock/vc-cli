require 'aruba/cucumber'

Before do
  @HOME = File.expand_path(File.join('~'))
  set_environment_variable 'HOME', @HOME

end

After do
  # Stop the server for every end of scenario
  @server.stop if (defined? @server) && @server
end
