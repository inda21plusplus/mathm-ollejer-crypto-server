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
    "ids": [base64string],
}
```

### Read file

Input:
```json
{
    "type": "read",
    "id": base64string,
}
```

Output:
```json
{
    "data": base64string,
    "signature": base64string,
    "hashes": [base64string], // den här är baklänges atm. alltså att hasharna längst upp kommer först
}
```

### Write file

Input:
```json
{
    "type": "write",
    "id": base64string,
    "data": base64string,
    "signature": base64string,
}
```

Output:
```json
{
    "hashes": [base64string], // den här är också baklänges
}
```

### Add file

TODO

### Remove file

TODO

### Move file

TODO?
