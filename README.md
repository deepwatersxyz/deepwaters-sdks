# deepwaters-sdks

SDKs for interacting with the [deepwaters](https://deepwaters.xyz) exchange.

These SDKS are mainly intended to help with two things:
* REST API authentication, which requires digital signing
* GraphQL websocket subscriptions, for those who are new to these

Complete examples of these, as well as sample calls to all REST API endpoints, are presented here.

The REST API ([docs](https://rest.docs.api.deepwaters.xyz)) and the GraphQL API ([docs](https://docs.api.deepwaters.xyz)) can both be used for trading, but the REST API is simpler to use and is designed solely for trading. The GraphQL websockets subscriptions must be used for realtime trade and order updates - no authentication is required for these. Note that both L2 and L3 subscriptions are available.
