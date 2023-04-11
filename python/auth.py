import datetime, json
from secp256k1_utils import get_signature_hex_str

def now_micros():
    return str(int(datetime.datetime.now().timestamp() * 1e6))

# some requests require authentication.
# those that modify the customer state require that the next nonce value be submitted, as part of authentication.
# for example, getting the API key status does not require nonce submission,
# but submitting an order does require nonce submission.
# note: upon receiving an error, the user should immediately resync the nonce
# see "sync_nonce", below
# this function increments the nonce if the nonce_d argument is applied
def get_authentication_headers(api_key, api_secret, verb, request_uri, nonce_d = None, payload = None):
    
    headers = {'X-DW-APIKEY': api_key}

    now_micros_str = now_micros()
    headers['X-DW-TSUS'] = now_micros_str

    to_hash_and_sign = verb + request_uri.lower() + now_micros_str

    if nonce_d is not None:
        nonce_str = str(nonce_d['nonce'])
        headers['X-DW-NONCE'] = nonce_str
        to_hash_and_sign += nonce_str
        nonce_d['nonce'] += 1

    if payload is not None:
        to_hash_and_sign += json.dumps(payload)

    headers['X-DW-SIGHEX'] = get_signature_hex_str(to_hash_and_sign, api_secret)

    return headers
