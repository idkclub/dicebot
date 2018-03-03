from circleci/golang:1.10

workdir /home/circleci
run curl 'https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-1.9.62.zip' > google_appengine.zip && unzip -q google_appengine.zip && rm google_appengine.zip
