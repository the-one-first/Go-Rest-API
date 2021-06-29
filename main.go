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
		fmt.Fprintf(w, "Kindly enter size of parking lot")
		return
	}

	json.Unmarshal(reqBody, &tempParkingLot)
	if parkingLot.TotalSlot > 0 {
		fmt.Fprintf(w, "Parking lot has been created")
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
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(parkingLot)

}

func parkCarToParkingSpot(w http.ResponseWriter, r *http.Request) {
	var newCarParked Car
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly input Car Object")
		return
	}

	json.Unmarshal(reqBody, &newCarParked)
	if parkingLot.TotalSlot == 0 {
		fmt.Fprintf(w, "Parking lot has not been created yet")
		return
	} else if parkingLot.TotalSlot == parkingLot.OccupiedSlot {
		fmt.Fprintf(w, "Parking lot is full")
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
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(newCarParked)
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
			fmt.Fprintf(w, "Parking lot has not been created yet")
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
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(carTemp)
					return
				} else {
					fmt.Fprintf(w, "No Car is Parked at Parking Slot No : %v", parkingSlotNoLeave)
					return
				}
			}
		}
	} else {
		fmt.Fprintf(w, "Please pass the parking slot number to leave")
		return
	}
}

func getParkingLot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(parkingLot)
}

func getCarByColour(w http.ResponseWriter, r *http.Request) {
	carColourToFind := mux.Vars(r)["CarColour"]
	var result []Car

	if carColourToFind != "" {
		if parkingLot.TotalSlot == 0 {
			fmt.Fprintf(w, "Parking lot has not been created yet")
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
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return
		} else {
			fmt.Fprintf(w, "Cannot find car with colour %v", carColourToFind)
			return
		}
	} else {
		fmt.Fprintf(w, "Please pass the car colour to find")
		return
	}
}

func getSlotNoByColour(w http.ResponseWriter, r *http.Request) {
	carColourToFind := mux.Vars(r)["CarColour"]
	var result []int

	if carColourToFind != "" {
		if parkingLot.TotalSlot == 0 {
			fmt.Fprintf(w, "Parking lot has not been created yet")
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
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return
		} else {
			fmt.Fprintf(w, "Cannot find slot no with colour %v", carColourToFind)
			return
		}
	} else {
		fmt.Fprintf(w, "Please pass the car colour to find")
		return
	}
}

func getSlotNoByLicenseNo(w http.ResponseWriter, r *http.Request) {
	licenseNoToFind := mux.Vars(r)["LicenseNoToFind"]
	var result []int

	if licenseNoToFind != "" {
		if parkingLot.TotalSlot == 0 {
			fmt.Fprintf(w, "Parking lot has not been created yet")
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
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)
			return
		} else {
			fmt.Fprintf(w, "Cannot find slot no with license no %v", licenseNoToFind)
			return
		}
	} else {
		fmt.Fprintf(w, "Please pass the car colour to find")
		return
	}
}

func homelink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homelink)

	router.HandleFunc("/parking-lot", createParkingLot).Methods("POST")
	router.HandleFunc("/parking-lot", getParkingLot).Methods("GET")
	router.HandleFunc("/parking-lot/park", parkCarToParkingSpot).Methods("POST")
	router.HandleFunc("/parking-lot/leave/{ParkingSlotNo}", leaveCarFromParkingLot).Methods("PUT")
	router.HandleFunc("/parking-lot/find-car-by-colour/{CarColour}", getCarByColour).Methods("GET")
	router.HandleFunc("/parking-lot/find-slot-no-by-colour/{CarColour}", getSlotNoByColour).Methods("GET")
	router.HandleFunc("/parking-lot/find-slot-no-by-license-no/{LicenseNoToFind}", getSlotNoByLicenseNo).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
