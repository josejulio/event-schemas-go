/*
 Copyright 2021 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

// Adapted from https://github.com/cloudevents/sdk-go/blob/a7187527ab3278128c1b2a8fe9856d49ecddf25d/samples/kafka/sender-receiver/main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RedHatInsights/event-schemas-go/apps/advisor/v1"
	"github.com/Shopify/sarama"
	"log"
	"sync/atomic"
	"time"

	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const (
	count = 1
)

func setup(topic string, structured bool) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	// With NewProtocol you can use the same client both to send and receive.
	protocol, err := kafka_sarama.NewProtocol([]string{"127.0.0.1:9092"}, saramaConfig, topic, topic)
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	defer protocol.Close(context.Background())

	c, err := cloudevents.NewClient(protocol, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Create a done channel to block until we've received (count) messages
	done := make(chan struct{})

	// Start the receiver
	go func() {
		log.Printf("will listen consuming topic %s\n", topic)
		var recvCount int32
		err = c.StartReceiver(context.TODO(), func(ctx context.Context, event cloudevents.Event) {
			receive(ctx, event)
			if atomic.AddInt32(&recvCount, 1) == count {
				done <- struct{}{}
			}
		})
		if err != nil {
			log.Fatalf("failed to start receiver: %s", err)
		} else {
			log.Printf("receiver stopped\n")
		}
	}()

	// Start sending the events
	for i := 0; i < count; i++ {
		e := cloudevents.NewEvent()
		e.SetType("com.redhat.console.advisor.new-recommendations")
		e.SetSource("urn:redhat:source:console:app:advisor")
		e.SetID("urn:redhat:console:event:5864ac25-4c52-4c87-bd28-9909a4fa3187")
		e.SetSubject("urn:redhat:subject:console:rhel:08e8ec2b-6a79-4f1d-bea4-a438da139493")
		e.SetTime(time.Now())
		e.SetExtension("redhatorgid", "org123")
		e.SetExtension("redhatconsolebundle", "rhel")
		e.SetDataSchema("https://console.redhat.com/api/schemas/apps/advisor/v1/advisor-recommendations.json")

		hostname := "rhel8desktop"
		hostUrl := "https://console.redhat.com/insights/inventory/08e8ec2b-6a79-4f1d-bea4-a438da139493"
		rhelVersion := "8.3"
		data := advisor.AdvisorRecommendations{
			System: advisor.RHELSystem{
				DisplayName: &hostname,
				HostURL:     &hostUrl,
				Hostname:    &hostname,
				InventoryID: "08e8ec2b-6a79-4f1d-bea4-a438da139493",
				RHELVersion: &rhelVersion,
				Tags: []advisor.RHELSystemTag{
					{
						Namespace: "insights-client",
						Key:       "Environment",
						Value:     "Production",
					},
				},
			},
			AdvisorRecommendations: []advisor.AdvisorRecommendation{
				{
					PublishDate:     "2021-03-13T18:44:00+00:00",
					RebootRequired:  false,
					RuleDescription: "System is not able to get the latest recommendations and may miss bug fixes when the Insights Client Core egg file is outdated",
					RuleID:          "insights_core_egg_not_up2date|INSIGHTS_CORE_EGG_NOT_UP2DATE",
					RuleURL:         "https://console.redhat.com/insights/advisor/recommendations/insights_core_egg_not_up2date|INSIGHTS_CORE_EGG_NOT_UP2DATE/",
					TotalRisk:       "2",
				},
			},
		}
		err := e.SetData(cloudevents.ApplicationJSON, data)
		if err != nil {
			return
		}

		var ctx context.Context
		ctx = context.Background()
		if structured {
			ctx = cloudevents.WithEncodingStructured(ctx)
		} else {
			ctx = cloudevents.WithEncodingBinary(ctx)
		}
		if result := c.Send(ctx, e); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send: %v", result)
		} else {
			log.Printf("sent: %d, accepted: %t", i, cloudevents.IsACK(result))
		}
	}

	<-done
}

func main() {
	setup("test-topic", false)
	setup("test-topic-structured", true)
}

func receive(ctx context.Context, event cloudevents.Event) {
	var err error
	serialized, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Fatalf("Unable to serialize to JSON\n")
	}
	fmt.Printf("CloudEvent in structured format: %s\n", serialized)
	result := advisor.AdvisorRecommendations{}
	err = event.DataAs(&result)
	if err != nil {
		log.Fatalf("Unable to parse data\n")
	}
	var orgId string
	err = event.ExtensionAs("redhatorgid", &orgId)
	if err != nil {
		log.Fatalf("Unable to extract orgId\n")
	}
	fmt.Printf("Message received for orgId: %s\n", orgId)
	fmt.Printf("DisplayName: %s\n", *result.System.DisplayName)
}
