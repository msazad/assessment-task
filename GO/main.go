package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MACResponse struct {
	PrimaryMac   string `json:"primaryMac"`
	SecondaryMac string `json:"secondaryMac"`
}

func generateMACAddresses(w http.ResponseWriter, r *http.Request) {
	var data struct {
		CID string `json:"cid"`
	}

	// Parse the JSON input
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the input length
	if len(data.CID) != 32 {
		http.Error(w, "Invalid CID length", http.StatusBadRequest)
		return
	}

	// Extract the last 64 bits (last 16 hex digits)
	last64Bits := data.CID[16:]

	// Convert hex to integer
	macInt, err := strconv.ParseUint(last64Bits, 16, 64)
	if err != nil {
		http.Error(w, "Invalid CID format", http.StatusBadRequest)
		return
	}

	// Generate the primary MAC address (48 bits)
	//primaryMac := (macInt & 0xFFFFFFFFFFFF00) | 0x000001 // Ensure unicast
	primaryMac := (macInt & 0xFFFFFFFFFFFF) | 0x020000000000

	// Generate the secondary MAC by incrementing
	secondaryMac := primaryMac + 1

	// Convert back to MAC address format
	primaryMacAddr := formatMACAddress(primaryMac)
	secondaryMacAddr := formatMACAddress(secondaryMac)

	// Respond with JSON
	response := MACResponse{
		PrimaryMac:   primaryMacAddr,
		SecondaryMac: secondaryMacAddr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to format MAC address
func formatMACAddress(mac uint64) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		(mac>>40)&0xFF, (mac>>32)&0xFF, (mac>>24)&0xFF, (mac>>16)&0xFF, (mac>>8)&0xFF, mac&0xFF)
}

func main() {
    // Serve the HTML file at the root endpoint
    http.Handle("/", http.FileServer(http.Dir("./")))
    
    // Handle the MAC generation at /generate
    http.HandleFunc("/generate", generateMACAddresses)

    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

