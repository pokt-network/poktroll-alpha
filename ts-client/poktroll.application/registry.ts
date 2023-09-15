import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgStakeApplication } from "./types/poktroll/application/tx";
import { MsgUnstakeApplication } from "./types/poktroll/application/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.application.MsgStakeApplication", MsgStakeApplication],
    ["/poktroll.application.MsgUnstakeApplication", MsgUnstakeApplication],
    
];

export { msgTypes }