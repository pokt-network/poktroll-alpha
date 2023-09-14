/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Application } from "./application";
import { Params } from "./params";

export const protobufPackage = "poktroll.application";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetApplicationRequest {
  address: string;
}

export interface QueryGetApplicationResponse {
  application: Application | undefined;
}

export interface QueryAllApplicationRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllApplicationResponse {
  application: Application[];
  pagination: PageResponse | undefined;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryGetApplicationRequest(): QueryGetApplicationRequest {
  return { address: "" };
}

export const QueryGetApplicationRequest = {
  encode(message: QueryGetApplicationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetApplicationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetApplicationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetApplicationRequest {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryGetApplicationRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetApplicationRequest>, I>>(object: I): QueryGetApplicationRequest {
    const message = createBaseQueryGetApplicationRequest();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGetApplicationResponse(): QueryGetApplicationResponse {
  return { application: undefined };
}

export const QueryGetApplicationResponse = {
  encode(message: QueryGetApplicationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.application !== undefined) {
      Application.encode(message.application, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetApplicationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetApplicationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.application = Application.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetApplicationResponse {
    return { application: isSet(object.application) ? Application.fromJSON(object.application) : undefined };
  },

  toJSON(message: QueryGetApplicationResponse): unknown {
    const obj: any = {};
    message.application !== undefined
      && (obj.application = message.application ? Application.toJSON(message.application) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetApplicationResponse>, I>>(object: I): QueryGetApplicationResponse {
    const message = createBaseQueryGetApplicationResponse();
    message.application = (object.application !== undefined && object.application !== null)
      ? Application.fromPartial(object.application)
      : undefined;
    return message;
  },
};

function createBaseQueryAllApplicationRequest(): QueryAllApplicationRequest {
  return { pagination: undefined };
}

export const QueryAllApplicationRequest = {
  encode(message: QueryAllApplicationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllApplicationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllApplicationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllApplicationRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllApplicationRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllApplicationRequest>, I>>(object: I): QueryAllApplicationRequest {
    const message = createBaseQueryAllApplicationRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllApplicationResponse(): QueryAllApplicationResponse {
  return { application: [], pagination: undefined };
}

export const QueryAllApplicationResponse = {
  encode(message: QueryAllApplicationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.application) {
      Application.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllApplicationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllApplicationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.application.push(Application.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllApplicationResponse {
    return {
      application: Array.isArray(object?.application)
        ? object.application.map((e: any) => Application.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllApplicationResponse): unknown {
    const obj: any = {};
    if (message.application) {
      obj.application = message.application.map((e) => e ? Application.toJSON(e) : undefined);
    } else {
      obj.application = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllApplicationResponse>, I>>(object: I): QueryAllApplicationResponse {
    const message = createBaseQueryAllApplicationResponse();
    message.application = object.application?.map((e) => Application.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a list of Application items. */
  Application(request: QueryGetApplicationRequest): Promise<QueryGetApplicationResponse>;
  ApplicationAll(request: QueryAllApplicationRequest): Promise<QueryAllApplicationResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.Application = this.Application.bind(this);
    this.ApplicationAll = this.ApplicationAll.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.application.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  Application(request: QueryGetApplicationRequest): Promise<QueryGetApplicationResponse> {
    const data = QueryGetApplicationRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.application.Query", "Application", data);
    return promise.then((data) => QueryGetApplicationResponse.decode(new _m0.Reader(data)));
  }

  ApplicationAll(request: QueryAllApplicationRequest): Promise<QueryAllApplicationResponse> {
    const data = QueryAllApplicationRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.application.Query", "ApplicationAll", data);
    return promise.then((data) => QueryAllApplicationResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
