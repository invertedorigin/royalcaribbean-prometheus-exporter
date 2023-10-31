package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Sailing struct {
	BookingLink string `json:"bookingLink"`
	ID          string `json:"id"`
	Itinerary   struct {
		Code     string `json:"code"`
		Typename string `json:"__typename"`
	} `json:"itinerary"`
	SailDate              string `json:"sailDate"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	StateroomClassPricing []struct {
		Price struct {
			Value    int    `json:"value"`
			Typename string `json:"__typename"`
		} `json:"price"`
		StateroomClass struct {
			ID       string `json:"id"`
			Typename string `json:"__typename"`
		} `json:"stateroomClass"`
		Typename string `json:"__typename"`
	} `json:"stateroomClassPricing"`
	Typename string `json:"__typename"`
}

type StateroomClassPricing struct {
	Price struct {
		Value    int    `json:"value"`
		Typename string `json:"__typename"`
	} `json:"price"`
	StateroomClass struct {
		ID       string `json:"id"`
		Typename string `json:"__typename"`
	} `json:"stateroomClass"`
	Typename string `json:"__typename"`
}

type CruiseSearch struct {
	Data struct {
		CruiseSearch struct {
			Results struct {
				Cruises []struct {
					ID                 string `json:"id"`
					ProductViewLink    string `json:"productViewLink"`
					LowestPriceSailing struct {
						BookingLink               string `json:"bookingLink"`
						ID                        string `json:"id"`
						LowestStateroomClassPrice struct {
							Price struct {
								Value    int    `json:"value"`
								Typename string `json:"__typename"`
							} `json:"price"`
							StateroomClass struct {
								ID       string `json:"id"`
								Typename string `json:"__typename"`
							} `json:"stateroomClass"`
							Typename string `json:"__typename"`
						} `json:"lowestStateroomClassPrice"`
						SailDate     string `json:"sailDate"`
						StartDate    string `json:"startDate"`
						EndDate      string `json:"endDate"`
						TaxesAndFees struct {
							Value    float64 `json:"value"`
							Typename string  `json:"__typename"`
						} `json:"taxesAndFees"`
						TaxesAndFeesIncluded bool   `json:"taxesAndFeesIncluded"`
						Typename             string `json:"__typename"`
					} `json:"lowestPriceSailing"`
					MasterSailing struct {
						Itinerary struct {
							Code  string `json:"code"`
							Media struct {
								Images []struct {
									Path     string `json:"path"`
									Typename string `json:"__typename"`
								} `json:"images"`
								Typename string `json:"__typename"`
							} `json:"media"`
							Days []struct {
								Number int    `json:"number"`
								Type   string `json:"type"`
								Ports  []struct {
									Activity      string `json:"activity"`
									DepartureTime string `json:"departureTime"`
									Port          struct {
										Code   string `json:"code"`
										Name   string `json:"name"`
										Region string `json:"region"`
										Media  struct {
											Images []struct {
												Path     string `json:"path"`
												Typename string `json:"__typename"`
											} `json:"images"`
											Typename string `json:"__typename"`
										} `json:"media"`
										Typename string `json:"__typename"`
									} `json:"port"`
									Typename string `json:"__typename"`
								} `json:"ports"`
								Typename string `json:"__typename"`
							} `json:"days"`
							DeparturePort struct {
								Code     string `json:"code"`
								Name     string `json:"name"`
								Region   string `json:"region"`
								Typename string `json:"__typename"`
							} `json:"departurePort"`
							Destination struct {
								Code     string `json:"code"`
								Name     string `json:"name"`
								Typename string `json:"__typename"`
							} `json:"destination"`
							Name          string `json:"name"`
							SailingNights int    `json:"sailingNights"`
							Ship          struct {
								Code             string `json:"code"`
								Name             string `json:"name"`
								StateroomClasses []struct {
									ID      string `json:"id"`
									Name    string `json:"name"`
									Content struct {
										Amenities   []string `json:"amenities"`
										Code        string   `json:"code"`
										MaxCapacity string   `json:"maxCapacity"`
										Media       struct {
											Images []struct {
												Path string `json:"path"`
												Meta struct {
													Description string `json:"description"`
													Title       string `json:"title"`
													Location    string `json:"location"`
													Typename    string `json:"__typename"`
												} `json:"meta"`
												Typename string `json:"__typename"`
											} `json:"images"`
											Typename string `json:"__typename"`
										} `json:"media"`
										SuperCategory string `json:"superCategory"`
										Typename      string `json:"__typename"`
									} `json:"content"`
									Typename string `json:"__typename"`
								} `json:"stateroomClasses"`
								Media struct {
									Images []struct {
										Path     string `json:"path"`
										Typename string `json:"__typename"`
									} `json:"images"`
									Typename string `json:"__typename"`
								} `json:"media"`
								Typename string `json:"__typename"`
							} `json:"ship"`
							TotalNights int    `json:"totalNights"`
							Type        string `json:"type"`
							Typename    string `json:"__typename"`
						} `json:"itinerary"`
						Typename string `json:"__typename"`
					} `json:"masterSailing"`
					Sailings []struct {
						BookingLink string `json:"bookingLink"`
						ID          string `json:"id"`
						Itinerary   struct {
							Code     string `json:"code"`
							Typename string `json:"__typename"`
						} `json:"itinerary"`
						SailDate              string `json:"sailDate"`
						StartDate             string `json:"startDate"`
						EndDate               string `json:"endDate"`
						StateroomClassPricing []struct {
							Price struct {
								Value    int    `json:"value"`
								Typename string `json:"__typename"`
							} `json:"price"`
							StateroomClass struct {
								ID       string `json:"id"`
								Typename string `json:"__typename"`
							} `json:"stateroomClass"`
							Typename string `json:"__typename"`
						} `json:"stateroomClassPricing"`
						Typename string `json:"__typename"`
					} `json:"sailings"`
					Typename string `json:"__typename"`
				} `json:"cruises"`
				Total    int    `json:"total"`
				Typename string `json:"__typename"`
			} `json:"results"`
			Typename string `json:"__typename"`
		} `json:"cruiseSearch"`
	} `json:"data"`
}

type CruiseDetails struct {
	StepType  string `json:"stepType"`
	URLParams []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"urlParams"`
	Steps []struct {
		Label                     string `json:"label"`
		StepType                  string `json:"stepType"`
		ServicePath               string `json:"servicePath"`
		Active                    bool   `json:"active"`
		Completed                 bool   `json:"completed"`
		BookingStepMobileProgress struct {
			MobileLabelPart1 string `json:"mobileLabelPart1"`
			MobileLabelPart2 string `json:"mobileLabelPart2"`
		} `json:"bookingStepMobileProgress"`
		Disabled bool `json:"disabled"`
		External bool `json:"external"`
	} `json:"steps"`
	StepDetails struct {
		BookingFooter struct {
			TermAndConditionsHeading string   `json:"termAndConditionsHeading"`
			TermAndConditions        []string `json:"termAndConditions"`

			PriceOccupancyAlert string `json:"priceOccupancyAlert"`
		} `json:"bookingFooter"`
		ResetSession           bool `json:"resetSession"`
		SelectedStateroomIndex int  `json:"selectedStateroomIndex"`

		Resources struct {
			BofaURL string `json:"bofaUrl"`
		} `json:"resources"`
		SelectedCabinPrices struct {
		} `json:"selectedCabinPrices"`
		StateroomsInfo struct {
		} `json:"stateroomsInfo"`
		Toggles []struct {
			Name string `json:"name"`
			On   bool   `json:"on"`
		} `json:"toggles"`
		InFinalPaymentPeriod   bool   `json:"inFinalPaymentPeriod"`
		DisplayState           string `json:"displayState"`
		MaxOccupancyCategories struct {
		} `json:"maxOccupancyCategories"`
		Description             string `json:"description"`
		CabinClass              string `json:"cabinClass"`
		StateroomCategoryGroups []struct {
			StateroomType       string `json:"stateroomType"`
			Title               string `json:"title"`
			Description         string `json:"description"`
			StateroomCategories []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Images      []struct {
					ImagePath string `json:"imagePath"`
					AltText   string `json:"altText"`
				} `json:"images"`
				MoreDetailLabel string `json:"moreDetailLabel"`
				HideDetailLabel string `json:"hideDetailLabel"`
				Details         []struct {
					Code  string   `json:"code"`
					Items []string `json:"items"`
					Title string   `json:"title"`
				} `json:"details"`
				DeckplanCodes struct {
					Label  string `json:"label"`
					Images []struct {
						ImagePath string `json:"imagePath"`
					} `json:"images"`
				} `json:"deckplanCodes"`
				ShipAccessibilityAction struct {
					Label           string `json:"label"`
					Target          string `json:"target"`
					Disabled        bool   `json:"disabled"`
					ServicePath     string `json:"servicePath"`
					HTTPRequestType string `json:"httpRequestType"`
					DisplayType     string `json:"displayType"`
				} `json:"shipAccessibilityAction"`
				CategoryAccessibilityAction struct {
					Label           string `json:"label"`
					Target          string `json:"target"`
					Disabled        bool   `json:"disabled"`
					ServicePath     string `json:"servicePath"`
					HTTPRequestType string `json:"httpRequestType"`
					DisplayType     string `json:"displayType"`
				} `json:"categoryAccessibilityAction"`
				PriceLockup struct {
					Currency                  string  `json:"currency"`
					CurrencySymbol            string  `json:"currencySymbol"`
					OriginalPriceNbr          float64     `json:"originalPriceNbr"`
					PerPerson                 string  `json:"perPerson"`
					PricePerPersonNbr         float64 `json:"pricePerPersonNbr"`
					PriceIncludesTaxesAndFees bool    `json:"priceIncludesTaxesAndFees"`
					Prefix                    string  `json:"prefix"`
					Price                     string  `json:"price"`
					TaxesAndFees              string  `json:"taxesAndFees"`
					TaxesAndFeesNbr           float64 `json:"taxesAndFeesNbr"`
					TotalPriceNbr             float64 `json:"totalPriceNbr"`
					TotalDiscountNbr          float64     `json:"totalDiscountNbr"`
					Value                     float64     `json:"value"`
					SelectedFareCode          string  `json:"selectedFareCode"`
					NetPrice                  float64     `json:"netPrice"`
					RawPrice                  float64     `json:"rawPrice"`
					Credit                    bool    `json:"credit"`
				} `json:"priceLockup"`
				ConnectedRooms        bool   `json:"connectedRooms"`
				GuaranteeCategory     bool   `json:"guaranteeCategory"`
				FamilyEligible        bool   `json:"familyEligible"`
				StateroomCategoryCode string `json:"stateroomCategoryCode"`
				StateroomSubtypeCode  string `json:"stateroomSubtypeCode"`
				Accessible            bool   `json:"accessible"`
			} `json:"stateroomCategories"`
		} `json:"stateroomCategoryGroups"`
	}
}

