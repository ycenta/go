## Pour setup le projet : 
- besoin de go, et en gros écouter ce qu'il vous dit de taper dans la console 👍👍
- une base de données postgresql (docker-compose.yaml si vous avez pas, j'ai ouvert le port sur 5431, y'a un adminer sur localhost:8080 si jamais, sinon je conseil dbeaver)
```docker-compose up``` pour lancer un container postgre
- http://localhost:8080/ pour acceder à adminer (ou http://127.0.0.1:8080/) 
    * Systeme : Postgresql
    * Serveur : db
    * Utilisateur : postgres
    * Mdp : example
- Créer la base/ les table (a la main ou en requête)

```
CREATE SEQUENCE product_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."product" (
    "id" integer DEFAULT nextval('product_id_seq') NOT NULL,
    "name" character varying NOT NULL,
    "price" numeric NOT NULL,
    "createdAt" time without time zone NOT NULL,
    "updatedAt" time without time zone NOT NULL,
    CONSTRAINT "product_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

```
    
    
- ```go run .```

## Les endpoints API

### Product

* **POST** localhost:3333/product       -> Pour créer un produit, besoin de form fields : name, price
* **GET** localhost:3333/product/{id}   -> Pour obtenir un produit, besoin d'url param : id (int)
* **PUT** localhost:3333/product/{id}/update    -> Pour update un produit, besoin de l'url param :id (int), et  form fields : name, price 
* **DELETE** localhost:3333/product/{id}/delete     -> Pour supprimer un produit, besoin de l'urm param :id
* **GET** /products     -> besoin de rien (envie de toi)

### Payement
La même chose mais c'est pas fini

## todo
* Créer un routeur (un peu mieux qu'en dur comme j'ai fait, histoire d'avoir un code propre), au final j'ai utilisé mux, bien plus pratique que tout faire en vanilla
* Peut être changer de framework pour les routes (j'utilise mux mais dès ce mois-ci il est finito, go passer sur gin ?) (donc besoin de réecrire le code)
* Utiliser une lib json pour renvoyer du json et pas du texte tout moche comme je fais
* Renvoyer des objets après un create/update (ou au moins toutes les données d'une ligne de la bd) et pas seulement l'id comme je fais actuellement
* Faire les fonctions handlers pour le payement (comme celle pour le product)
* Rajouter de la verif sur les champs pour pas faire crash le serveur
* Créer la methode stream(payement) avec SSE pour renvoyer en direct dès qu'un paiement est crée (avec les channels & goroutines)
* Nettoyer le code oopsi : split en plusieurs fichiers pour que ça soit plus propre

