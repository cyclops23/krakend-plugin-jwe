Plugin to allow krakend to process Core's JWE tokens

Requires use of same Go version as krakend - currently 1.20.11

Requires following environment variables to be set in order to validate JWEs from Core:
```
CORE_JWE_SYMMETRIC_KEY
CORE_JWE_SHARED_SECRET
```

Note Core does not currently include expiry time in the JWT...
