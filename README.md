# go-sampleapi

## Executar

Para rodar os containeres, utilize docker-compose up -d

### Endpoints

**POST** /accounts *(criação de uma conta)*
Request body:

```
{ 
    "document_number": "12345678900"
}
```

**GET** /accounts/:accountId *(consulta de informações de uma conta)* 
Response Body: 
```
{
    "account_id": 1, 
    "document_number": "12345678900"
}
```

**POST** /transactions *(criação de uma transação)* 
Request Body: 
```
{
    "account_id": 1, 
    "operation_type_id": 4, 
    "amount": 123.45
}
```
