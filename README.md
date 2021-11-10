# Crypto File Server

Av Mathias

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
