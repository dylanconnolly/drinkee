go/gin backend service

# API

- [GET /drinks](#get-drinks)
- [POST /drinks](#post-drinks)
- [GET /drinks/:id](#get-drinksid)
- [POST generateDrinks](#post-generatedrinks)

## Drinks Endpoints
### `GET drinks`

Request:
```
curl -X GET "localhost:8080/api/v1/drinks"
```

Response:
```
[
  ...,
  {
    "id": 85,
    "name": "dirty martini",
    "displayName": "Dirty Martini",
    "instructions": "Pour the vodka, dry vermouth and olive brine into a cocktail shaker with a handful of ice and shake well.\r\nRub the rim of a martini glass with the wedge of lemon.\r\nStrain the contents of the cocktail shaker into the glass and add the olive.\r\nA dirty Martini contains a splash of olive brine or olive juice and is typically garnished with an olive.",
    "drinkIngredients": [
        {
            "name": "vodka",
            "displayName": "Vodka",
            "measurement": "70ml/2fl oz"
        },
        {
            "name": "dry vermouth",
            "displayName": "Dry Vermouth",
            "measurement": "1 tbsp"
        },
        {
            "name": "olive brine",
            "displayName": "Olive Brine",
            "measurement": "2 tbsp"
        },
        {
            "name": "lemon",
            "displayName": "Lemon",
            "measurement": "1 wedge"
        },
        {
            "name": "olive",
            "displayName": "Olive",
            "measurement": "1"
        }
    ]
  },
  ...
]
```

### `POST drinks`

Request:
```curl
curl -X POST "localhost:8080/api/v1/drinks" \
  -H "Content-Type: application/json" \
  -d '{
        "name": "create new drink",
        "displayName": "Create New Drink",
        "description": "New Drink Description",
        "instructions": "Combine all ingredients in a glass of your choosing.",
        "drinkIngredients": [
          {"name": "Mango", "measurement": "3 parts"},
          {"name": "Vodka", "measurement": "1 part"}
        ]
      }'
```

### `GET drinks/:id`

Request:
```
curl -X GET "localhost:8080/api/v1/drinks/14"
```

Response:
```
{
  "id": 14,
  "name": "acapulco",
  "displayName": "Acapulco",
  "instructions": "Combine and shake all ingredients (except mint) with ice and strain into an old-fashioned glass over ice cubes. Add the sprig of mint and serve.",
  "drinkIngredients": [
    {
        "name": "egg white",
        "displayName": "Egg White",
        "measurement": "1"
    },
    {
        "name": "light rum",
        "displayName": "Light rum",
        "measurement": "1 1/2 oz"
    },
    {
        "name": "triple sec",
        "displayName": "Triple sec",
        "measurement": "1 1/2 tsp"
    },
    {
        "name": "lime juice",
        "displayName": "Lime juice",
        "measurement": "1 tblsp"
    },
    {
        "name": "sugar",
        "displayName": "Sugar",
        "measurement": "1 tsp"
    },
    {
        "name": "mint",
        "displayName": "Mint",
        "measurement": "1"
    }
  ]
}
```

### `POST generateDrinks`

Request:
```
curl -X POST "localhost:8080/api/v1/generateDrinks" \ 
  -H "Content-Type: application/json" \
  -d '{
    "ingredients": [
      {"id": 2, "name": "Vodka"},
      {"id": 7, "name": "Ginger beer"},
      {"id": 8, "name": "Lime"},
      {"id": 9, "name: "Vermouth"},
      {"id": 10, "name": "Olive"},
      {"id": 1, "name": "Mango"}
    ]
  }'
```

Response:
```
[
  {
      "id": 3,
      "name": "ace",
      "displayName": "Ace",
      "instructions": "Shake all the ingredients in a cocktail shaker and ice then strain in a cold glass.",
      "missingIngredientCount": 1,
      "haveIngredientCount": 4,
      "drinkIngredients": [
          {
              "name": "gin",
              "displayName": "Gin",
              "measurement": "2 shots"
          },
          {
              "name": "grenadine",
              "displayName": "Grenadine",
              "measurement": "1/2 shot"
          },
          {
              "name": "heavy cream",
              "displayName": "Heavy cream",
              "measurement": "1/2 shot"
          },
          {
              "name": "milk",
              "displayName": "Milk",
              "measurement": "1/2 shot"
          },
          {
              "name": "egg white",
              "displayName": "Egg White",
              "measurement": "1/2 fresh"
          }
      ]
  },
  {
      "id": 12,
      "name": "addison",
      "displayName": "Addison",
      "instructions": "Shake together all the ingredients and strain into a cold glass.",
      "missingIngredientCount": 1,
      "haveIngredientCount": 1,
      "drinkIngredients": [
          {
              "name": "gin",
              "displayName": "Gin",
              "measurement": "1 1/2 shot"
          },
          {
              "name": "vermouth",
              "displayName": "Vermouth",
              "measurement": "1 1/2 shot"
          }
      ]
  },
  ...
]
```

### Migrations

Postgres with sqlx + migrate

To create new migration
```
migrate create -ext sql -dir db/migrations -seq <migration_name>
```
To run migration
```
migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations up
```
To force DB to version
```
migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations force <version>
```
e.g. `migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations force 10`


### Testing
Integration tests use dockertest to spin up and tear down containers running instances of postgres DBs.

Test a single file
```
go test -race /path/to/file
```

for verbose output on tests:
```
go test /path/to/file -v
```