type customMetric struct {
	url                  string
	status               float64
	totalMS              float64
	dnsMS                float64
	firstbyteMS          float64
	connectMS            float64
	price                float64
	cruiseID             string
	itinerary            string
	stateroomClass       string
	dateLabel            string
	ship                 string
	departurePort        string
	days                 string
	shipCode             string
	destinationCode      string
	stateroomDetailClass string
}

type Exporter struct {
	ctx                   context.Context
	urlStatus             *prometheus.GaugeVec
	urlMs                 *prometheus.GaugeVec
	urlDNS                *prometheus.GaugeVec
	urlFirstByte          *prometheus.GaugeVec
	urlConnectTime        *prometheus.GaugeVec
	royalPrice            *prometheus.GaugeVec
	urls                  []string
	healthcheck_invertval time.Duration
}

func NewExporter(ctx context.Context, inverval time.Duration, urls []string) (hc *Exporter) {
	hc = &Exporter{
		ctx: ctx,
		urlStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "url_status",
			Help:      "Status of the URL as a integer value",
		}, []string{"url"}),
		urlMs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "url_response_ms",
			Help:      "Response time in milliseconds it took for the URL to respond.",
		}, []string{"url"}),
		urlDNS: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "url_dns_ms",
			Help:      "Response time in milliseconds it took for the DNS request to take place.",
		}, []string{"url"}),
		urlFirstByte: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "url_first_byte_ms",
			Help:      "Response time in milliseconds it took to retrive the first byte.",
		}, []string{"url"}),
		urlConnectTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "url_connect_time_ms",
			Help:      "Response time in milliseconds it took to establish the inital connection.",
		}, []string{"url"}),
		royalPrice: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "royal",
			Subsystem: "detail",
			Name:      "price",
			Help:      "cabin price with labels",
		}, []string{"url", "cruiseid", "itinerary", "stateroomclass", "datelabel", "ship", "departureport", "days", "shipcode", "destinationcode", "stateroomdetailclass"}),
		healthcheck_invertval: inverval,
		urls:                  urls,
	}
	prometheus.MustRegister(hc.urlStatus, hc.urlMs, hc.urlDNS, hc.urlConnectTime, hc.urlFirstByte, hc.royalPrice)
	http.Handle("/metrics", promhttp.Handler())
	return hc
}

