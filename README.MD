## Pour setup le projet : 
- go
- une base de données mysql (docker-compose.yml, adminer sur localhost:8080)
```docker-compose up``` pour lancer un container mysql
- Adminer :  
    * Systeme : mysql
    * Serveur : db
    * Utilisateur : user
    * Mdp : password
    * Database : goapi
    
- Se rendre dans ```src``` et lancer ```go run .```

## Les endpoints API

### Product

* **POST** localhost:3333/api/products:
    * fields : name (string), price(float)
    * creation product
* **PUT** localhost:3333/api/products/:id
    * id : int
    * update product
* **DELETE** localhost:3333/api/products/:id
    * id : int
    * supprime product
* **GET** localhost:333/api/products/
    * liste tous les products
* **GET** localhost:3333/api/products/:id
    * id : int
    * renvoie le produit demandé
### Payement

* **POST** localhost:3333/api/payments 
    * fields : productid(int), pricepaid(float)
    * creation payment
* **PUT** localhost:3333/api/payments/:id
    * fields : productid(int), pricepaid(float)
    * update payment
* **DELETE** localhost:3333/api/payments/:id
    * id : int
    * supprime payment
* **GET** localhost:3333/api/payments
    * liste tous les payments
* **GET** localhost:3333/api/payments/:id
    * id : int
    * renvoie le payment demandé
* **GET** localhost:3333/api/payments/stream
    * Header nécessaire : 
        - Accept: text/event-stream

* Commande pour écouter le SSE depuis un terminal :
    ```curl -H "Accept: text/event-stream" -N http://localhost:3333/api/payments/stream```