/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Application } from "../application/application";
import { Servicers } from "../servicer/servicers";

export const protobufPackage = "poktroll.session";

export interface Session {
  application: Application | undefined;
  servicers: Servicers[];
}

function createBaseSession(): Session {
  return { application: undefined, servicers: [] };
}

export const Session = {
  encode(message: Session, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.application !== undefined) {
      Application.encode(message.application, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.servicers) {
      Servicers.encode(v!, writer.uint32(18).fork()).ldelim();
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
          message.application = Application.decode(reader, reader.uint32());
          break;
        case 2:
          message.servicers.push(Servicers.decode(reader, reader.uint32()));
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
      application: isSet(object.application) ? Application.fromJSON(object.application) : undefined,
      servicers: Array.isArray(object?.servicers) ? object.servicers.map((e: any) => Servicers.fromJSON(e)) : [],
    };
  },

  toJSON(message: Session): unknown {
    const obj: any = {};
    message.application !== undefined
      && (obj.application = message.application ? Application.toJSON(message.application) : undefined);
    if (message.servicers) {
      obj.servicers = message.servicers.map((e) => e ? Servicers.toJSON(e) : undefined);
    } else {
      obj.servicers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Session>, I>>(object: I): Session {
    const message = createBaseSession();
    message.application = (object.application !== undefined && object.application !== null)
      ? Application.fromPartial(object.application)
      : undefined;
    message.servicers = object.servicers?.map((e) => Servicers.fromPartial(e)) || [];
    return message;
  },
};

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
