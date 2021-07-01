package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ResponseBody struct {
	ResponseCode int         `json:"ResponseCode"`
	Data         interface{} `json:"Data"`
}

type Car struct {
	LicenseNo string `json:"LicenseNo"`
	Colour    string `json:"Colour"`
}

type ParkingSlot struct {
	SlotNo int  `json:"SlotNo"`
	Car    *Car `json:"Car"`
}

type ParkingLot struct {
	TotalSlot      int           `json:"TotalSlot"`
	OccupiedSlot   int           `json:"OccupiedSlot"`
	NextFreeSlotNo int           `json:"NextFreeSlotNo"`
	ParkingSlot    []ParkingSlot `json:"ParkingSlot"`
}

var parkingLot = ParkingLot{}

func createParkingLot(w http.ResponseWriter, r *http.Request) {
	var tempParkingLot ParkingLot
	var totalSlot int
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Kindly enter size of parking lot",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}

	json.Unmarshal(reqBody, &tempParkingLot)
	if parkingLot.TotalSlot > 0 {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Parking lot has been created",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}
	totalSlot = tempParkingLot.TotalSlot
	parkingSlot := make([]ParkingSlot, totalSlot)

	for i := 0; i < len(parkingSlot); i++ {
		parkingSlot[i] = ParkingSlot{
			SlotNo: i,
			Car:    nil,
		}
	}

	parkingLot = ParkingLot{
		TotalSlot:      totalSlot,
		OccupiedSlot:   0,
		NextFreeSlotNo: 0,
		ParkingSlot:    parkingSlot,
	}

	responseBody := ResponseBody{
		ResponseCode: http.StatusCreated,
		Data:         parkingLot,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseBody)

}

func parkCarToParkingSpot(w http.ResponseWriter, r *http.Request) {
	var newCarParked Car
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Kindly input Car Object",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}

	json.Unmarshal(reqBody, &newCarParked)
	if parkingLot.TotalSlot == 0 {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Parking lot has not been created yet",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	} else if parkingLot.TotalSlot == parkingLot.OccupiedSlot {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Parking lot is full",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	} else {
		for i := 0; i < len(parkingLot.ParkingSlot); i++ {
			if parkingLot.ParkingSlot[i].SlotNo == parkingLot.NextFreeSlotNo {
				parkingLot.ParkingSlot[i].Car = &Car{
					LicenseNo: newCarParked.LicenseNo,
					Colour:    newCarParked.Colour,
				}
				parkingLot.OccupiedSlot = parkingLot.OccupiedSlot + 1

				for j := parkingLot.NextFreeSlotNo; j < len(parkingLot.ParkingSlot); j++ {
					if parkingLot.ParkingSlot[j].Car == nil {
						parkingLot.NextFreeSlotNo = j
						break
					}
				}
				responseBody := ResponseBody{
					ResponseCode: http.StatusCreated,
					Data:         newCarParked,
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(responseBody)
				return
			}
		}
	}

}

func leaveCarFromParkingLot(w http.ResponseWriter, r *http.Request) {
	parkingSlotNoLeave := mux.Vars(r)["ParkingSlotNo"]
	parkingSlotNumber, err := strconv.Atoi(parkingSlotNoLeave)

	if err == nil {
		if parkingLot.TotalSlot == 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusInternalServerError,
				Data:         "Parking lot has not been created yet",
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
		for i := 0; i < len(parkingLot.ParkingSlot); i++ {
			if parkingLot.ParkingSlot[i].SlotNo == parkingSlotNumber {
				if parkingLot.ParkingSlot[i].Car != nil {
					carTemp := parkingLot.ParkingSlot[i].Car
					parkingLot.ParkingSlot[i].Car = nil
					parkingLot.OccupiedSlot = parkingLot.OccupiedSlot - 1

					for j := 0; j < len(parkingLot.ParkingSlot); j++ {
						if parkingLot.ParkingSlot[j].Car == nil {
							parkingLot.NextFreeSlotNo = j
							break
						}
					}
					responseBody := ResponseBody{
						ResponseCode: http.StatusOK,
						Data:         carTemp,
					}
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(responseBody)
					return
				} else {
					responseBody := ResponseBody{
						ResponseCode: http.StatusNotFound,
						Data:         "No Car is Parked at Parking Slot No : " + parkingSlotNoLeave,
					}
					json.NewEncoder(w).Encode(responseBody)
					return
				}
			}
		}
	} else {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Please pass the parking slot number to leave",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}
}

func getParkingLot(w http.ResponseWriter, r *http.Request) {
	if parkingLot.TotalSlot == 0 {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Parking lot has not been created yet",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}

	responseBody := ResponseBody{
		ResponseCode: http.StatusOK,
		Data:         parkingLot,
	}
	json.NewEncoder(w).Encode(responseBody)
}

func getCarByColour(w http.ResponseWriter, r *http.Request) {
	carColourToFind := mux.Vars(r)["CarColour"]
	var result []Car

	if carColourToFind != "" {
		if parkingLot.TotalSlot == 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusInternalServerError,
				Data:         "Parking lot has not been created yet",
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
		for i := 0; i < len(parkingLot.ParkingSlot); i++ {
			if parkingLot.ParkingSlot[i].Car != nil {
				if strings.EqualFold(parkingLot.ParkingSlot[i].Car.Colour, carColourToFind) {
					result = append(result, *parkingLot.ParkingSlot[i].Car)
				}
			}
		}
		if len(result) > 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusOK,
				Data:         result,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseBody)
			return
		} else {
			responseBody := ResponseBody{
				ResponseCode: http.StatusNotFound,
				Data:         "Cannot find car with colour " + carColourToFind,
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
	} else {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Please pass the car colour to find",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}
}

func getSlotNoByColour(w http.ResponseWriter, r *http.Request) {
	carColourToFind := mux.Vars(r)["CarColour"]
	var result []int

	if carColourToFind != "" {
		if parkingLot.TotalSlot == 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusInternalServerError,
				Data:         "Parking lot has not been created yet",
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
		for i := 0; i < len(parkingLot.ParkingSlot); i++ {
			if parkingLot.ParkingSlot[i].Car != nil {
				if strings.EqualFold(parkingLot.ParkingSlot[i].Car.Colour, carColourToFind) {
					result = append(result, parkingLot.ParkingSlot[i].SlotNo)
				}
			}
		}
		if len(result) > 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusOK,
				Data:         result,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseBody)
			return
		} else {
			responseBody := ResponseBody{
				ResponseCode: http.StatusNotFound,
				Data:         "Cannot find slot no with colour " + carColourToFind,
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
	} else {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Please pass the car colour to find",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}
}

func getSlotNoByLicenseNo(w http.ResponseWriter, r *http.Request) {
	licenseNoToFind := mux.Vars(r)["LicenseNoToFind"]
	var result []int

	if licenseNoToFind != "" {
		if parkingLot.TotalSlot == 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusInternalServerError,
				Data:         "Parking lot has not been created yet",
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
		for i := 0; i < len(parkingLot.ParkingSlot); i++ {
			if parkingLot.ParkingSlot[i].Car != nil {
				if strings.EqualFold(parkingLot.ParkingSlot[i].Car.LicenseNo, licenseNoToFind) {
					result = append(result, parkingLot.ParkingSlot[i].SlotNo)
				}
			}
		}
		if len(result) > 0 {
			responseBody := ResponseBody{
				ResponseCode: http.StatusOK,
				Data:         result,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseBody)
			return
		} else {
			responseBody := ResponseBody{
				ResponseCode: http.StatusNotFound,
				Data:         "Cannot find slot no with license no " + licenseNoToFind,
			}
			json.NewEncoder(w).Encode(responseBody)
			return
		}
	} else {
		responseBody := ResponseBody{
			ResponseCode: http.StatusInternalServerError,
			Data:         "Please pass the car license no to find",
		}
		json.NewEncoder(w).Encode(responseBody)
		return
	}
}

func homelink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to parking lot API!")
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.HandleFunc("/", homelink)

	router.HandleFunc("/api/v1/parking-lot", createParkingLot).Methods("POST")
	router.HandleFunc("/api/v1/parking-lot", getParkingLot).Methods("GET")
	router.HandleFunc("/api/v1/parking-lot/park", parkCarToParkingSpot).Methods("POST")
	router.HandleFunc("/api/v1/parking-lot/leave/{ParkingSlotNo}", leaveCarFromParkingLot).Methods("PATCH")
	router.HandleFunc("/api/v1/parking-lot/find-car-by-colour/{CarColour}", getCarByColour).Methods("GET")
	router.HandleFunc("/api/v1/parking-lot/find-slot-no-by-colour/{CarColour}", getSlotNoByColour).Methods("GET")
	router.HandleFunc("/api/v1/parking-lot/find-slot-no-by-license-no/{LicenseNoToFind}", getSlotNoByLicenseNo).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
