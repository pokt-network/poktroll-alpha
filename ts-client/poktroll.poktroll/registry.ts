import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUnstake } from "./types/poktroll/poktroll/tx";
import { MsgStake } from "./types/poktroll/poktroll/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/poktroll.poktroll.MsgUnstake", MsgUnstake],
    ["/poktroll.poktroll.MsgStake", MsgStake],
    
];

export { msgTypes }