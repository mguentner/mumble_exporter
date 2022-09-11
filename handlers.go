package main

import (
	"encoding/binary"
	"net/http"
	"fmt"
	"net"
	"time"
	"bytes"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	TimeoutError = errors.New("Timeout Error")
)

const (
	MAX_UDP_PACKET_SIZE = 65536
)

func validateHost(input string) (string, error) {
	validate := validator.New()
	err := validate.Var(&input, "required,hostname_port")
	if err != nil {
		return "", err
	}
	return input, nil
}

type MumblePing struct {
	version string
	ts uint64
	users int32
	maxUsers int32
	bandwidthBitsPerSecond int32
	latencyMicroSeconds uint64
}

func decodeMumblePingFromBinary(buf []byte, timeOfArrival time.Time) (*MumblePing, error) {
	if len(buf) != 24 {
		return nil, errors.New("Expected exactly 24 bytes")
	}
	result := MumblePing{}
	result.version = string(buf[0:4])
	reader := bytes.NewReader(buf[4:])
	err := binary.Read(reader, binary.BigEndian, &result.ts)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &result.users)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &result.maxUsers)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &result.bandwidthBitsPerSecond)
	if err != nil {
		return nil, err
	}
	result.latencyMicroSeconds = uint64(timeOfArrival.Nanosecond() / 1000) - result.ts
	return &result, nil
}

func sendPing(host string) (*MumblePing, error) {
	conn, err := net.ListenUDP("udp", nil)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	destination, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, err
	}
	sendBuffer := &bytes.Buffer{}
	zero := uint32(0)
	now := uint64(time.Now().Nanosecond() / 1000)
	binary.Write(sendBuffer, binary.BigEndian, zero)
	binary.Write(sendBuffer, binary.BigEndian, now)
	_, err = conn.WriteTo(sendBuffer.Bytes(), destination)
	if err != nil {
		return nil, err
	}
	readBuffer := make([]byte, MAX_UDP_PACKET_SIZE)
	err = conn.SetReadDeadline(time.Now().Add(3*time.Second))
	if err != nil {
		return nil, err
	}
	readBytes, _, err := conn.ReadFromUDP(readBuffer)
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return nil, TimeoutError
	}
	if err != nil {
		return nil, err
	}
	timeOfArrival := time.Now()
	decoded, err := decodeMumblePingFromBinary(readBuffer[0:readBytes], timeOfArrival)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func AddMetricsToRegistry(mumblePingResult *MumblePing, host string, registry *prometheus.Registry) error {
	addAndSetGauge := func(name string, help string, value float64) error {
		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:  "mumble",
			Name:        name,
			Help:        help,
			ConstLabels: map[string]string{
				"host": host,
			},
		})
		err := registry.Register(g)
		if err != nil {
			return err
		}
		g.Set(value)
		return nil
	}
	err := addAndSetGauge("current_users", "currently connected users", float64(mumblePingResult.users))
	if err != nil {
		return err
	}
	err = addAndSetGauge("max_users", "maximum users (server limit)", float64(mumblePingResult.maxUsers))
	if err != nil {
		return err
	}
	err = addAndSetGauge("latency_microseconds", "latency to the server in micro seconds", float64(mumblePingResult.latencyMicroSeconds))
	if err != nil {
		return err
	}
	return nil
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	hostInput := r.URL.Query().Get("host")
	hostValidated, err := validateHost(hostInput)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err) , http.StatusBadRequest)
		return
	}
	result, err := sendPing(hostValidated)
	if err != nil {
		if err == TimeoutError {
			log.Warn().Msgf("Connection to %s timed out", hostValidated)
			http.Error(w, "Could not reach remote server", http.StatusGatewayTimeout)
			return
		}
		http.Error(w, fmt.Sprintf("%v", err) , http.StatusBadRequest)
		return
	}
	registry := prometheus.NewRegistry()
	err = AddMetricsToRegistry(result, hostValidated, registry)
	if err != nil {
		log.Error().Msgf("Error during registry creation: %v", err)
		http.Error(w, "ServerError", http.StatusInternalServerError)
		return
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
