import * as jspb from 'google-protobuf'



export class RegisterEventStreamRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): RegisterEventStreamRequest;

  getEventStream(): EventStream | undefined;
  setEventStream(value?: EventStream): RegisterEventStreamRequest;
  hasEventStream(): boolean;
  clearEventStream(): RegisterEventStreamRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterEventStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterEventStreamRequest): RegisterEventStreamRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterEventStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterEventStreamRequest;
  static deserializeBinaryFromReader(message: RegisterEventStreamRequest, reader: jspb.BinaryReader): RegisterEventStreamRequest;
}

export namespace RegisterEventStreamRequest {
  export type AsObject = {
    path: string,
    eventStream?: EventStream.AsObject,
  }
}

export class RegisterEventStreamResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): RegisterEventStreamResponse;

  getError(): string;
  setError(value: string): RegisterEventStreamResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterEventStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterEventStreamResponse): RegisterEventStreamResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterEventStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterEventStreamResponse;
  static deserializeBinaryFromReader(message: RegisterEventStreamResponse, reader: jspb.BinaryReader): RegisterEventStreamResponse;
}

export namespace RegisterEventStreamResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class RegisterFileRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): RegisterFileRequest;

  getIndividualFile(): IndividualFile | undefined;
  setIndividualFile(value?: IndividualFile): RegisterFileRequest;
  hasIndividualFile(): boolean;
  clearIndividualFile(): RegisterFileRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterFileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterFileRequest): RegisterFileRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterFileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterFileRequest;
  static deserializeBinaryFromReader(message: RegisterFileRequest, reader: jspb.BinaryReader): RegisterFileRequest;
}

export namespace RegisterFileRequest {
  export type AsObject = {
    path: string,
    individualFile?: IndividualFile.AsObject,
  }
}

export class RegisterFileResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): RegisterFileResponse;

  getError(): string;
  setError(value: string): RegisterFileResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterFileResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterFileResponse): RegisterFileResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterFileResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterFileResponse;
  static deserializeBinaryFromReader(message: RegisterFileResponse, reader: jspb.BinaryReader): RegisterFileResponse;
}

export namespace RegisterFileResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class RegisterDirectoryRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): RegisterDirectoryRequest;

  getDirectory(): Directory | undefined;
  setDirectory(value?: Directory): RegisterDirectoryRequest;
  hasDirectory(): boolean;
  clearDirectory(): RegisterDirectoryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDirectoryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDirectoryRequest): RegisterDirectoryRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterDirectoryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDirectoryRequest;
  static deserializeBinaryFromReader(message: RegisterDirectoryRequest, reader: jspb.BinaryReader): RegisterDirectoryRequest;
}

export namespace RegisterDirectoryRequest {
  export type AsObject = {
    path: string,
    directory?: Directory.AsObject,
  }
}

export class RegisterDirectoryResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): RegisterDirectoryResponse;

  getError(): string;
  setError(value: string): RegisterDirectoryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDirectoryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDirectoryResponse): RegisterDirectoryResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterDirectoryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDirectoryResponse;
  static deserializeBinaryFromReader(message: RegisterDirectoryResponse, reader: jspb.BinaryReader): RegisterDirectoryResponse;
}

export namespace RegisterDirectoryResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class RegisterDatabaseTableRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): RegisterDatabaseTableRequest;

  getDatabaseTable(): DatabaseTable | undefined;
  setDatabaseTable(value?: DatabaseTable): RegisterDatabaseTableRequest;
  hasDatabaseTable(): boolean;
  clearDatabaseTable(): RegisterDatabaseTableRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDatabaseTableRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDatabaseTableRequest): RegisterDatabaseTableRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterDatabaseTableRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDatabaseTableRequest;
  static deserializeBinaryFromReader(message: RegisterDatabaseTableRequest, reader: jspb.BinaryReader): RegisterDatabaseTableRequest;
}

export namespace RegisterDatabaseTableRequest {
  export type AsObject = {
    path: string,
    databaseTable?: DatabaseTable.AsObject,
  }
}

