import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUnstakeServicer } from "./types/poktroll/servicer/tx";
import { MsgStakeServicer } from "./types/poktroll/servicer/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.servicer.MsgUnstakeServicer", MsgUnstakeServicer],
    ["/poktroll.servicer.MsgStakeServicer", MsgStakeServicer],
    
];

export { msgTypes }