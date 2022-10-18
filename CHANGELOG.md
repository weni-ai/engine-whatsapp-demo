v1.2.3
---------- 
 * Fix bug of downloading files in WhatsApp Demo channels #24

v1.2.2
---------- 
 * Add welcome message as environment variable #21

v1.2.1
----------
 * fix channel creation validation #19.

v1.2.0
----------
 * add http rest endpoint for channel creation #16.

v1.1.0
----------
 * add prometheus metrics for channel creations, contacts activations, contact messages and defaults.

v1.0.1
----------
 * added build workflows

v1.0.0
----------
 * models for contact, channel and config
 * handle whatsapp api webhook event callbacks
 * structured logging
 * mongoDB as database 
 * gRPC service to handle weni-integrations Channel creation
 * http service to handle courier requests and whatsapp api webhook event callback requests
 * contact token confirmation
 * docker image build
 * sentry integration
