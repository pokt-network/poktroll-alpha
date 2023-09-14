/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "poktroll.application";

export interface MsgStakeApplication {
  address: string;
  stakeAmount: Coin | undefined;
}

export interface MsgStakeApplicationResponse {
}

export interface MsgUnstakeApplication {
  address: string;
}

export interface MsgUnstakeApplicationResponse {
}

function createBaseMsgStakeApplication(): MsgStakeApplication {
  return { address: "", stakeAmount: undefined };
}

export const MsgStakeApplication = {
  encode(message: MsgStakeApplication, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.stakeAmount !== undefined) {
      Coin.encode(message.stakeAmount, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeApplication {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeApplication();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.stakeAmount = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgStakeApplication {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      stakeAmount: isSet(object.stakeAmount) ? Coin.fromJSON(object.stakeAmount) : undefined,
    };
  },

  toJSON(message: MsgStakeApplication): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.stakeAmount !== undefined
      && (obj.stakeAmount = message.stakeAmount ? Coin.toJSON(message.stakeAmount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeApplication>, I>>(object: I): MsgStakeApplication {
    const message = createBaseMsgStakeApplication();
    message.address = object.address ?? "";
    message.stakeAmount = (object.stakeAmount !== undefined && object.stakeAmount !== null)
      ? Coin.fromPartial(object.stakeAmount)
      : undefined;
    return message;
  },
};

function createBaseMsgStakeApplicationResponse(): MsgStakeApplicationResponse {
  return {};
}

export const MsgStakeApplicationResponse = {
  encode(_: MsgStakeApplicationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeApplicationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeApplicationResponse();
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

  fromJSON(_: any): MsgStakeApplicationResponse {
    return {};
  },

  toJSON(_: MsgStakeApplicationResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeApplicationResponse>, I>>(_: I): MsgStakeApplicationResponse {
    const message = createBaseMsgStakeApplicationResponse();
    return message;
  },
};

function createBaseMsgUnstakeApplication(): MsgUnstakeApplication {
  return { address: "" };
}

export const MsgUnstakeApplication = {
  encode(message: MsgUnstakeApplication, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeApplication {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeApplication();
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

  fromJSON(object: any): MsgUnstakeApplication {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: MsgUnstakeApplication): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeApplication>, I>>(object: I): MsgUnstakeApplication {
    const message = createBaseMsgUnstakeApplication();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseMsgUnstakeApplicationResponse(): MsgUnstakeApplicationResponse {
  return {};
}

export const MsgUnstakeApplicationResponse = {
  encode(_: MsgUnstakeApplicationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeApplicationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeApplicationResponse();
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

  fromJSON(_: any): MsgUnstakeApplicationResponse {
    return {};
  },

  toJSON(_: MsgUnstakeApplicationResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeApplicationResponse>, I>>(_: I): MsgUnstakeApplicationResponse {
    const message = createBaseMsgUnstakeApplicationResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  StakeApplication(request: MsgStakeApplication): Promise<MsgStakeApplicationResponse>;
  UnstakeApplication(request: MsgUnstakeApplication): Promise<MsgUnstakeApplicationResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.StakeApplication = this.StakeApplication.bind(this);
    this.UnstakeApplication = this.UnstakeApplication.bind(this);
  }
  StakeApplication(request: MsgStakeApplication): Promise<MsgStakeApplicationResponse> {
    const data = MsgStakeApplication.encode(request).finish();
    const promise = this.rpc.request("poktroll.application.Msg", "StakeApplication", data);
    return promise.then((data) => MsgStakeApplicationResponse.decode(new _m0.Reader(data)));
  }

  UnstakeApplication(request: MsgUnstakeApplication): Promise<MsgUnstakeApplicationResponse> {
    const data = MsgUnstakeApplication.encode(request).finish();
    const promise = this.rpc.request("poktroll.application.Msg", "UnstakeApplication", data);
    return promise.then((data) => MsgUnstakeApplicationResponse.decode(new _m0.Reader(data)));
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
