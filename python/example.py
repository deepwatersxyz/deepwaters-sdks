import requests, datetime
import auth

host_url = 'https://testnet.api.deepwaters.xyz'
api_route = '/rest/v1/'
# request these in the webapp
api_key = None
api_secret = None
nonce_d = {'nonce': None}
base_asset_id = 'WBTC.GOERLI.5.TESTNET.PROD'
quote_asset_id = 'USDC.GOERLI.5.TESTNET.PROD'

def sync_nonce(api_key, api_secret, nonce_d):
    extension = 'customer/api-key-status'
    request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)
    headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)
    r = requests.get(url, headers=headers)
    response = r.json()
    nonce_d['nonce'] = response['result']['nonce']

sync_nonce(api_key, api_secret, nonce_d)

# per-user
# requires authentication, without nonce
# GET /customer

extension = 'customer'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

# per API key
# requires authentication, without nonce
# GET /customer/api-key-status

extension = 'customer/api-key-status'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

# GET /pairs

request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, 'pairs')

print('GET %s ... ' % url)
r = requests.get(url)
response = r.json()
print(response)
print()

# GET /pairs/{pair_name}

pair_name = base_asset_id + '-' + quote_asset_id
extension = f'pairs/{pair_name}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

print('GET %s ... ' % url)
r = requests.get(url)
response = r.json()
print(response)
print()

# GET /pairs/{pair_name}/orderbook?depth={depth}

extension = f'pairs/{pair_name}/orderbook?depth=12'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

print('GET %s ... ' % url)
r = requests.get(url)
response = r.json()
print(response)
print()

# GET /assets

url = host_url + api_route + 'assets'

print('GET %s ... ' % url)
r = requests.get(url)
response = r.json()
print(response)
print()

# authentication required, including nonce
# POST /orders

# for fetching with filtering later
created_at_or_after_micros = auth.now_micros()

request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, 'orders')

payload = {"type": "LIMIT", "side": "BUY", "quantity": "1.00000", "price": "17000.00", "baseAssetID": base_asset_id, "quoteAssetID": quote_asset_id, "durationType": "GOOD_TILL_EXPIRY", "customerObjectID": auth.now_micros(), "expiresAtMicros": int(int(auth.now_micros()) + 10*1e6)}

headers = auth.get_authentication_headers(api_key, api_secret, 'POST', request_uri, nonce_d, payload)

print('POST %s ... ' % url)
r = requests.post(url, headers=headers, json=payload)
response = r.json()
print(response)
print()

# authentication required, including nonce
# delete all orders (for only a specific pair, when specified)
# DELETE /orders?pair={pair}

extension = f'orders?pair={pair_name}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'DELETE', request_uri, nonce_d)

print('DELETE %s with get params ... ' % url)
r = requests.delete(url, headers=headers)
response = r.json()
print(response)
print()

extension = 'orders'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'DELETE', request_uri, nonce_d)

print('DELETE %s without get params ... ' % url)
r = requests.delete(url, headers=headers)
response = r.json()
print(response)
print()

# per-user
# requires authentication, without nonce
# GET /orders?pair={pair}&type={type}&side={side}&status-in={status-in}&created-at-or-after-micros={created_at_or_after_micros}&created-before-micros={created_before_micros}&skip={skip}&limit={limit}

created_before_micros = auth.now_micros()

extension = f'orders?pair={pair_name}&type=LIMIT&side=BUY&status-in=ACTIVE-PARTIALLY_FILLED-FILLED&created-at-or-after-micros={created_at_or_after_micros}&created-before-micros={created_before_micros}&skip=0&limit=5'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

# per-user
# requires authentication, without nonce
# GET /trades?pair={pair}&type={type}&created-at-or-after-micros={created_at_or_after_micros}&created-before-micros={created_before_micros}&skip={skip}&limit={limit}

extension = f'trades?pair={pair_name}&type=FILL&created-at-or-after-micros={created_at_or_after_micros}&created-before-micros={created_before_micros}&skip=0&limit=1'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

# requires authentication
# GET /orders/by-venue-order-id/{venue_order_id} <- without nonce
# DELETE /orders/by-venue-order-id/{venue_order_id} <- with nonce

# first create an order to GET and then DELETE

extension = 'orders'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

payload = {'type': 'LIMIT', 'side': 'BUY',
'quantity': '1.0000', 'price': '17000.00', 'baseAssetID': base_asset_id, 'quoteAssetID': quote_asset_id}

headers = auth.get_authentication_headers(api_key, api_secret, 'POST', request_uri, nonce_d, payload)

print('POST %s ... ' % url)
r = requests.post(url, headers=headers, json=payload)
response = r.json()
print(response)
print()

venue_order_id = response['result']['venueOrderID']

extension = f'orders/by-venue-order-id/{venue_order_id}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

extension = f'orders/by-venue-order-id/{venue_order_id}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'DELETE', request_uri, nonce_d)

print('DELETE %s ... ' % url)
r = requests.delete(url, headers=headers)
response = r.json()
print(response)
print()

# requires authentication
# GET /orders/by-customer-object-id/{customer_object_id} <- without nonce
# DELETE /orders/by-customer-object-id/{customer_object_id} <- with nonce

# first create an order to GET and then DELETE

extension = 'orders'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

customer_object_id = str(datetime.datetime.now().timestamp())

payload = {'type': 'LIMIT', 'side': 'BUY',
'customerObjectID': customer_object_id,
'quantity': '10.0000', 'price': '17000.00', 'baseAssetID': base_asset_id, 'quoteAssetID': quote_asset_id}

headers = auth.get_authentication_headers(api_key, api_secret, 'POST', request_uri, nonce_d, payload)

print('POST %s ... ' % url)
r = requests.post(url, headers=headers, json=payload)
response = r.json()
print(response)
print()

extension = f'orders/by-customer-object-id/{customer_object_id}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'GET', request_uri)

print('GET %s ... ' % url)
r = requests.get(url, headers=headers)
response = r.json()
print(response)
print()

extension = f'orders/by-customer-object-id/{customer_object_id}'
request_uri, url = auth.get_request_uri_and_url_from_extension(host_url, api_route, extension)

headers = auth.get_authentication_headers(api_key, api_secret, 'DELETE', request_uri, nonce_d)

print('DELETE %s ... ' % url)
r = requests.delete(url, headers=headers)
response = r.json()
print(response)
print()
