#!/bin/sh

# Create the output directory if it doesn't exist
mkdir -p src/proto

# Generate JavaScript code
protoc \
  --js_out=import_style=commonjs,binary:src/proto \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:src/proto \
  -I../proto \
  ../proto/nexus.proto

# Create ES module wrapper for nexus_pb.js
echo 'import "./nexus_pb.js";

const proto = window.proto;
export { proto };' > src/proto/nexus_pb_es.js

# Create ES module wrapper for nexus service
echo 'import { NexusServiceClient } from "./NexusServiceClientPb";

const NexusService = {
  serviceName: "nexus.NexusService",
  GetChildren: {
    methodName: "GetChildren",
    service: NexusServiceClient,
  },
  GetNode: {
    methodName: "GetNode",
    service: NexusServiceClient,
  },
  PublishValue: {
    methodName: "PublishValue",
    service: NexusServiceClient,
  },
  PublishEventStream: {
    methodName: "PublishEventStream",
    service: NexusServiceClient,
  },
  PublishIndividualFile: {
    methodName: "PublishIndividualFile",
    service: NexusServiceClient,
  },
  PublishDirectory: {
    methodName: "PublishDirectory",
    service: NexusServiceClient,
  },
  PublishDatabaseTable: {
    methodName: "PublishDatabaseTable",
    service: NexusServiceClient,
  },
  DeletePath: {
    methodName: "DeletePath",
    service: NexusServiceClient,
  }
};

export { NexusService };' > src/proto/nexus_pb_service_es.js