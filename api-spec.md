# Tilgangsstyring

## Roller

- Authenticated user (f.eks. meg selv)
- Owner (er med i teamet som eier dataproduktet)

## Endepunkt

- `/dataproducts`
    - `POST /`: Authenticated user kan opprette et dataprodukt
    - `PUT /{productID}`: Owner kan oppdatere et dataprodukt
    - `DELETE /{productID}`: Owner kan slette et dataprodukt
- `/access`
    - `POST /{productID}`: Authenticated user kan opprette access til et dataprodukt for hvem som helst
    - `DELETE /{productID}`: Owner kan slette vilkårlig access, Authenticated user kan slette access for seg selv
