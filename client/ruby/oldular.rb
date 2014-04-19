require 'net/http'
require 'base64'
require 'rest_client'
require 'json'


def help
  p 'call me with 1/2 parameters - action and id'
  exit
end

server = 'localhost:8080'
room = 'test'
action = ARGV[0]

help unless action

id = ARGV[1] if ARGV[1]
url = "#{server}/#{room}/#{action}" unless id
url = "#{server}/#{room}/#{action}/#{id}" if id

auth = 'Basic ' + Base64.encode64( 'test:test2' ).chomp

p url

resource = RestClient::Resource.new( url )
response = resource.get( :Authorization => auth )


my_hash = JSON.parse(response)

puts JSON.pretty_generate(my_hash)