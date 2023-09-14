/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "poktroll.servicer";

export interface MsgStakeServicer {
  address: string;
}

export interface MsgStakeServicerResponse {
}

function createBaseMsgStakeServicer(): MsgStakeServicer {
  return { address: "" };
}

export const MsgStakeServicer = {
  encode(message: MsgStakeServicer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeServicer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeServicer();
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

  fromJSON(object: any): MsgStakeServicer {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: MsgStakeServicer): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeServicer>, I>>(object: I): MsgStakeServicer {
    const message = createBaseMsgStakeServicer();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseMsgStakeServicerResponse(): MsgStakeServicerResponse {
  return {};
}

export const MsgStakeServicerResponse = {
  encode(_: MsgStakeServicerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeServicerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeServicerResponse();
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

  fromJSON(_: any): MsgStakeServicerResponse {
    return {};
  },

  toJSON(_: MsgStakeServicerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeServicerResponse>, I>>(_: I): MsgStakeServicerResponse {
    const message = createBaseMsgStakeServicerResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  StakeServicer(request: MsgStakeServicer): Promise<MsgStakeServicerResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.StakeServicer = this.StakeServicer.bind(this);
  }
  StakeServicer(request: MsgStakeServicer): Promise<MsgStakeServicerResponse> {
    const data = MsgStakeServicer.encode(request).finish();
    const promise = this.rpc.request("poktroll.servicer.Msg", "StakeServicer", data);
    return promise.then((data) => MsgStakeServicerResponse.decode(new _m0.Reader(data)));
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
