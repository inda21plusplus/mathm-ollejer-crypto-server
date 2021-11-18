# Crypto File Server

Av Mathias

## Key exchange (RSA)

E - encryption key   (public)
D - decryption key   (private)
V - verification key (public)
S - signing key      (private)

c - Client
s - Server
a - 3rd party authority

Va - "fick på papper i det mörka gränden", hårdkodas / konfigureras i klienten
Sa(Es) sparas hos servern

1. client sends `(Ec + Vc)`

   T.ex: `{"e":{"n":"1234324","e":"3"},"v":{"n":"54894264","e":"65533"}}}`
2. server sends `Ec(Es + Vs + Sa(Es))`

   T.ex: Ec(`{"e":{"n":"5456232","e":"3"},"v":{"n":"69871564","e":"65533"},"s":"54ue56489uu5156i4i56464i1"}`)
3. client verifies `Es with Sa(Es)`
4. client sends `Es(Sc(8 rand bytes = R1))`
5. server sends `Ec(Ss(8 rand bytes = R2))`
6. `(R1+R2)` is used for ChaCha20-Poly1305 for all data send from here and
   onwards

## Authentication

TODO

## API

### List ids

Input:
```json
{
    "type": "list"
}
```

Output:
```json
{
    "ids": ["<base64data>"],
}
```

### Read file

Input:
```json
{
    "type": "read",
    "id": "<base64data>",
}
```

Output:
```json
{
    "data": "<base64data>",
    "signature": "<base64data>",
    "hashes": ["<base64data>"],
}
```

`hashes` är baklänges. Alltså att hasharna längst upp kommer först.

### Write file

Input:
```json
{
    "type": "write",
    "id": "<base64data>",
    "data": "<base64data>",
    "signature": "<base64data>",
}
```

Output:
```json
{
    "hashes": ["<base64data>"],
}
```

`hashes` är baklänges här också.

### Add file

TODO

### Remove file

TODO

### Move file

TODO?
