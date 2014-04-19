require 'net/http'
require 'base64'
require 'rest_client'
require 'json'


server = "localhost:8080"
room = "test"
action = "look"
url = "#{server}/#{room}/#{action}"
auth = 'Basic ' + Base64.encode64( 'test:test2' ).chomp


resource = RestClient::Resource.new( url )
response = resource.get( :Authorization => auth )

p response