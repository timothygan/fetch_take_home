# fetch_take_home
Timothy Gan's submission for Take Home Assignment for Fetch Rewards.

## Technology
* Go `v1.23`
* `gin` REST API framework.

## Prerequisites
* Docker

## Project Structure
* [http](https://github.com/timothygan/fetch_take_home/tree/main/internal/transport/http)
  contains the transport layer logic.
* [receipts](https://github.com/timothygan/fetch_take_home/tree/main/internal/receipts)
  contains the service layer logic.
* [db](https://github.com/timothygan/fetch_take_home/tree/main/internal/db)
  contains the db layer logic.
* [errors](https://github.com/timothygan/fetch_take_home/tree/main/errors)
  contains application error codes.
* [cmd](https://github.com/timothygan/fetch_take_home/blob/main/cmd/server/main.go)
  starts the actual server.

## Running the service
To build and start the service:
```docker build -t fetch . && docker run --rm -p 8080:8080 -t fetch```.
The service should now be running on `http://localhost:8080`, change `[your port here]:8080`
in the command accordingly if you want to use a different port for the application.

## Endpoints
### Endpoint: Process Receipts

* Path: `/receipts/process`
* Method: `POST`
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Takes in a JSON receipt and returns a JSON object with a UUID.
Example Payload:
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
Example Response:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```
If an invalid receipt is provided, the endpoint will return a `400` status code. 

### Endpoint: Get Points

* Path: `/receipts/{id}/points`
* Method: `GET`
* Response: A JSON object containing the number of points awarded.

Takes in a receipt ID and returns the number of points awarded for that receipt.

Example Response:
```json
{ "points": 28 }
```
If an invalid id is provided, the endpoint will return a `404` status code.

## Rules

These rules collectively define how many points should be awarded to a receipt.

* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents.
* 25 points if the total is a multiple of `0.25`.
* 5 points for every two items on the receipt.
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
* 6 points if the day in the purchase date is odd.
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.

## Some Extra Info
This was my first Go application! Patterns largely taken from [Go's tutorials](https://go.dev/),
Elliot Forbes's [example repo](https://github.com/TutorialEdge/go-rest-api-course), and Kristian Ott's
[example repo](https://github.com/kott/go-service-example/tree/main).