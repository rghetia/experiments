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

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/helloworld.proto

package main

import (
	"log"
	"math/rand"
        "os"
	"net"
	"net/http"
	"time"

	"context"
        "contrib.go.opencensus.io/exporter/stackdriver"
        "contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
//	"go.opencensus.io/examples/exporter"
	pb "github.com/rghetia/experiments/grpc/proto"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
)

const port = ":50051"

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	ctx, span := trace.StartSpan(ctx, "sleep")
	time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
	span.End()
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// Start z-Pages server.
	go func() {
		mux := http.NewServeMux()
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
	}()

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
        trace.RegisterExporter(exporter)

        // Set reporting period to report data at 60 seconds.
        // The recommended reporting period by Stackdriver Monitoring is >= 1 minute:
        // https://cloud.google.com/monitoring/custom-metrics/creating-metrics#writing-ts.
        trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

        view.SetReportingPeriod(60 * time.Second)


	// Register the views to collect server request count.
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Set up a new server with the OpenCensus
	// stats handler to enable stats and tracing.
	s := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	pb.RegisterGreeterServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
