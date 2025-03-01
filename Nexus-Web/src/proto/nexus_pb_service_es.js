import { NexusServiceClient } from "./NexusServiceClientPb";

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

export { NexusService };
