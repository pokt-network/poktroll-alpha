/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "poktroll.session";

export interface MsgGetSession {
  address: string;
}

export interface MsgGetSessionResponse {
}

function createBaseMsgGetSession(): MsgGetSession {
  return { address: "" };
}

export const MsgGetSession = {
  encode(message: MsgGetSession, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGetSession {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGetSession();
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

  fromJSON(object: any): MsgGetSession {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: MsgGetSession): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgGetSession>, I>>(object: I): MsgGetSession {
    const message = createBaseMsgGetSession();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseMsgGetSessionResponse(): MsgGetSessionResponse {
  return {};
}

export const MsgGetSessionResponse = {
  encode(_: MsgGetSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgGetSessionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGetSessionResponse();
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

  fromJSON(_: any): MsgGetSessionResponse {
    return {};
  },

  toJSON(_: MsgGetSessionResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgGetSessionResponse>, I>>(_: I): MsgGetSessionResponse {
    const message = createBaseMsgGetSessionResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  GetSession(request: MsgGetSession): Promise<MsgGetSessionResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetSession = this.GetSession.bind(this);
  }
  GetSession(request: MsgGetSession): Promise<MsgGetSessionResponse> {
    const data = MsgGetSession.encode(request).finish();
    const promise = this.rpc.request("poktroll.session.Msg", "GetSession", data);
    return promise.then((data) => MsgGetSessionResponse.decode(new _m0.Reader(data)));
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