export class RegisterDatabaseTableResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): RegisterDatabaseTableResponse;

  getError(): string;
  setError(value: string): RegisterDatabaseTableResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDatabaseTableResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDatabaseTableResponse): RegisterDatabaseTableResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterDatabaseTableResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDatabaseTableResponse;
  static deserializeBinaryFromReader(message: RegisterDatabaseTableResponse, reader: jspb.BinaryReader): RegisterDatabaseTableResponse;
}

export namespace RegisterDatabaseTableResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class StoreValueRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): StoreValueRequest;

  getStringValue(): StringValue | undefined;
  setStringValue(value?: StringValue): StoreValueRequest;
  hasStringValue(): boolean;
  clearStringValue(): StoreValueRequest;

  getIntValue(): IntValue | undefined;
  setIntValue(value?: IntValue): StoreValueRequest;
  hasIntValue(): boolean;
  clearIntValue(): StoreValueRequest;

  getFloatValue(): FloatValue | undefined;
  setFloatValue(value?: FloatValue): StoreValueRequest;
  hasFloatValue(): boolean;
  clearFloatValue(): StoreValueRequest;

  getValueCase(): StoreValueRequest.ValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StoreValueRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StoreValueRequest): StoreValueRequest.AsObject;
  static serializeBinaryToWriter(message: StoreValueRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StoreValueRequest;
  static deserializeBinaryFromReader(message: StoreValueRequest, reader: jspb.BinaryReader): StoreValueRequest;
}

export namespace StoreValueRequest {
  export type AsObject = {
    path: string,
    stringValue?: StringValue.AsObject,
    intValue?: IntValue.AsObject,
    floatValue?: FloatValue.AsObject,
  }

  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    STRING_VALUE = 2,
    INT_VALUE = 3,
    FLOAT_VALUE = 4,
  }
}

export class StoreValueResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): StoreValueResponse;

  getError(): string;
  setError(value: string): StoreValueResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StoreValueResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StoreValueResponse): StoreValueResponse.AsObject;
  static serializeBinaryToWriter(message: StoreValueResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StoreValueResponse;
  static deserializeBinaryFromReader(message: StoreValueResponse, reader: jspb.BinaryReader): StoreValueResponse;
}

export namespace StoreValueResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class DeletePathRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): DeletePathRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePathRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePathRequest): DeletePathRequest.AsObject;
  static serializeBinaryToWriter(message: DeletePathRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePathRequest;
  static deserializeBinaryFromReader(message: DeletePathRequest, reader: jspb.BinaryReader): DeletePathRequest;
}

export namespace DeletePathRequest {
  export type AsObject = {
    path: string,
  }
}

export class DeletePathResponse extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): DeletePathResponse;

  getError(): string;
  setError(value: string): DeletePathResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeletePathResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeletePathResponse): DeletePathResponse.AsObject;
  static serializeBinaryToWriter(message: DeletePathResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeletePathResponse;
  static deserializeBinaryFromReader(message: DeletePathResponse, reader: jspb.BinaryReader): DeletePathResponse;
}

export namespace DeletePathResponse {
  export type AsObject = {
    success: boolean,
    error: string,
  }
}

export class EventStream extends jspb.Message {
  getServer(): string;
  setServer(value: string): EventStream;

  getTopic(): string;
  setTopic(value: string): EventStream;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EventStream.AsObject;
  static toObject(includeInstance: boolean, msg: EventStream): EventStream.AsObject;
  static serializeBinaryToWriter(message: EventStream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EventStream;
  static deserializeBinaryFromReader(message: EventStream, reader: jspb.BinaryReader): EventStream;
}

export namespace EventStream {
  export type AsObject = {
    server: string,
    topic: string,
  }
}

export class Dataset extends jspb.Message {
  getIndividualFile(): IndividualFile | undefined;
  setIndividualFile(value?: IndividualFile): Dataset;
  hasIndividualFile(): boolean;
  clearIndividualFile(): Dataset;

  getDirectory(): Directory | undefined;
  setDirectory(value?: Directory): Dataset;
  hasDirectory(): boolean;
  clearDirectory(): Dataset;

  getDatabaseTable(): DatabaseTable | undefined;
  setDatabaseTable(value?: DatabaseTable): Dataset;
  hasDatabaseTable(): boolean;
  clearDatabaseTable(): Dataset;

