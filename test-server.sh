mkdir test

# register device
curl -X POST -d '{"FirmwareVersion":"1"}' localhost/devices/dummy/register

# query device
curl "localhost/devices/dummy"

# upload firmware
curl -X POST --data-binary @firmware.bin "localhost/firmware/dummy/image"

# query firmware meta
curl "localhost/firmware/dummy"

# download firmware first 1024 bytes
curl "localhost/firmware/dummy/image/chunk?start=0&len=1024" > test/firmware_chunk.bin

# download whole firmware
curl "localhost/firmware/dummy/image" > test/firmware.bin
