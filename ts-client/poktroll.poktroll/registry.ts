import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgStake } from "./types/poktroll/poktroll/tx";
import { MsgUnstake } from "./types/poktroll/poktroll/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.poktroll.MsgStake", MsgStake],
    ["/poktroll.poktroll.MsgUnstake", MsgUnstake],
    
];

export { msgTypes }