  getDatasetCase(): Dataset.DatasetCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Dataset.AsObject;
  static toObject(includeInstance: boolean, msg: Dataset): Dataset.AsObject;
  static serializeBinaryToWriter(message: Dataset, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Dataset;
  static deserializeBinaryFromReader(message: Dataset, reader: jspb.BinaryReader): Dataset;
}

export namespace Dataset {
  export type AsObject = {
    individualFile?: IndividualFile.AsObject,
    directory?: Directory.AsObject,
    databaseTable?: DatabaseTable.AsObject,
  }

  export enum DatasetCase { 
    DATASET_NOT_SET = 0,
    INDIVIDUAL_FILE = 1,
    DIRECTORY = 2,
    DATABASE_TABLE = 3,
  }
}

export class IndividualFile extends jspb.Message {
  getFileType(): string;
  setFileType(value: string): IndividualFile;

  getFilePath(): string;
  setFilePath(value: string): IndividualFile;

  getColumnNamesList(): Array<string>;
  setColumnNamesList(value: Array<string>): IndividualFile;
  clearColumnNamesList(): IndividualFile;
  addColumnNames(value: string, index?: number): IndividualFile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IndividualFile.AsObject;
  static toObject(includeInstance: boolean, msg: IndividualFile): IndividualFile.AsObject;
  static serializeBinaryToWriter(message: IndividualFile, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IndividualFile;
  static deserializeBinaryFromReader(message: IndividualFile, reader: jspb.BinaryReader): IndividualFile;
}

export namespace IndividualFile {
  export type AsObject = {
    fileType: string,
    filePath: string,
    columnNamesList: Array<string>,
  }
}

export class Directory extends jspb.Message {
  getFileType(): string;
  setFileType(value: string): Directory;

  getDirectoryPath(): string;
  setDirectoryPath(value: string): Directory;

  getFileCount(): number;
  setFileCount(value: number): Directory;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Directory.AsObject;
  static toObject(includeInstance: boolean, msg: Directory): Directory.AsObject;
  static serializeBinaryToWriter(message: Directory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Directory;
  static deserializeBinaryFromReader(message: Directory, reader: jspb.BinaryReader): Directory;
}

export namespace Directory {
  export type AsObject = {
    fileType: string,
    directoryPath: string,
    fileCount: number,
  }
}

export class DatabaseTable extends jspb.Message {
  getDbType(): string;
  setDbType(value: string): DatabaseTable;

  getHost(): string;
  setHost(value: string): DatabaseTable;

  getPort(): number;
  setPort(value: number): DatabaseTable;

  getDbName(): string;
  setDbName(value: string): DatabaseTable;

  getTableName(): string;
  setTableName(value: string): DatabaseTable;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatabaseTable.AsObject;
  static toObject(includeInstance: boolean, msg: DatabaseTable): DatabaseTable.AsObject;
  static serializeBinaryToWriter(message: DatabaseTable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DatabaseTable;
  static deserializeBinaryFromReader(message: DatabaseTable, reader: jspb.BinaryReader): DatabaseTable;
}

export namespace DatabaseTable {
  export type AsObject = {
    dbType: string,
    host: string,
    port: number,
    dbName: string,
    tableName: string,
  }
}

export class StringValue extends jspb.Message {
  getValue(): string;
  setValue(value: string): StringValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StringValue.AsObject;
  static toObject(includeInstance: boolean, msg: StringValue): StringValue.AsObject;
  static serializeBinaryToWriter(message: StringValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StringValue;
  static deserializeBinaryFromReader(message: StringValue, reader: jspb.BinaryReader): StringValue;
}

export namespace StringValue {
  export type AsObject = {
    value: string,
  }
}

export class IntValue extends jspb.Message {
  getValue(): number;
  setValue(value: number): IntValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IntValue.AsObject;
  static toObject(includeInstance: boolean, msg: IntValue): IntValue.AsObject;
  static serializeBinaryToWriter(message: IntValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IntValue;
  static deserializeBinaryFromReader(message: IntValue, reader: jspb.BinaryReader): IntValue;
}

export namespace IntValue {
  export type AsObject = {
    value: number,
  }
}

export class FloatValue extends jspb.Message {
  getValue(): number;
  setValue(value: number): FloatValue;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FloatValue.AsObject;
  static toObject(includeInstance: boolean, msg: FloatValue): FloatValue.AsObject;
  static serializeBinaryToWriter(message: FloatValue, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FloatValue;
  static deserializeBinaryFromReader(message: FloatValue, reader: jspb.BinaryReader): FloatValue;
}

export namespace FloatValue {
  export type AsObject = {
    value: number,
  }
}

export class SubscribeRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): SubscribeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeRequest): SubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: SubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeRequest;
  static deserializeBinaryFromReader(message: SubscribeRequest, reader: jspb.BinaryReader): SubscribeRequest;
}

export namespace SubscribeRequest {
  export type AsObject = {
    path: string,
  }
}

export class Event extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    data: Uint8Array | string,
  }
}

export class GetPathRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): GetPathRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPathRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPathRequest): GetPathRequest.AsObject;
  static serializeBinaryToWriter(message: GetPathRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPathRequest;
  static deserializeBinaryFromReader(message: GetPathRequest, reader: jspb.BinaryReader): GetPathRequest;
}

export namespace GetPathRequest {
  export type AsObject = {
    path: string,
  }
}

export class GetNodeResponse extends jspb.Message {
  getStringValue(): StringValue | undefined;
  setStringValue(value?: StringValue): GetNodeResponse;
  hasStringValue(): boolean;
  clearStringValue(): GetNodeResponse;

