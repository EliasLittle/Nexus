import { useMemo } from 'react';
import { grpc } from '@improbable-eng/grpc-web';
import { proto } from '../proto/nexus_pb_es';
import { NexusService } from '../proto/nexus_pb_service_es';

const DEFAULT_CONNECTION = 'http://localhost:8080';

const createEventStream = (topic, server = 'localhost:9092') => {
  const stream = new proto.nexus.EventStream();
  stream.setServer(server);
  stream.setTopic(topic);
  return stream;
};

const createIndividualFile = (filePath) => {
  const file = new proto.nexus.IndividualFile();
  const fileType = filePath.split('.').pop() || '';
  file.setFiletype(fileType);
  file.setFilepath(filePath);
  file.setColumnnamesList([]);
  return file;
};

const createDirectory = (directoryPath, fileType, fileCount) => {
  const dir = new proto.nexus.Directory();
  dir.setFiletype(fileType);
  dir.setDirectorypath(directoryPath);
  dir.setFilecount(fileCount);
  return dir;
};

const createDatabaseTable = (dbType, host, port, dbName, tableName) => {
  const table = new proto.nexus.DatabaseTable();
  table.setDbtype(dbType);
  table.setHost(host);
  table.setPort(port);
  table.setDbname(dbName);
  table.setTablename(tableName);
  return table;
};

export const useNexusClient = (host = DEFAULT_CONNECTION) => {
  console.log('Host:', host);
  return useMemo(() => {
    const client = {
      getChildren: (path) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.GetChildrenRequest();
          request.setPath(path);

          grpc.unary(NexusService.GetChildren, {
            request,
            host,
            onEnd: ({ status, statusMessage, message }) => {
              if (status === grpc.Code.OK && message) {
                resolve(message.getChildrenList());
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      getFull: (path) => {
        console.log('Getting full path:', path);
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.GetPathRequest();
          request.setPath(path);

          console.log('Request:', request);

          grpc.unary(NexusService.GetNode, {
            request,
            host,
            onEnd: ({ status, statusMessage, message }) => {
              if (status === grpc.Code.OK && message) {
                const valueType = message.getValuetype();
                let value = null;

                try {
                  switch (valueType) {
                    case 'InternalNode':
                      value = null;
                      break;
                    case 'StringValue':
                      value = message.getStringvalue();
                      break;
                    case 'IntValue':
                      value = message.getIntvalue();
                      break;
                    case 'FloatValue':
                      value = message.getFloatvalue();
                      break;
                    case 'DatabaseTable':
                      value = message.getDatabasetable();
                      break;
                    case 'Directory':
                      value = message.getDirectory();
                      break;
                    case 'IndividualFile':
                      value = message.getIndividualfile();
                      break;
                    case 'EventStream':
                      value = message.getEventstream();
                      break;
                    default:
                      throw new Error(`Unknown value type: ${valueType}`);
                  }
                  resolve({ value, type: valueType });
                } catch (error) {
                  reject(error);
                }
              } else {
                reject(new Error(statusMessage || 'Unknown error occurred'));
              }
            },
          });
        });
      },

      publishValue: (path, value) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.PublishValueRequest();
          request.setPath(path);

          // Determine value type and set accordingly
          if (typeof value === 'string') {
            const stringValue = new proto.nexus.StringValue();
            stringValue.setValue(value);
            request.setStringvalue(stringValue);
          } else if (Number.isInteger(value)) {
            const intValue = new proto.nexus.IntValue();
            intValue.setValue(value);
            request.setIntvalue(intValue);
          } else if (typeof value === 'number') {
            const floatValue = new proto.nexus.FloatValue();
            floatValue.setValue(value);
            request.setFloatvalue(floatValue);
          } else {
            reject(new Error('Unsupported value type'));
            return;
          }

          grpc.unary(NexusService.PublishValue, {
            request,
            host,
            onEnd: ({ status, statusMessage }) => {
              if (status === grpc.Code.OK) {
                resolve();
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      publishEventStream: (path, config) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.PublishEventStreamRequest();
          request.setPath(path);
          request.setEventstream(createEventStream(config.topic, config.server));

          grpc.unary(NexusService.PublishEventStream, {
            request,
            host,
            onEnd: ({ status, statusMessage }) => {
              if (status === grpc.Code.OK) {
                resolve();
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      publishIndividualFile: (path, config) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.PublishIndividualFileRequest();
          request.setPath(path);
          request.setIndividualfile(createIndividualFile(config.filePath));

          grpc.unary(NexusService.PublishIndividualFile, {
            request,
            host,
            onEnd: ({ status, statusMessage }) => {
              if (status === grpc.Code.OK) {
                resolve();
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      publishDirectory: (path, config) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.PublishDirectoryRequest();
          request.setPath(path);
          request.setDirectory(createDirectory(config.directoryPath, config.fileType, config.fileCount));

          grpc.unary(NexusService.PublishDirectory, {
            request,
            host,
            onEnd: ({ status, statusMessage }) => {
              if (status === grpc.Code.OK) {
                resolve();
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      publishDatabaseTable: (path, config) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.PublishDatabaseTableRequest();
          request.setPath(path);
          request.setDatabasetable(createDatabaseTable(
            config.dbType,
            config.host,
            config.port,
            config.dbName,
            config.tableName
          ));

          grpc.unary(NexusService.PublishDatabaseTable, {
            request,
            host,
            onEnd: ({ status, statusMessage }) => {
              if (status === grpc.Code.OK) {
                resolve();
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },

      delete: (path) => {
        return new Promise((resolve, reject) => {
          const request = new proto.nexus.DeletePathRequest();
          request.setPath(path);

          grpc.unary(NexusService.DeletePath, {
            request,
            host,
            onEnd: ({ status, statusMessage, message }) => {
              if (status === grpc.Code.OK) {
                if (message && !message.getSuccess()) {
                  reject(new Error(message.getError()));
                } else {
                  resolve();
                }
              } else {
                reject(new Error(statusMessage));
              }
            },
          });
        });
      },
    };

    return client;
  }, [host]);
}; 