# Parking Lot API

## How to compile

Go to root directory where main.go is located. Then run this command in cmd terminal :

```bash
go build
```

## How to run application

Go to root directory where go-rest-api.exe is located. Then run this command in cmd terminal :

```
.\go-rest-api.exe
```

## How to test service

1. Open postman
2. Import Go-Rest-API.postman_collection.json
3. Test the service by click Send button

## List of Service

1. CreateParkingLot (/api/v1/parking-lot) ; Method POST
2. GetParkingLot (/api/v1/parking-lot) ; Method GET
3. ParkCarInParkingLot (/api/v1/parking-lot/park) ; Method POST
4. LeaveCarFromParkingLot (/api/v1/parking-lot/leave/{ParkingSlotNo}) ; Method PATCH
5. GetCarByColour (/api/v1/parking-lot/find-car-by-colour/{CarColour}) ; Method GET
6. GetSlotNoByColour (/api/v1/parking-lot/find-slot-no-by-colour/{CarColour}) ; Method GET
7. GetSlotNoByLicenseNo (/api/v1/parking-lot/find-slot-no-by-license-no/{LicenseNoToFind}) ; Method GET