# AGILE Device Manager service

An AGILE [DeviceManager](http://agile-iot.github.io/agile-api-spec/docs/html/api.html#iot_agile_DeviceManager) implementation.

***Currently under development***

##Requisites

- Go (>= v1.6) installed
- [Glide](https://github.com/Masterminds/glide) dependency manager

##Installation

- Clone the repository
- run `glide install` to sync the dependencies
- `go build`

##Running

Running the server
`./device-manager server`

Running the client (Currently just create a new Device instance)
`./device-manager client`

##License

MIT