func paramsToQueryString(params map[string]string) string {
	var queryString string
	for key, value := range params {
		queryString += fmt.Sprintf("%s=%s&", key, value)
	}
	return queryString[:len(queryString)-1]
}

func (hc *Exporter) fetchDetails(ctx context.Context, shipName string, shipDestinationCode string, departurePort string, days int, ID string, shipCode string, sc Sailing, stateroom StateroomClassPricing, ch chan<- string) {
	params := map[string]string{
		"groupId":              ID,
		"packageCode":          sc.Itinerary.Code,
		"roomIndex":            "1",
		"sailDate":             sc.SailDate,
		"selectedCurrencyCode": "USD",
		"shipCode":             shipCode,
		"startDate":            sc.StartDate,
		"stepSubtypeFlow":      "selectAndContinueSailDate",
		"country":              "USA",
		"landing":              "false",
		"language":             "en",
		"market":               "usa",
		"browser":              "safari",
		"browserVersion":       "17.0.0",
		"screenWidth":          "1680",
		"browserWidth":         "1680",
		"device":               "desktop",
	}

	jsonData := map[string]interface{}{
		"acceptedWipeState":              false,
		"continueConnectedStateroomFlow": false,
		"sailDate":                       sc.SailDate,
		"stateroom":                      stateroom.StateroomClass.ID,
		"stateroomCategoryCode":          nil,
		"stateroomSubType":               nil,
		"bookingCategory":                "DEPOSIT_NOT_REFUNDABLE",
		"index":                          1,
	}

	if stateroom.Price.Value == 0 {
		ch <- fmt.Sprintf("SUCCESS: %s %s %s %s %s", shipName, ID, sc.Itinerary.Code, sc.SailDate, stateroom.StateroomClass.ID)
		return
	}
	if shipName != "Harmony of the Seas" && shipName != "Ovation of the Seas"{
		ch <- fmt.Sprintf("SUCCESS: %s %s %s %s %s", shipName, ID, sc.Itinerary.Code, sc.SailDate, stateroom.StateroomClass.ID)
		return
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		ch <- "error marshaling json"
	}

	client := &http.Client{}

	// Make the HTTP request
	url := "https://www.royalcaribbean.com/mcb/api/booking/step/sailDate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		ch <- "error building request"
		return
	}

	req.URL.RawQuery = paramsToQueryString(params)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")

	req = req.WithContext(ctx)

	var details CruiseDetails
	resp, err := client.Do(req)
	if err != nil {
		ch <- "error downloading"
		return
	}

	defer resp.Body.Close()



	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- "error downloading"
		return
	}

	test := json.Unmarshal(bodyText, &details)
	if test != nil {
		ch <- fmt.Sprintf("ERROR: %s %s %s %s %s", shipName, ID, sc.Itinerary.Code, sc.SailDate, stateroom.StateroomClass.ID)
		return
	}
	for _, srd := range details.StepDetails.StateroomCategoryGroups {
		for _, srcg := range srd.StateroomCategories {
			hc.updateCustomMetrics(
				&customMetric{
					url:                  "test",
					dnsMS:                10,
					connectMS:            10,
					firstbyteMS:          10,
					totalMS:              10,
					status:               1,
					price:                float64(srcg.PriceLockup.TotalPriceNbr),
					cruiseID:             ID,
					itinerary:            sc.Itinerary.Code,
					stateroomClass:       stateroom.StateroomClass.ID,
					dateLabel:            sc.SailDate,
					ship:                 shipName,
					departurePort:        departurePort,
					days:                 strconv.Itoa(days),
					shipCode:             shipCode,
					destinationCode:      shipDestinationCode,
					stateroomDetailClass: srcg.Title,
				},
			)
		}
	}
	ch <- fmt.Sprintf("SUCCESS: %s %s %s %s %s", shipName, ID, sc.Itinerary.Code, sc.SailDate, stateroom.StateroomClass.ID)
}