  getIntValue(): IntValue | undefined;
  setIntValue(value?: IntValue): GetNodeResponse;
  hasIntValue(): boolean;
  clearIntValue(): GetNodeResponse;

  getFloatValue(): FloatValue | undefined;
  setFloatValue(value?: FloatValue): GetNodeResponse;
  hasFloatValue(): boolean;
  clearFloatValue(): GetNodeResponse;

  getIndividualFile(): IndividualFile | undefined;
  setIndividualFile(value?: IndividualFile): GetNodeResponse;
  hasIndividualFile(): boolean;
  clearIndividualFile(): GetNodeResponse;

  getDirectory(): Directory | undefined;
  setDirectory(value?: Directory): GetNodeResponse;
  hasDirectory(): boolean;
  clearDirectory(): GetNodeResponse;

  getDatabaseTable(): DatabaseTable | undefined;
  setDatabaseTable(value?: DatabaseTable): GetNodeResponse;
  hasDatabaseTable(): boolean;
  clearDatabaseTable(): GetNodeResponse;

  getEventStream(): EventStream | undefined;
  setEventStream(value?: EventStream): GetNodeResponse;
  hasEventStream(): boolean;
  clearEventStream(): GetNodeResponse;

  getValueType(): string;
  setValueType(value: string): GetNodeResponse;

  getIsEndOfPath(): boolean;
  setIsEndOfPath(value: boolean): GetNodeResponse;

  getError(): string;
  setError(value: string): GetNodeResponse;

  getValueCase(): GetNodeResponse.ValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNodeResponse): GetNodeResponse.AsObject;
  static serializeBinaryToWriter(message: GetNodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNodeResponse;
  static deserializeBinaryFromReader(message: GetNodeResponse, reader: jspb.BinaryReader): GetNodeResponse;
}

export namespace GetNodeResponse {
  export type AsObject = {
    stringValue?: StringValue.AsObject,
    intValue?: IntValue.AsObject,
    floatValue?: FloatValue.AsObject,
    individualFile?: IndividualFile.AsObject,
    directory?: Directory.AsObject,
    databaseTable?: DatabaseTable.AsObject,
    eventStream?: EventStream.AsObject,
    valueType: string,
    isEndOfPath: boolean,
    error: string,
  }

  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    STRING_VALUE = 1,
    INT_VALUE = 2,
    FLOAT_VALUE = 3,
    INDIVIDUAL_FILE = 4,
    DIRECTORY = 5,
    DATABASE_TABLE = 6,
    EVENT_STREAM = 7,
  }
}

export class AccessInfo extends jspb.Message {
  getFile(): FileInfo | undefined;
  setFile(value?: FileInfo): AccessInfo;
  hasFile(): boolean;
  clearFile(): AccessInfo;

  getDatabase(): DatabaseInfo | undefined;
  setDatabase(value?: DatabaseInfo): AccessInfo;
  hasDatabase(): boolean;
  clearDatabase(): AccessInfo;

