/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "poktroll.poktroll";

export interface MsgStake {
  /** TODO: Reconsider naming to `from` or `address` or `signer` */
  creator: string;
  /** TODO: Figure out denomation for this and consider renaming to `amountToStake` */
  amount: string;
  /** TODO: Change `actorType` to an enum defined in `actor.proto` */
  actorType: string;
}

export interface MsgStakeResponse {
  /** TODO: Need to return success/error codes */
  success: boolean;
}

export interface MsgUnstake {
  /** TODO: Reconsider naming to `from` or `address` or `signer` */
  creator: string;
  /** TODO: Figure out denomation for this */
  amount: string;
  /** TODO: Change `actorType` to an enum defined in `actor.proto` */
  actorType: string;
}

export interface MsgUnstakeResponse {
  /** TODO: Need to return success/error codes */
  success: boolean;
}

function createBaseMsgStake(): MsgStake {
  return { creator: "", amount: "", actorType: "" };
}

export const MsgStake = {
  encode(message: MsgStake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.amount !== "") {
      writer.uint32(18).string(message.amount);
    }
    if (message.actorType !== "") {
      writer.uint32(26).string(message.actorType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.amount = reader.string();
          break;
        case 3:
          message.actorType = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgStake {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      amount: isSet(object.amount) ? String(object.amount) : "",
      actorType: isSet(object.actorType) ? String(object.actorType) : "",
    };
  },

  toJSON(message: MsgStake): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.amount !== undefined && (obj.amount = message.amount);
    message.actorType !== undefined && (obj.actorType = message.actorType);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStake>, I>>(object: I): MsgStake {
    const message = createBaseMsgStake();
    message.creator = object.creator ?? "";
    message.amount = object.amount ?? "";
    message.actorType = object.actorType ?? "";
    return message;
  },
};

function createBaseMsgStakeResponse(): MsgStakeResponse {
  return { success: false };
}

export const MsgStakeResponse = {
  encode(message: MsgStakeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.success === true) {
      writer.uint32(8).bool(message.success);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgStakeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgStakeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.success = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgStakeResponse {
    return { success: isSet(object.success) ? Boolean(object.success) : false };
  },

  toJSON(message: MsgStakeResponse): unknown {
    const obj: any = {};
    message.success !== undefined && (obj.success = message.success);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeResponse>, I>>(object: I): MsgStakeResponse {
    const message = createBaseMsgStakeResponse();
    message.success = object.success ?? false;
    return message;
  },
};

function createBaseMsgUnstake(): MsgUnstake {
  return { creator: "", amount: "", actorType: "" };
}

export const MsgUnstake = {
  encode(message: MsgUnstake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.amount !== "") {
      writer.uint32(18).string(message.amount);
    }
    if (message.actorType !== "") {
      writer.uint32(26).string(message.actorType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.amount = reader.string();
          break;
        case 3:
          message.actorType = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnstake {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      amount: isSet(object.amount) ? String(object.amount) : "",
      actorType: isSet(object.actorType) ? String(object.actorType) : "",
    };
  },

  toJSON(message: MsgUnstake): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.amount !== undefined && (obj.amount = message.amount);
    message.actorType !== undefined && (obj.actorType = message.actorType);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstake>, I>>(object: I): MsgUnstake {
    const message = createBaseMsgUnstake();
    message.creator = object.creator ?? "";
    message.amount = object.amount ?? "";
    message.actorType = object.actorType ?? "";
    return message;
  },
};

function createBaseMsgUnstakeResponse(): MsgUnstakeResponse {
  return { success: false };
}

export const MsgUnstakeResponse = {
  encode(message: MsgUnstakeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.success === true) {
      writer.uint32(8).bool(message.success);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.success = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnstakeResponse {
    return { success: isSet(object.success) ? Boolean(object.success) : false };
  },

  toJSON(message: MsgUnstakeResponse): unknown {
    const obj: any = {};
    message.success !== undefined && (obj.success = message.success);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeResponse>, I>>(object: I): MsgUnstakeResponse {
    const message = createBaseMsgUnstakeResponse();
    message.success = object.success ?? false;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  Stake(request: MsgStake): Promise<MsgStakeResponse>;
  Unstake(request: MsgUnstake): Promise<MsgUnstakeResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Stake = this.Stake.bind(this);
    this.Unstake = this.Unstake.bind(this);
  }
  Stake(request: MsgStake): Promise<MsgStakeResponse> {
    const data = MsgStake.encode(request).finish();
    const promise = this.rpc.request("poktroll.poktroll.Msg", "Stake", data);
    return promise.then((data) => MsgStakeResponse.decode(new _m0.Reader(data)));
  }

  Unstake(request: MsgUnstake): Promise<MsgUnstakeResponse> {
    const data = MsgUnstake.encode(request).finish();
    const promise = this.rpc.request("poktroll.poktroll.Msg", "Unstake", data);
    return promise.then((data) => MsgUnstakeResponse.decode(new _m0.Reader(data)));
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
