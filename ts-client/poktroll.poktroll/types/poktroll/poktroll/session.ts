/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "poktroll.poktroll";

export interface Session {
  /** a universally unique ID for the session */
  id: string;
  /** a monotonically increasing number representing the # on the chain */
  sessionNumber: number;
  /** the height at which the session starts */
  sessionHeight: number;
  /** the number of blocks the session is valid from */
  numSessionBlocks: number;
  /**
   * CONSIDERATION: Should we add a `RelayChain` enum and use it across the board?
   * CONSIDERATION: Should a single session support multiple relay chains?
   * TECHDEBT: Do we need backwards with v0? https://docs.pokt.network/supported-blockchains/
   */
  relayChain: string;
  /** CONSIDERATION: Should a single session support multiple geo zones? */
  geoZone: string;
}

function createBaseSession(): Session {
  return { id: "", sessionNumber: 0, sessionHeight: 0, numSessionBlocks: 0, relayChain: "", geoZone: "" };
}

export const Session = {
  encode(message: Session, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.sessionNumber !== 0) {
      writer.uint32(16).int64(message.sessionNumber);
    }
    if (message.sessionHeight !== 0) {
      writer.uint32(24).int64(message.sessionHeight);
    }
    if (message.numSessionBlocks !== 0) {
      writer.uint32(32).int64(message.numSessionBlocks);
    }
    if (message.relayChain !== "") {
      writer.uint32(42).string(message.relayChain);
    }
    if (message.geoZone !== "") {
      writer.uint32(50).string(message.geoZone);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Session {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSession();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.sessionNumber = longToNumber(reader.int64() as Long);
          break;
        case 3:
          message.sessionHeight = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.numSessionBlocks = longToNumber(reader.int64() as Long);
          break;
        case 5:
          message.relayChain = reader.string();
          break;
        case 6:
          message.geoZone = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Session {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      sessionNumber: isSet(object.sessionNumber) ? Number(object.sessionNumber) : 0,
      sessionHeight: isSet(object.sessionHeight) ? Number(object.sessionHeight) : 0,
      numSessionBlocks: isSet(object.numSessionBlocks) ? Number(object.numSessionBlocks) : 0,
      relayChain: isSet(object.relayChain) ? String(object.relayChain) : "",
      geoZone: isSet(object.geoZone) ? String(object.geoZone) : "",
    };
  },

  toJSON(message: Session): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.sessionNumber !== undefined && (obj.sessionNumber = Math.round(message.sessionNumber));
    message.sessionHeight !== undefined && (obj.sessionHeight = Math.round(message.sessionHeight));
    message.numSessionBlocks !== undefined && (obj.numSessionBlocks = Math.round(message.numSessionBlocks));
    message.relayChain !== undefined && (obj.relayChain = message.relayChain);
    message.geoZone !== undefined && (obj.geoZone = message.geoZone);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Session>, I>>(object: I): Session {
    const message = createBaseSession();
    message.id = object.id ?? "";
    message.sessionNumber = object.sessionNumber ?? 0;
    message.sessionHeight = object.sessionHeight ?? 0;
    message.numSessionBlocks = object.numSessionBlocks ?? 0;
    message.relayChain = object.relayChain ?? "";
    message.geoZone = object.geoZone ?? "";
    return message;
  },
};

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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