  getInfoCase(): AccessInfo.InfoCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AccessInfo.AsObject;
  static toObject(includeInstance: boolean, msg: AccessInfo): AccessInfo.AsObject;
  static serializeBinaryToWriter(message: AccessInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AccessInfo;
  static deserializeBinaryFromReader(message: AccessInfo, reader: jspb.BinaryReader): AccessInfo;
}

export namespace AccessInfo {
  export type AsObject = {
    file?: FileInfo.AsObject,
    database?: DatabaseInfo.AsObject,
  }

  export enum InfoCase { 
    INFO_NOT_SET = 0,
    FILE = 1,
    DATABASE = 2,
  }
}

export class FileInfo extends jspb.Message {
  getFilepath(): string;
  setFilepath(value: string): FileInfo;

  getFormat(): string;
  setFormat(value: string): FileInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileInfo.AsObject;
  static toObject(includeInstance: boolean, msg: FileInfo): FileInfo.AsObject;
  static serializeBinaryToWriter(message: FileInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileInfo;
  static deserializeBinaryFromReader(message: FileInfo, reader: jspb.BinaryReader): FileInfo;
}

export namespace FileInfo {
  export type AsObject = {
    filepath: string,
    format: string,
  }
}

export class DatabaseInfo extends jspb.Message {
  getConnectionString(): string;
  setConnectionString(value: string): DatabaseInfo;

  getTable(): string;
  setTable(value: string): DatabaseInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DatabaseInfo.AsObject;
  static toObject(includeInstance: boolean, msg: DatabaseInfo): DatabaseInfo.AsObject;
  static serializeBinaryToWriter(message: DatabaseInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DatabaseInfo;
  static deserializeBinaryFromReader(message: DatabaseInfo, reader: jspb.BinaryReader): DatabaseInfo;
}

export namespace DatabaseInfo {
  export type AsObject = {
    connectionString: string,
    table: string,
  }
}

export class GetChildrenRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): GetChildrenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetChildrenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetChildrenRequest): GetChildrenRequest.AsObject;
  static serializeBinaryToWriter(message: GetChildrenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetChildrenRequest;
  static deserializeBinaryFromReader(message: GetChildrenRequest, reader: jspb.BinaryReader): GetChildrenRequest;
}

export namespace GetChildrenRequest {
  export type AsObject = {
    path: string,
  }
}

export class GetChildrenResponse extends jspb.Message {
  getChildrenList(): Array<ChildInfo>;
  setChildrenList(value: Array<ChildInfo>): GetChildrenResponse;
  clearChildrenList(): GetChildrenResponse;
  addChildren(value?: ChildInfo, index?: number): ChildInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetChildrenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetChildrenResponse): GetChildrenResponse.AsObject;
  static serializeBinaryToWriter(message: GetChildrenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetChildrenResponse;
  static deserializeBinaryFromReader(message: GetChildrenResponse, reader: jspb.BinaryReader): GetChildrenResponse;
}

export namespace GetChildrenResponse {
  export type AsObject = {
    childrenList: Array<ChildInfo.AsObject>,
  }
}

export class ChildInfo extends jspb.Message {
  getName(): string;
  setName(value: string): ChildInfo;

  getType(): string;
  setType(value: string): ChildInfo;

  getNumchildren(): number;
  setNumchildren(value: number): ChildInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChildInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ChildInfo): ChildInfo.AsObject;
  static serializeBinaryToWriter(message: ChildInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChildInfo;
  static deserializeBinaryFromReader(message: ChildInfo, reader: jspb.BinaryReader): ChildInfo;
}

export namespace ChildInfo {
  export type AsObject = {
    name: string,
    type: string,
    numchildren: number,
  }
}

export class GetPathTypeResponse extends jspb.Message {
  getPathType(): string;
  setPathType(value: string): GetPathTypeResponse;

  getError(): string;
  setError(value: string): GetPathTypeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPathTypeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPathTypeResponse): GetPathTypeResponse.AsObject;
  static serializeBinaryToWriter(message: GetPathTypeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPathTypeResponse;
  static deserializeBinaryFromReader(message: GetPathTypeResponse, reader: jspb.BinaryReader): GetPathTypeResponse;
}

export namespace GetPathTypeResponse {
  export type AsObject = {
    pathType: string,
    error: string,
  }
}

