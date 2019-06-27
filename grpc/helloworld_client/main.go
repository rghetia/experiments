// Copyright 2017, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"time"

	"context"
        "contrib.go.opencensus.io/exporter/stackdriver"
        "contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	//"go.opencensus.io/examples/exporter"
	pb "github.com/rghetia/experiments/grpc/proto"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Register stats and trace exporters to export
	// the collected data.
	// view.RegisterExporter(&exporter.PrintExporter{})
        exporter, err := stackdriver.NewExporter(stackdriver.Options{
                ProjectID:         os.Getenv("GCP_PROJECT_ID"), // Google Cloud Console project ID for stackdriver.
                MonitoredResource: monitoredresource.Autodetect(),
        })
        if err != nil {
                log.Fatal(err)
        }
        view.RegisterExporter(exporter)

        // Set reporting period to report data at 60 seconds.
        // The recommended reporting period by Stackdriver Monitoring is >= 1 minute:
        // https://cloud.google.com/monitoring/custom-metrics/creating-metrics#writing-ts.
        view.SetReportingPeriod(60 * time.Second)


	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	// Set up a connection to the server with the OpenCensus
	// stats handler to enable stats and tracing.
	conn, err := grpc.Dial(address, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
//	view.SetReportingPeriod(time.Second)
	for {
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("Could not greet: %v", err)
		} else {
			//log.Printf("Greeting: %s", r.Message)
		}
		time.Sleep(2 * time.Second)
	}
}
