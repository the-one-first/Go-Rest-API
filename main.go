package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func allocateCarToParkingSpot(w http.ResponseWriter, r *http.Request) {
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

func getParkingLot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(parkingLot)
}

func homelink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homelink)

	router.HandleFunc("/parking-lot", createParkingLot).Methods("POST")
	router.HandleFunc("/parking-lot", getParkingLot).Methods("GET")
	router.HandleFunc("/parking-a-car", allocateCarToParkingSpot).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
