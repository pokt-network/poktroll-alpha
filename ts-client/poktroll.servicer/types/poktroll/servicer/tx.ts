/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "poktroll.servicer";

export interface MsgStakeServicer {
  address: string;
  stakeAmount: Coin | undefined;
}

export interface MsgStakeServicerResponse {
}

export interface MsgUnstakeServicer {
  address: string;
}

export interface MsgUnstakeServicerResponse {
}

export interface MsgClaim {
  creator: string;
  smtRootHash: Uint8Array;
}

export interface MsgClaimResponse {
}

export interface MsgProof {
  creator: string;
  root: string;
  path: string;
  valueHash: string;
  sum: number;
  proofBz: string;
}

export interface MsgProofResponse {
}

function createBaseMsgStakeServicer(): MsgStakeServicer {
  return { address: "", stakeAmount: undefined };
}

export const MsgStakeServicer = {
  encode(message: MsgStakeServicer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.stakeAmount !== undefined) {
      Coin.encode(message.stakeAmount, writer.uint32(18).fork()).ldelim();
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

  fromJSON(object: any): MsgStakeServicer {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      stakeAmount: isSet(object.stakeAmount) ? Coin.fromJSON(object.stakeAmount) : undefined,
    };
  },

  toJSON(message: MsgStakeServicer): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.stakeAmount !== undefined
      && (obj.stakeAmount = message.stakeAmount ? Coin.toJSON(message.stakeAmount) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgStakeServicer>, I>>(object: I): MsgStakeServicer {
    const message = createBaseMsgStakeServicer();
    message.address = object.address ?? "";
    message.stakeAmount = (object.stakeAmount !== undefined && object.stakeAmount !== null)
      ? Coin.fromPartial(object.stakeAmount)
      : undefined;
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

function createBaseMsgUnstakeServicer(): MsgUnstakeServicer {
  return { address: "" };
}

export const MsgUnstakeServicer = {
  encode(message: MsgUnstakeServicer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeServicer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeServicer();
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

  fromJSON(object: any): MsgUnstakeServicer {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: MsgUnstakeServicer): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeServicer>, I>>(object: I): MsgUnstakeServicer {
    const message = createBaseMsgUnstakeServicer();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseMsgUnstakeServicerResponse(): MsgUnstakeServicerResponse {
  return {};
}

export const MsgUnstakeServicerResponse = {
  encode(_: MsgUnstakeServicerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnstakeServicerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnstakeServicerResponse();
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

  fromJSON(_: any): MsgUnstakeServicerResponse {
    return {};
  },

  toJSON(_: MsgUnstakeServicerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnstakeServicerResponse>, I>>(_: I): MsgUnstakeServicerResponse {
    const message = createBaseMsgUnstakeServicerResponse();
    return message;
  },
};

function createBaseMsgClaim(): MsgClaim {
  return { creator: "", smtRootHash: new Uint8Array() };
}

export const MsgClaim = {
  encode(message: MsgClaim, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.smtRootHash.length !== 0) {
      writer.uint32(18).bytes(message.smtRootHash);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaim {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaim();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.smtRootHash = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgClaim {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      smtRootHash: isSet(object.smtRootHash) ? bytesFromBase64(object.smtRootHash) : new Uint8Array(),
    };
  },

  toJSON(message: MsgClaim): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.smtRootHash !== undefined
      && (obj.smtRootHash = base64FromBytes(
        message.smtRootHash !== undefined ? message.smtRootHash : new Uint8Array(),
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaim>, I>>(object: I): MsgClaim {
    const message = createBaseMsgClaim();
    message.creator = object.creator ?? "";
    message.smtRootHash = object.smtRootHash ?? new Uint8Array();
    return message;
  },
};

function createBaseMsgClaimResponse(): MsgClaimResponse {
  return {};
}

export const MsgClaimResponse = {
  encode(_: MsgClaimResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimResponse();
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

  fromJSON(_: any): MsgClaimResponse {
    return {};
  },

  toJSON(_: MsgClaimResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaimResponse>, I>>(_: I): MsgClaimResponse {
    const message = createBaseMsgClaimResponse();
    return message;
  },
};

function createBaseMsgProof(): MsgProof {
  return { creator: "", root: "", path: "", valueHash: "", sum: 0, proofBz: "" };
}

export const MsgProof = {
  encode(message: MsgProof, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.root !== "") {
      writer.uint32(18).string(message.root);
    }
    if (message.path !== "") {
      writer.uint32(26).string(message.path);
    }
    if (message.valueHash !== "") {
      writer.uint32(34).string(message.valueHash);
    }
    if (message.sum !== 0) {
      writer.uint32(40).int32(message.sum);
    }
    if (message.proofBz !== "") {
      writer.uint32(50).string(message.proofBz);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgProof {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgProof();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.root = reader.string();
          break;
        case 3:
          message.path = reader.string();
          break;
        case 4:
          message.valueHash = reader.string();
          break;
        case 5:
          message.sum = reader.int32();
          break;
        case 6:
          message.proofBz = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgProof {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      root: isSet(object.root) ? String(object.root) : "",
      path: isSet(object.path) ? String(object.path) : "",
      valueHash: isSet(object.valueHash) ? String(object.valueHash) : "",
      sum: isSet(object.sum) ? Number(object.sum) : 0,
      proofBz: isSet(object.proofBz) ? String(object.proofBz) : "",
    };
  },

  toJSON(message: MsgProof): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.root !== undefined && (obj.root = message.root);
    message.path !== undefined && (obj.path = message.path);
    message.valueHash !== undefined && (obj.valueHash = message.valueHash);
    message.sum !== undefined && (obj.sum = Math.round(message.sum));
    message.proofBz !== undefined && (obj.proofBz = message.proofBz);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgProof>, I>>(object: I): MsgProof {
    const message = createBaseMsgProof();
    message.creator = object.creator ?? "";
    message.root = object.root ?? "";
    message.path = object.path ?? "";
    message.valueHash = object.valueHash ?? "";
    message.sum = object.sum ?? 0;
    message.proofBz = object.proofBz ?? "";
    return message;
  },
};

function createBaseMsgProofResponse(): MsgProofResponse {
  return {};
}

export const MsgProofResponse = {
  encode(_: MsgProofResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgProofResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgProofResponse();
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

  fromJSON(_: any): MsgProofResponse {
    return {};
  },

  toJSON(_: MsgProofResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgProofResponse>, I>>(_: I): MsgProofResponse {
    const message = createBaseMsgProofResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  StakeServicer(request: MsgStakeServicer): Promise<MsgStakeServicerResponse>;
  UnstakeServicer(request: MsgUnstakeServicer): Promise<MsgUnstakeServicerResponse>;
  Claim(request: MsgClaim): Promise<MsgClaimResponse>;
  Proof(request: MsgProof): Promise<MsgProofResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.StakeServicer = this.StakeServicer.bind(this);
    this.UnstakeServicer = this.UnstakeServicer.bind(this);
    this.Claim = this.Claim.bind(this);
    this.Proof = this.Proof.bind(this);
  }
  StakeServicer(request: MsgStakeServicer): Promise<MsgStakeServicerResponse> {
    const data = MsgStakeServicer.encode(request).finish();
    const promise = this.rpc.request("poktroll.servicer.Msg", "StakeServicer", data);
    return promise.then((data) => MsgStakeServicerResponse.decode(new _m0.Reader(data)));
  }

  UnstakeServicer(request: MsgUnstakeServicer): Promise<MsgUnstakeServicerResponse> {
    const data = MsgUnstakeServicer.encode(request).finish();
    const promise = this.rpc.request("poktroll.servicer.Msg", "UnstakeServicer", data);
    return promise.then((data) => MsgUnstakeServicerResponse.decode(new _m0.Reader(data)));
  }

  Claim(request: MsgClaim): Promise<MsgClaimResponse> {
    const data = MsgClaim.encode(request).finish();
    const promise = this.rpc.request("poktroll.servicer.Msg", "Claim", data);
    return promise.then((data) => MsgClaimResponse.decode(new _m0.Reader(data)));
  }

  Proof(request: MsgProof): Promise<MsgProofResponse> {
    const data = MsgProof.encode(request).finish();
    const promise = this.rpc.request("poktroll.servicer.Msg", "Proof", data);
    return promise.then((data) => MsgProofResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
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