func (hc *Exporter) updateCustomMetrics(cm *customMetric) {
	// log.Printf("Updating custom metrics: url: %s, connectMS: %.0f, dnsMS: %.0f, firstbyteMS: %.0f, totalMS: %.0f, status: %.0f",
	// 	cm.url,
	// 	cm.connectMS,
	// 	cm.dnsMS,
	// 	cm.firstbyteMS,
	// 	cm.totalMS,
	// 	cm.status,
	// )
	hc.urlDNS.With(prometheus.Labels{
		"url": cm.url,
	}).Set(cm.dnsMS)
	hc.urlConnectTime.With(prometheus.Labels{
		"url": cm.url,
	}).Set(cm.connectMS)
	hc.urlMs.With(prometheus.Labels{
		"url": cm.url,
	}).Set(cm.totalMS)
	hc.urlFirstByte.With(prometheus.Labels{
		"url": cm.url,
	}).Set(cm.firstbyteMS)
	hc.urlStatus.With(prometheus.Labels{
		"url": cm.url,
	}).Set(cm.status)
	hc.royalPrice.With(prometheus.Labels{
		"url":                  cm.url,
		"cruiseid":             cm.cruiseID,
		"itinerary":            cm.itinerary,
		"stateroomclass":       cm.stateroomClass,
		"datelabel":            cm.dateLabel,
		"ship":                 cm.ship,
		"departureport":        cm.departurePort,
		"days":                 cm.days,
		"shipcode":             cm.shipCode,
		"destinationcode":      cm.destinationCode,
		"stateroomdetailclass": cm.stateroomDetailClass,
	}).Set(cm.price)
}

