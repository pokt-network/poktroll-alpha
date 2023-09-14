/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Servicer } from "./actor";
import { Params } from "./params";

export const protobufPackage = "poktroll.poktroll";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetSessionRequest {
}

export interface QueryGetSessionResponse {
  placeholderText: string;
}

export interface QueryServicersRequest {
  pagination: PageRequest | undefined;
}

export interface QueryServicersResponse {
  servicers: Servicer[];
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

function createBaseQueryGetSessionRequest(): QueryGetSessionRequest {
  return {};
}

export const QueryGetSessionRequest = {
  encode(_: QueryGetSessionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSessionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSessionRequest();
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

  fromJSON(_: any): QueryGetSessionRequest {
    return {};
  },

  toJSON(_: QueryGetSessionRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSessionRequest>, I>>(_: I): QueryGetSessionRequest {
    const message = createBaseQueryGetSessionRequest();
    return message;
  },
};

function createBaseQueryGetSessionResponse(): QueryGetSessionResponse {
  return { placeholderText: "" };
}

export const QueryGetSessionResponse = {
  encode(message: QueryGetSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.placeholderText !== "") {
      writer.uint32(10).string(message.placeholderText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSessionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSessionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.placeholderText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSessionResponse {
    return { placeholderText: isSet(object.placeholderText) ? String(object.placeholderText) : "" };
  },

  toJSON(message: QueryGetSessionResponse): unknown {
    const obj: any = {};
    message.placeholderText !== undefined && (obj.placeholderText = message.placeholderText);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSessionResponse>, I>>(object: I): QueryGetSessionResponse {
    const message = createBaseQueryGetSessionResponse();
    message.placeholderText = object.placeholderText ?? "";
    return message;
  },
};

function createBaseQueryServicersRequest(): QueryServicersRequest {
  return { pagination: undefined };
}

export const QueryServicersRequest = {
  encode(message: QueryServicersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryServicersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryServicersRequest();
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

  fromJSON(object: any): QueryServicersRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryServicersRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryServicersRequest>, I>>(object: I): QueryServicersRequest {
    const message = createBaseQueryServicersRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryServicersResponse(): QueryServicersResponse {
  return { servicers: [], pagination: undefined };
}

export const QueryServicersResponse = {
  encode(message: QueryServicersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.servicers) {
      Servicer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryServicersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryServicersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.servicers.push(Servicer.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QueryServicersResponse {
    return {
      servicers: Array.isArray(object?.servicers) ? object.servicers.map((e: any) => Servicer.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryServicersResponse): unknown {
    const obj: any = {};
    if (message.servicers) {
      obj.servicers = message.servicers.map((e) => e ? Servicer.toJSON(e) : undefined);
    } else {
      obj.servicers = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryServicersResponse>, I>>(object: I): QueryServicersResponse {
    const message = createBaseQueryServicersResponse();
    message.servicers = object.servicers?.map((e) => Servicer.fromPartial(e)) || [];
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
  /** Queries a list of GetSession items. */
  GetSession(request: QueryGetSessionRequest): Promise<QueryGetSessionResponse>;
  /** Queries for all staked Servicers */
  Servicers(request: QueryServicersRequest): Promise<QueryServicersResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.GetSession = this.GetSession.bind(this);
    this.Servicers = this.Servicers.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.poktroll.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  GetSession(request: QueryGetSessionRequest): Promise<QueryGetSessionResponse> {
    const data = QueryGetSessionRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.poktroll.Query", "GetSession", data);
    return promise.then((data) => QueryGetSessionResponse.decode(new _m0.Reader(data)));
  }

  Servicers(request: QueryServicersRequest): Promise<QueryServicersResponse> {
    const data = QueryServicersRequest.encode(request).finish();
    const promise = this.rpc.request("poktroll.poktroll.Query", "Servicers", data);
    return promise.then((data) => QueryServicersResponse.decode(new _m0.Reader(data)));
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
