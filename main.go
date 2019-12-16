package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	s "strings"
)

type Addr struct {
	Id      string  `json:"id"`
	Address string  `json:"addr"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Position struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type HereView struct {
	Results []HereResult `json:"Result"`
}

type HereResult struct {
	Location HereLocation `json:"Location"`
}

type HereLocation struct {
	Pos Position `json:"DisplayPosition"`
}

type HereResponse struct {
	Views []HereView `json:"View"`
}

type HereData struct {
	Response HereResponse `json:"Response"`
}

const baseApi string = "https://geocoder.ls.hereapi.com/6.2/geocode.json?"

func geocode(apiKey string, a Addr) Addr {
	v := url.Values{}
	v.Set("searchtext", a.Address)
	v.Set("gen", "9")
	v.Set("apiKey", apiKey)
	u := s.Join([]string{baseApi, v.Encode()}, "")

	res, err := http.Get(u)

	if err != nil {
		log.Fatal(err)
	}

	bod, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var hereData HereData
	err = json.Unmarshal(bod, &hereData)

	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("%v", hereData)

	if len(hereData.Response.Views) == 0 {
		log.Printf("bad response. View missing.")
	} else if len(hereData.Response.Views[0].Results) == 0 {
		log.Printf("No result for address")
	} else {
		a.Lat = hereData.Response.Views[0].Results[0].Location.Pos.Latitude
		a.Lon = hereData.Response.Views[0].Results[0].Location.Pos.Longitude
	}
	return a

}

func main() {
	in_path := flag.String("in", "", "path to the file with addresses to geocode. csv with pk id, address")
	out_path := flag.String("out", "./geocoded.json", "path to the file to write the coords. output pk id, coords")
	apiKey := flag.String("apikey", "", "HERE maps api key")

	flag.Parse()

	fmt.Printf("Reading %s\n", *in_path)

	infile, errIn := os.Open(*in_path)
	defer infile.Close()

	if errIn != nil {
		log.Panic(errIn)
	}

	fmt.Printf("Writing %s\n", *out_path)

	outfile, errOut := os.Create(*out_path)
	defer outfile.Close()

	if errOut != nil {
		log.Panic(errOut)
	}

	reader := csv.NewReader(bufio.NewReader(infile))
	writer := bufio.NewWriter(outfile)
	writer.WriteString("[")

	cnt := 0

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		addr := geocode(*apiKey, Addr{line[0], line[1], 0, 0})

		if cnt > 0 {
			writer.WriteString(",")
		}

		jsonAddr, err := json.Marshal(addr)

		if err != nil {
			log.Panic(err)
		}

		writer.Write(jsonAddr)
		cnt++
	}

	writer.WriteString("]")
	writer.Flush()
}