func (hc *Exporter) fetchStats(url string) {

    hc.royalPrice.Reset()

	//var start, connect, dns time.Time

	//var connectMS, dnsMS, firstbyteMS, totalMS, status float64

	trace := &httptrace.ClientTrace{
		// //	DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		// //	DNSDone: func(ddi httptrace.DNSDoneInfo) {
		// 		//dnsMS = float64(time.Since(dns).Milliseconds())
		// 	},

		// 	ConnectStart: func(network, addr string) { connect = time.Now() },
		// 	ConnectDone: func(network, addr string, err error) {
		// 		//connectMS = float64(time.Since(connect).Milliseconds())
		// 	},

		// 	GotFirstResponseByte: func() {
		// 		//firstbyteMS = float64(time.Since(start).Milliseconds())
		// 	},
	}

	count := 20 // Set the number of results per page
	skip := 0  // Start with the first page

	for {
		jsonData := map[string]interface{}{
			"operationName": "cruiseSearch_Cruises",
			"variables": map[string]interface{}{
				"sort": map[string]interface{}{
					"by": "RECOMMENDED",
				},
				"pagination": map[string]interface{}{
					"count": count,
					"skip":  skip,
				},
			},
			"query": "query cruiseSearch_Cruises($filters: String, $qualifiers: String, $sort: CruiseSearchSort, $pagination: CruiseSearchPagination) { cruiseSearch( filters: $filters qualifiers: $qualifiers sort: $sort pagination: $pagination ) { results { cruises { id productViewLink lowestPriceSailing { bookingLink id lowestStateroomClassPrice { price { value __typename } stateroomClass { id __typename } __typename } sailDate startDate endDate taxesAndFees { value __typename } taxesAndFeesIncluded __typename } masterSailing { itinerary { code media { images { path __typename } __typename } days { number type ports { activity arrivalTime departureTime port { code name region media { images { path __typename } __typename } __typename } __typename } __typename } departurePort { code name region __typename } destination { code name __typename } name postTour { days { number type ports { activity arrivalTime departureTime port { code name region __typename } __typename } __typename } duration __typename } preTour { days { number type ports { activity arrivalTime departureTime port { code name region __typename } __typename } __typename } duration __typename } sailingNights ship { code name stateroomClasses { id name content { amenities area code maxCapacity media { images { path meta { description title location __typename } __typename } __typename } superCategory __typename } __typename } media { images { path __typename } __typename } __typename } totalNights type __typename } __typename } sailings { bookingLink id itinerary { code __typename } sailDate startDate endDate stateroomClassPricing { price { value __typename } stateroomClass { id __typename } __typename } __typename } __typename } cruiseRecommendationId total __typename } __typename } }",
		}

		jsonValue, _ := json.Marshal(jsonData)

		// Create an HTTP request with the JSON data and custom User-Agent header.
		req, err := http.NewRequestWithContext(httptrace.WithClientTrace(hc.ctx, trace), "POST", url, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")

		//start = time.Now()
		// Send the HTTP request.
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response:", err)
		}

		var data CruiseSearch
		json.Unmarshal(bodyText, &data)

		for _, s := range data.Data.CruiseSearch.Results.Cruises {
			for _, sc := range s.Sailings {
				ch := make(chan string, len(sc.StateroomClassPricing))
				defer close(ch)

				test, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				for _, stateroom := range sc.StateroomClassPricing {
					go hc.fetchDetails(
						test,
						s.MasterSailing.Itinerary.Ship.Name, 
						s.MasterSailing.Itinerary.Destination.Code, 
						s.MasterSailing.Itinerary.DeparturePort.Name, 
						s.MasterSailing.Itinerary.TotalNights, 
						s.ID, s.MasterSailing.Itinerary.Ship.Code, 
						sc, 
						stateroom, 
						ch)
				}
				for i := 0; i < len(sc.StateroomClassPricing); i++ {
					select {
					case result := <-ch:
						fmt.Println(result)
					case <-test.Done():
						fmt.Println("Request timed out")
					}
				}
			}
		}
		if skip < (data.Data.CruiseSearch.Results.Total - 20) {
			skip = skip + 20
		} else {
			break
		}
		log.Printf("pulled down %d skipping the first %d of %d total", count, skip, data.Data.CruiseSearch.Results.Total)
	}
}

func (hc *Exporter) StartCollector() {
	ticker := time.NewTicker(hc.healthcheck_invertval)
	log.Println("starting exporter")
	for _, u := range hc.urls {
		hc.fetchStats(u)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, u := range hc.urls {
					hc.fetchStats(u)
				}
			case <-hc.ctx.Done():
				log.Println("Gracefully stopping exporter")
				return
			}
		}
	}()
}
