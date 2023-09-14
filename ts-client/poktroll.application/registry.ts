import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUnstakeApplication } from "./types/poktroll/application/tx";
import { MsgStakeApplication } from "./types/poktroll/application/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.application.MsgUnstakeApplication", MsgUnstakeApplication],
    ["/poktroll.application.MsgStakeApplication", MsgStakeApplication],
    
];

export { msgTypes }